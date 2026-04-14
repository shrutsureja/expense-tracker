package repository

import (
	"database/sql"
	"time"

	"expense-tracker/internal/models"
)

type FamilyRepository struct {
	db *sql.DB
}

func NewFamilyRepository(db *sql.DB) *FamilyRepository {
	return &FamilyRepository{db: db}
}

func (r *FamilyRepository) Create(family *models.Family) error {
	result, err := r.db.Exec(
		`INSERT INTO families (name, created_by, owner_id, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		family.Name, family.CreatedBy, family.OwnerID, time.Now(), time.Now(),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	family.ID = id
	return nil
}

func (r *FamilyRepository) GetAll() ([]models.Family, error) {
	rows, err := r.db.Query(
		`SELECT id, name, created_by, owner_id, created_at, updated_at FROM families ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families []models.Family
	for rows.Next() {
		var f models.Family
		if err := rows.Scan(&f.ID, &f.Name, &f.CreatedBy, &f.OwnerID, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		families = append(families, f)
	}
	return families, rows.Err()
}

func (r *FamilyRepository) GetByID(id int64) (*models.Family, error) {
	f := &models.Family{}
	err := r.db.QueryRow(
		`SELECT id, name, created_by, owner_id, created_at, updated_at FROM families WHERE id = ?`, id,
	).Scan(&f.ID, &f.Name, &f.CreatedBy, &f.OwnerID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *FamilyRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM families WHERE id = ?`, id)
	return err
}

func (r *FamilyRepository) UpdateOwner(familyID, ownerID int64) error {
	_, err := r.db.Exec(`UPDATE families SET owner_id = ?, updated_at = ? WHERE id = ?`, ownerID, time.Now(), familyID)
	return err
}

func (r *FamilyRepository) DB() *sql.DB {
	return r.db
}
