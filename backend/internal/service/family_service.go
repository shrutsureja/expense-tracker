package service

import (
	"errors"
	"regexp"
	"time"

	"expense-tracker/internal/database"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var pinRegex = regexp.MustCompile(`^\d{6}$`)

func validatePIN(pin string) error {
	if !pinRegex.MatchString(pin) {
		return errors.New("PIN must be exactly 6 digits")
	}
	return nil
}

type FamilyService struct {
	familyRepo *repository.FamilyRepository
	userRepo   *repository.UserRepository
}

func NewFamilyService(familyRepo *repository.FamilyRepository, userRepo *repository.UserRepository) *FamilyService {
	return &FamilyService{
		familyRepo: familyRepo,
		userRepo:   userRepo,
	}
}

func (s *FamilyService) CreateFamily(name string, adminID int64, ownerUsername, ownerPassword, ownerDisplayName string) (*models.Family, *models.User, error) {
	if name == "" || ownerUsername == "" || ownerPassword == "" || ownerDisplayName == "" {
		return nil, nil, errors.New("all fields are required")
	}
	if err := validatePIN(ownerPassword); err != nil {
		return nil, nil, err
	}

	db := s.familyRepo.DB()
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Create family
	now := time.Now()
	result, err := tx.Exec(
		`INSERT INTO families (name, created_by, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		name, adminID, now, now,
	)
	if err != nil {
		return nil, nil, err
	}
	familyID, _ := result.LastInsertId()

	// Create owner user
	hash, err := bcrypt.GenerateFromPassword([]byte(ownerPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	result, err = tx.Exec(
		`INSERT INTO users (username, password_hash, display_name, role, family_id, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, 1, ?, ?)`,
		ownerUsername, string(hash), ownerDisplayName, models.RoleFamilyOwner, familyID, now, now,
	)
	if err != nil {
		return nil, nil, err
	}
	ownerID, _ := result.LastInsertId()

	// Link owner to family
	if _, err := tx.Exec(`UPDATE families SET owner_id = ? WHERE id = ?`, ownerID, familyID); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	// Seed default tags for the new family (best-effort, don't fail family creation)
	_ = database.SeedDefaultTagsForFamily(db, familyID)

	family := &models.Family{ID: familyID, Name: name, CreatedBy: adminID, OwnerID: &ownerID, CreatedAt: now, UpdatedAt: now}
	user := &models.User{ID: ownerID, Username: ownerUsername, DisplayName: ownerDisplayName, Role: models.RoleFamilyOwner, FamilyID: &familyID, IsActive: true, CreatedAt: now, UpdatedAt: now}

	return family, user, nil
}

func (s *FamilyService) ListFamilies() ([]models.FamilyWithOwner, error) {
	return s.familyRepo.GetAllWithOwners()
}

func (s *FamilyService) GetFamily(id int64) (*models.Family, error) {
	return s.familyRepo.GetByID(id)
}

func (s *FamilyService) DeleteFamily(id int64) error {
	return s.familyRepo.Delete(id)
}

func (s *FamilyService) AddMember(familyID int64, username, pin, displayName string) (*models.User, error) {
	if username == "" || pin == "" || displayName == "" {
		return nil, errors.New("username, pin, and display name are required")
	}
	if err := validatePIN(pin); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Role:         models.RoleFamilyMember,
		FamilyID:     &familyID,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *FamilyService) ListMembers(familyID int64) ([]models.User, error) {
	return s.userRepo.GetByFamilyID(familyID)
}

func (s *FamilyService) UpdateMember(memberID int64, displayName, pin string, familyID int64) error {
	user, err := s.userRepo.GetByID(memberID)
	if err != nil {
		return err
	}

	if user.FamilyID == nil || *user.FamilyID != familyID {
		return errors.New("member does not belong to this family")
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if pin != "" {
		if err := validatePIN(pin); err != nil {
			return err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hash)
	}

	return s.userRepo.Update(user)
}

func (s *FamilyService) DeactivateMember(memberID, familyID int64) error {
	user, err := s.userRepo.GetByID(memberID)
	if err != nil {
		return err
	}

	if user.FamilyID == nil || *user.FamilyID != familyID {
		return errors.New("member does not belong to this family")
	}

	if user.Role == models.RoleFamilyOwner {
		return errors.New("cannot deactivate the family owner")
	}

	return s.userRepo.Deactivate(memberID)
}
