package models

import "time"

type Expense struct {
	ID              int64     `json:"id"`
	FamilyID        int64     `json:"family_id"`
	UserID          int64     `json:"user_id"`
	Amount          float64   `json:"amount"`
	TagID           int64     `json:"tag_id"`
	PaymentMethod   string    `json:"payment_method"`
	Note            string    `json:"note"`
	ExpenseDate     string    `json:"expense_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserDisplayName string    `json:"user_display_name,omitempty"`
	TagName         string    `json:"tag_name,omitempty"`
	TagIcon         string    `json:"tag_icon,omitempty"`
}

type ExpenseFilters struct {
	UserID        *int64
	TagID         *int64
	PaymentMethod *string
	DateFrom      *string
	DateTo        *string
	Page          int
	PageSize      int
}
