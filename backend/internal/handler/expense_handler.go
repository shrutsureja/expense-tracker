package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"expense-tracker/internal/models"
	"expense-tracker/internal/service"

	"github.com/go-chi/chi/v5"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

type addExpenseRequest struct {
	Amount        float64 `json:"amount"`
	TagID         int64   `json:"tag_id"`
	PaymentMethod string  `json:"payment_method"`
	Note          string  `json:"note"`
	ExpenseDate   string  `json:"expense_date"`
}

func (h *ExpenseHandler) AddExpense(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	var req addExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	expense, err := h.expenseService.AddExpense(
		*claims.FamilyID, claims.UserID, req.Amount, req.TagID,
		req.PaymentMethod, req.Note, req.ExpenseDate,
	)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, expense)
}

func (h *ExpenseHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	filters := models.ExpenseFilters{
		Page:     1,
		PageSize: 20,
	}

	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			filters.Page = p
		}
	}
	if v := r.URL.Query().Get("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil {
			filters.PageSize = ps
		}
	}
	if v := r.URL.Query().Get("user_id"); v != "" {
		if uid, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.UserID = &uid
		}
	}
	if v := r.URL.Query().Get("tag_id"); v != "" {
		if tid, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.TagID = &tid
		}
	}
	if v := r.URL.Query().Get("payment_method"); v != "" {
		filters.PaymentMethod = &v
	}
	if v := r.URL.Query().Get("date_from"); v != "" {
		filters.DateFrom = &v
	}
	if v := r.URL.Query().Get("date_to"); v != "" {
		filters.DateTo = &v
	}

	expenses, total, err := h.expenseService.ListExpenses(*claims.FamilyID, filters)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list expenses"})
		return
	}
	if expenses == nil {
		expenses = []models.Expense{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"expenses": expenses,
		"total":    total,
		"page":     filters.Page,
		"page_size": filters.PageSize,
	})
}

func (h *ExpenseHandler) RecentExpenses(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	expenses, err := h.expenseService.GetRecentExpenses(*claims.FamilyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get recent expenses"})
		return
	}
	if expenses == nil {
		expenses = []models.Expense{}
	}

	writeJSON(w, http.StatusOK, expenses)
}

type updateExpenseRequest struct {
	Amount        float64 `json:"amount"`
	TagID         int64   `json:"tag_id"`
	PaymentMethod string  `json:"payment_method"`
	Note          string  `json:"note"`
	ExpenseDate   string  `json:"expense_date"`
}

func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
		return
	}

	var req updateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.expenseService.UpdateExpense(id, claims.UserID, claims.Role, *claims.FamilyID, req.Amount, req.TagID, req.PaymentMethod, req.Note, req.ExpenseDate); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "expense updated"})
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid expense id"})
		return
	}

	if err := h.expenseService.DeleteExpense(id, claims.UserID, claims.Role, *claims.FamilyID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "expense deleted"})
}
