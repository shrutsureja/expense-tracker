package service

import (
	"errors"

	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
)

type ExpenseService struct {
	expenseRepo *repository.ExpenseRepository
}

func NewExpenseService(expenseRepo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) AddExpense(familyID, userID int64, amount float64, tagID int64, paymentMethod, note, expenseDate string) (*models.Expense, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if expenseDate == "" {
		return nil, errors.New("expense date is required")
	}

	expense := &models.Expense{
		FamilyID:      familyID,
		UserID:        userID,
		Amount:        amount,
		TagID:         tagID,
		PaymentMethod: paymentMethod,
		Note:          note,
		ExpenseDate:   expenseDate,
	}

	if err := s.expenseRepo.Create(expense); err != nil {
		return nil, err
	}

	// Fetch the full expense with joined fields
	return s.expenseRepo.GetByID(expense.ID)
}

func (s *ExpenseService) UpdateExpense(expenseID int64, userID int64, userRole string, familyID int64, amount float64, tagID int64, paymentMethod, note, expenseDate string) error {
	existing, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return errors.New("expense not found")
	}

	if existing.FamilyID != familyID {
		return errors.New("expense does not belong to your family")
	}

	// Only creator or family owner can update
	if existing.UserID != userID && userRole != models.RoleFamilyOwner {
		return errors.New("only the creator or family owner can update this expense")
	}

	existing.Amount = amount
	existing.TagID = tagID
	existing.PaymentMethod = paymentMethod
	existing.Note = note
	existing.ExpenseDate = expenseDate

	return s.expenseRepo.Update(existing)
}

func (s *ExpenseService) DeleteExpense(expenseID, userID int64, userRole string, familyID int64) error {
	existing, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return errors.New("expense not found")
	}

	if existing.FamilyID != familyID {
		return errors.New("expense does not belong to your family")
	}

	if existing.UserID != userID && userRole != models.RoleFamilyOwner {
		return errors.New("only the creator or family owner can delete this expense")
	}

	return s.expenseRepo.Delete(expenseID)
}

func (s *ExpenseService) ListExpenses(familyID int64, filters models.ExpenseFilters) ([]models.Expense, int, error) {
	return s.expenseRepo.List(familyID, filters)
}

func (s *ExpenseService) GetRecentExpenses(familyID int64) ([]models.Expense, error) {
	return s.expenseRepo.GetRecent(familyID, 20)
}
