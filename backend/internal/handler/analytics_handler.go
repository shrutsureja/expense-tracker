package handler

import (
	"net/http"

	"expense-tracker/internal/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) Summary(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	rangeType := r.URL.Query().Get("range")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	dateFrom, dateTo := service.ResolveDateRange(rangeType, from, to)

	summary, err := h.analyticsService.GetSummary(*claims.FamilyID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get summary"})
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

func (h *AnalyticsHandler) ByPerson(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	dateFrom, dateTo := service.ResolveDateRange(r.URL.Query().Get("range"), r.URL.Query().Get("from"), r.URL.Query().Get("to"))

	result, err := h.analyticsService.GetByPerson(*claims.FamilyID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get person breakdown"})
		return
	}
	if result == nil {
		result = []service.PersonBreakdown{}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) ByCategory(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	dateFrom, dateTo := service.ResolveDateRange(r.URL.Query().Get("range"), r.URL.Query().Get("from"), r.URL.Query().Get("to"))

	result, err := h.analyticsService.GetByCategory(*claims.FamilyID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get category breakdown"})
		return
	}
	if result == nil {
		result = []service.CategoryBreakdown{}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) ByPaymentMethod(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	dateFrom, dateTo := service.ResolveDateRange(r.URL.Query().Get("range"), r.URL.Query().Get("from"), r.URL.Query().Get("to"))

	result, err := h.analyticsService.GetByPaymentMethod(*claims.FamilyID, dateFrom, dateTo)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get payment method breakdown"})
		return
	}
	if result == nil {
		result = []service.PaymentMethodBreakdown{}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) DailyTrend(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	month := r.URL.Query().Get("month")
	if month == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "month parameter is required (YYYY-MM)"})
		return
	}

	result, err := h.analyticsService.GetDailyTrend(*claims.FamilyID, month)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get daily trend"})
		return
	}
	if result == nil {
		result = []service.DailyTotal{}
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *AnalyticsHandler) MonthlyComparison(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims.FamilyID == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no family associated"})
		return
	}

	result, err := h.analyticsService.GetMonthlyComparison(*claims.FamilyID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get monthly comparison"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}
