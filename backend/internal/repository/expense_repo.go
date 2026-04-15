package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"expense-tracker/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func scanExpense(row interface {
	Scan(dest ...any) error
}) (*models.Expense, error) {
	e := &models.Expense{}
	var note sql.NullString
	if err := row.Scan(
		&e.ID, &e.FamilyID, &e.UserID, &e.Amount, &e.TagID, &e.PaymentMethod, &note,
		&e.ExpenseDate, &e.CreatedAt, &e.UpdatedAt,
		&e.UserDisplayName, &e.TagName, &e.TagIcon,
	); err != nil {
		return nil, err
	}
	if note.Valid {
		e.Note = note.String
	}
	return e, nil
}

const expenseSelectCols = `
	e.id, e.family_id, e.user_id, e.amount, e.tag_id, e.payment_method, e.note,
	e.expense_date, e.created_at, e.updated_at,
	u.display_name, t.name, t.icon`

const expenseJoins = `
	FROM expenses e
	JOIN users u ON e.user_id = u.id
	JOIN tags t ON e.tag_id = t.id`

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	var noteVal interface{}
	if expense.Note == "" {
		noteVal = nil
	} else {
		noteVal = expense.Note
	}

	result, err := tx.Exec(
		`INSERT INTO expenses (family_id, user_id, amount, tag_id, payment_method, note, expense_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		expense.FamilyID, expense.UserID, expense.Amount, expense.TagID,
		expense.PaymentMethod, noteVal, expense.ExpenseDate, now, now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	expense.ID = id

	// Update tag frequency
	_, err = tx.Exec(
		`INSERT INTO user_tag_frequency (user_id, tag_id, usage_count, last_used_at)
		 VALUES (?, ?, 1, ?)
		 ON CONFLICT(user_id, tag_id) DO UPDATE SET usage_count = usage_count + 1, last_used_at = ?`,
		expense.UserID, expense.TagID, now, now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ExpenseRepository) GetByID(id int64) (*models.Expense, error) {
	row := r.db.QueryRow(
		`SELECT `+expenseSelectCols+expenseJoins+` WHERE e.id = ?`, id,
	)
	return scanExpense(row)
}

func (r *ExpenseRepository) Update(expense *models.Expense) error {
	var noteVal interface{}
	if expense.Note == "" {
		noteVal = nil
	} else {
		noteVal = expense.Note
	}
	_, err := r.db.Exec(
		`UPDATE expenses SET amount = ?, tag_id = ?, payment_method = ?, note = ?, expense_date = ?, updated_at = ?
		 WHERE id = ?`,
		expense.Amount, expense.TagID, expense.PaymentMethod, noteVal, expense.ExpenseDate, time.Now(), expense.ID,
	)
	return err
}

func (r *ExpenseRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM expenses WHERE id = ?`, id)
	return err
}

func (r *ExpenseRepository) List(familyID int64, filters models.ExpenseFilters) ([]models.Expense, int, error) {
	where := []string{"e.family_id = ?"}
	args := []interface{}{familyID}

	if filters.UserID != nil {
		where = append(where, "e.user_id = ?")
		args = append(args, *filters.UserID)
	}
	if filters.TagID != nil {
		where = append(where, "e.tag_id = ?")
		args = append(args, *filters.TagID)
	}
	if filters.PaymentMethod != nil {
		where = append(where, "e.payment_method = ?")
		args = append(args, *filters.PaymentMethod)
	}
	if filters.DateFrom != nil {
		where = append(where, "e.expense_date >= ?")
		args = append(args, *filters.DateFrom)
	}
	if filters.DateTo != nil {
		where = append(where, "e.expense_date <= ?")
		args = append(args, *filters.DateTo)
	}

	whereClause := strings.Join(where, " AND ")

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM expenses e WHERE %s", whereClause)
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf(
		`SELECT `+expenseSelectCols+expenseJoins+` WHERE %s ORDER BY e.expense_date DESC, e.created_at DESC LIMIT ? OFFSET ?`,
		whereClause,
	)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		e, err := scanExpense(rows)
		if err != nil {
			return nil, 0, err
		}
		expenses = append(expenses, *e)
	}

	return expenses, total, rows.Err()
}

func (r *ExpenseRepository) GetRecent(familyID int64, limit int) ([]models.Expense, error) {
	rows, err := r.db.Query(
		`SELECT `+expenseSelectCols+expenseJoins+`
		 WHERE e.family_id = ?
		 ORDER BY e.created_at DESC LIMIT ?`,
		familyID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		e, err := scanExpense(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *e)
	}
	return expenses, rows.Err()
}

func (r *ExpenseRepository) DB() *sql.DB {
	return r.db
}
