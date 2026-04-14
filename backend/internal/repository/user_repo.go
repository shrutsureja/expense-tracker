package repository

import (
	"database/sql"
	"time"

	"expense-tracker/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, display_name, role, family_id, is_active, created_at, updated_at
		 FROM users WHERE username = ?`, username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName,
		&user.Role, &user.FamilyID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, display_name, role, family_id, is_active, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.DisplayName,
		&user.Role, &user.FamilyID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	result, err := r.db.Exec(
		`INSERT INTO users (username, password_hash, display_name, role, family_id, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.PasswordHash, user.DisplayName, user.Role,
		user.FamilyID, user.IsActive, time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) GetByFamilyID(familyID int64) ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, password_hash, display_name, role, family_id, is_active, created_at, updated_at
		 FROM users WHERE family_id = ? ORDER BY created_at`, familyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.DisplayName,
			&u.Role, &u.FamilyID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET display_name = ?, password_hash = ?, is_active = ?, updated_at = ? WHERE id = ?`,
		user.DisplayName, user.PasswordHash, user.IsActive, time.Now(), user.ID,
	)
	return err
}

func (r *UserRepository) Deactivate(id int64) error {
	_, err := r.db.Exec(`UPDATE users SET is_active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	return err
}
