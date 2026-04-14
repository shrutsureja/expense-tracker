package service

import (
	"database/sql"
	"errors"
	"time"

	"expense-tracker/internal/auth"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errors.New("invalid username or password")
		}
		return "", nil, err
	}

	if !user.IsActive {
		return "", nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	token, err := auth.GenerateToken(user.ID, user.Username, user.Role, user.FamilyID, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
