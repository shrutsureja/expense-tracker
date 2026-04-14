package service

import (
	"database/sql"
	"fmt"
	"time"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

type Summary struct {
	TotalAmount float64 `json:"total_amount"`
	Count       int     `json:"count"`
	AvgPerDay   float64 `json:"avg_per_day"`
}

type CategoryBreakdown struct {
	TagID      int64   `json:"tag_id"`
	TagName    string  `json:"tag_name"`
	TagIcon    string  `json:"tag_icon"`
	Total      float64 `json:"total"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type PersonBreakdown struct {
	UserID      int64   `json:"user_id"`
	DisplayName string  `json:"display_name"`
	Total       float64 `json:"total"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

type PaymentMethodBreakdown struct {
	PaymentMethod string  `json:"payment_method"`
	Total         float64 `json:"total"`
	Count         int     `json:"count"`
	Percentage    float64 `json:"percentage"`
}

type DailyTotal struct {
	Date   string  `json:"date"`
	Total  float64 `json:"total"`
	Count  int     `json:"count"`
}

type MonthlyComparison struct {
	ThisMonth      float64 `json:"this_month"`
	LastMonth      float64 `json:"last_month"`
	ThisMonthCount int     `json:"this_month_count"`
	LastMonthCount int     `json:"last_month_count"`
	ChangePercent  float64 `json:"change_percent"`
}

// ResolveDateRange converts range names to actual date strings
func ResolveDateRange(rangeType, from, to string) (string, string) {
	now := time.Now()
	today := now.Format("2006-01-02")

	switch rangeType {
	case "today":
		return today, today
	case "week":
		weekStart := now.AddDate(0, 0, -int(now.Weekday()))
		return weekStart.Format("2006-01-02"), today
	case "month":
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return monthStart.Format("2006-01-02"), today
	case "custom":
		if from != "" && to != "" {
			return from, to
		}
		// Default to this month
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return monthStart.Format("2006-01-02"), today
	default:
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return monthStart.Format("2006-01-02"), today
	}
}

func (s *AnalyticsService) GetSummary(familyID int64, dateFrom, dateTo string) (*Summary, error) {
	var totalAmount sql.NullFloat64
	var count int
	err := s.db.QueryRow(
		`SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM expenses
		 WHERE family_id = ? AND expense_date >= ? AND expense_date <= ?`,
		familyID, dateFrom, dateTo,
	).Scan(&totalAmount, &count)
	if err != nil {
		return nil, err
	}

	// Calculate days in range for avg
	from, _ := time.Parse("2006-01-02", dateFrom)
	to, _ := time.Parse("2006-01-02", dateTo)
	days := to.Sub(from).Hours()/24 + 1
	if days < 1 {
		days = 1
	}

	total := 0.0
	if totalAmount.Valid {
		total = totalAmount.Float64
	}

	return &Summary{
		TotalAmount: total,
		Count:       count,
		AvgPerDay:   total / days,
	}, nil
}

func (s *AnalyticsService) GetByPerson(familyID int64, dateFrom, dateTo string) ([]PersonBreakdown, error) {
	rows, err := s.db.Query(
		`SELECT e.user_id, u.display_name, SUM(e.amount) as total, COUNT(*) as cnt
		 FROM expenses e JOIN users u ON e.user_id = u.id
		 WHERE e.family_id = ? AND e.expense_date >= ? AND e.expense_date <= ?
		 GROUP BY e.user_id ORDER BY total DESC`,
		familyID, dateFrom, dateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PersonBreakdown
	var grandTotal float64
	for rows.Next() {
		var pb PersonBreakdown
		if err := rows.Scan(&pb.UserID, &pb.DisplayName, &pb.Total, &pb.Count); err != nil {
			return nil, err
		}
		grandTotal += pb.Total
		results = append(results, pb)
	}

	for i := range results {
		if grandTotal > 0 {
			results[i].Percentage = (results[i].Total / grandTotal) * 100
		}
	}

	return results, rows.Err()
}

func (s *AnalyticsService) GetByCategory(familyID int64, dateFrom, dateTo string) ([]CategoryBreakdown, error) {
	rows, err := s.db.Query(
		`SELECT e.tag_id, t.name, t.icon, SUM(e.amount) as total, COUNT(*) as cnt
		 FROM expenses e JOIN tags t ON e.tag_id = t.id
		 WHERE e.family_id = ? AND e.expense_date >= ? AND e.expense_date <= ?
		 GROUP BY e.tag_id ORDER BY total DESC`,
		familyID, dateFrom, dateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CategoryBreakdown
	var grandTotal float64
	for rows.Next() {
		var cb CategoryBreakdown
		if err := rows.Scan(&cb.TagID, &cb.TagName, &cb.TagIcon, &cb.Total, &cb.Count); err != nil {
			return nil, err
		}
		grandTotal += cb.Total
		results = append(results, cb)
	}

	for i := range results {
		if grandTotal > 0 {
			results[i].Percentage = (results[i].Total / grandTotal) * 100
		}
	}

	return results, rows.Err()
}

func (s *AnalyticsService) GetByPaymentMethod(familyID int64, dateFrom, dateTo string) ([]PaymentMethodBreakdown, error) {
	rows, err := s.db.Query(
		`SELECT payment_method, SUM(amount) as total, COUNT(*) as cnt
		 FROM expenses
		 WHERE family_id = ? AND expense_date >= ? AND expense_date <= ?
		 GROUP BY payment_method ORDER BY total DESC`,
		familyID, dateFrom, dateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PaymentMethodBreakdown
	var grandTotal float64
	for rows.Next() {
		var pm PaymentMethodBreakdown
		if err := rows.Scan(&pm.PaymentMethod, &pm.Total, &pm.Count); err != nil {
			return nil, err
		}
		grandTotal += pm.Total
		results = append(results, pm)
	}

	for i := range results {
		if grandTotal > 0 {
			results[i].Percentage = (results[i].Total / grandTotal) * 100
		}
	}

	return results, rows.Err()
}

func (s *AnalyticsService) GetDailyTrend(familyID int64, month string) ([]DailyTotal, error) {
	// month is "YYYY-MM"
	dateFrom := fmt.Sprintf("%s-01", month)
	t, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return nil, err
	}
	dateTo := t.AddDate(0, 1, -1).Format("2006-01-02")

	rows, err := s.db.Query(
		`SELECT expense_date, SUM(amount) as total, COUNT(*) as cnt
		 FROM expenses
		 WHERE family_id = ? AND expense_date >= ? AND expense_date <= ?
		 GROUP BY expense_date ORDER BY expense_date`,
		familyID, dateFrom, dateTo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DailyTotal
	for rows.Next() {
		var dt DailyTotal
		if err := rows.Scan(&dt.Date, &dt.Total, &dt.Count); err != nil {
			return nil, err
		}
		results = append(results, dt)
	}

	return results, rows.Err()
}

func (s *AnalyticsService) GetMonthlyComparison(familyID int64) (*MonthlyComparison, error) {
	now := time.Now()
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)
	lastMonthEnd := thisMonthStart.AddDate(0, 0, -1)

	mc := &MonthlyComparison{}

	err := s.db.QueryRow(
		`SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM expenses
		 WHERE family_id = ? AND expense_date >= ? AND expense_date <= ?`,
		familyID, thisMonthStart.Format("2006-01-02"), now.Format("2006-01-02"),
	).Scan(&mc.ThisMonth, &mc.ThisMonthCount)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(
		`SELECT COALESCE(SUM(amount), 0), COUNT(*) FROM expenses
		 WHERE family_id = ? AND expense_date >= ? AND expense_date <= ?`,
		familyID, lastMonthStart.Format("2006-01-02"), lastMonthEnd.Format("2006-01-02"),
	).Scan(&mc.LastMonth, &mc.LastMonthCount)
	if err != nil {
		return nil, err
	}

	if mc.LastMonth > 0 {
		mc.ChangePercent = ((mc.ThisMonth - mc.LastMonth) / mc.LastMonth) * 100
	}

	return mc, nil
}
