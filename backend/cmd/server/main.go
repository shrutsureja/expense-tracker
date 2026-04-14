package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"expense-tracker/internal/config"
	"expense-tracker/internal/database"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"expense-tracker/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := database.SeedSuperAdmin(db, cfg.Auth.SuperAdmin.Username, cfg.Auth.SuperAdmin.Password); err != nil {
		log.Fatalf("Failed to seed super admin: %v", err)
	}
	if err := database.SeedDefaultTags(db); err != nil {
		log.Fatalf("Failed to seed default tags: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	familyRepo := repository.NewFamilyRepository(db)
	tagRepo := repository.NewTagRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.TokenExpiryDuration())
	familyService := service.NewFamilyService(familyRepo, userRepo)
	expenseService := service.NewExpenseService(expenseRepo)
	analyticsService := service.NewAnalyticsService(db)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	familyHandler := handler.NewFamilyHandler(familyService)
	tagHandler := handler.NewTagHandler(tagRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})

		// Public routes
		r.Post("/auth/login", authHandler.Login)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware(cfg.Auth.JWTSecret))

			// Auth
			r.Get("/auth/me", authHandler.Me)

			// Super admin routes
			r.Route("/admin", func(r chi.Router) {
				r.Use(handler.RequireRole(models.RoleSuperAdmin))
				r.Get("/families", familyHandler.ListFamilies)
				r.Post("/families", familyHandler.CreateFamily)
				r.Delete("/families/{id}", familyHandler.DeleteFamily)
			})

			// Family routes (owner + members)
			r.Route("/family", func(r chi.Router) {
				r.Use(handler.RequireRole(models.RoleFamilyOwner, models.RoleFamilyMember))
				r.Get("/", familyHandler.GetFamily)
				r.Get("/members", familyHandler.ListMembers)

				// Owner only
				r.Group(func(r chi.Router) {
					r.Use(handler.RequireRole(models.RoleFamilyOwner))
					r.Post("/members", familyHandler.AddMember)
					r.Put("/members/{id}", familyHandler.UpdateMember)
					r.Delete("/members/{id}", familyHandler.DeactivateMember)
				})
			})

			// Tags
			r.Get("/tags", tagHandler.ListTags)
			r.Get("/tags/suggested", tagHandler.SuggestedTags)
			r.Group(func(r chi.Router) {
				r.Use(handler.RequireRole(models.RoleFamilyOwner))
				r.Post("/tags", tagHandler.CreateTag)
			})

			// Expenses
			r.Route("/expenses", func(r chi.Router) {
				r.Use(handler.RequireRole(models.RoleFamilyOwner, models.RoleFamilyMember))
				r.Post("/", expenseHandler.AddExpense)
				r.Get("/", expenseHandler.ListExpenses)
				r.Get("/recent", expenseHandler.RecentExpenses)
				r.Put("/{id}", expenseHandler.UpdateExpense)
				r.Delete("/{id}", expenseHandler.DeleteExpense)
			})

			// Analytics
			r.Route("/analytics", func(r chi.Router) {
				r.Use(handler.RequireRole(models.RoleFamilyOwner, models.RoleFamilyMember))
				r.Get("/summary", analyticsHandler.Summary)
				r.Get("/by-person", analyticsHandler.ByPerson)
				r.Get("/by-category", analyticsHandler.ByCategory)
				r.Get("/by-payment-method", analyticsHandler.ByPaymentMethod)
				r.Get("/daily-trend", analyticsHandler.DailyTrend)
				r.Get("/monthly-comparison", analyticsHandler.MonthlyComparison)
			})
		})
	})

	// Serve React SPA from embedded static files
	r.Handle("/*", spaHandler())

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

//go:embed static
var staticFS embed.FS

func spaHandler() http.HandlerFunc {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		// static directory doesn't exist yet (dev mode) — serve a placeholder
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html><body><h1>Frontend not built yet</h1><p>Run <code>make build-frontend</code> or <code>cd frontend && npm run dev</code></p></body></html>"))
		}
	}

	fileServer := http.FileServer(http.FS(sub))

	return func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		if _, err := fs.Stat(sub, path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for all non-file routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	}
}
