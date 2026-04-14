package repository

import (
	"database/sql"
	"time"

	"expense-tracker/internal/models"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) GetAllForFamily(familyID int64) ([]models.Tag, error) {
	rows, err := r.db.Query(
		`SELECT id, name, icon, family_id, is_active, sort_order, created_at
		 FROM tags WHERE (family_id IS NULL OR family_id = ?) AND is_active = 1
		 ORDER BY sort_order, name`, familyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Icon, &t.FamilyID, &t.IsActive, &t.SortOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepository) GetSuggestedForUser(userID int64, limit int) ([]models.Tag, error) {
	rows, err := r.db.Query(
		`SELECT t.id, t.name, t.icon, t.family_id, t.is_active, t.sort_order, t.created_at
		 FROM tags t
		 JOIN user_tag_frequency utf ON t.id = utf.tag_id
		 WHERE utf.user_id = ? AND t.is_active = 1
		 ORDER BY utf.usage_count DESC
		 LIMIT ?`, userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Icon, &t.FamilyID, &t.IsActive, &t.SortOrder, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (r *TagRepository) Create(tag *models.Tag) error {
	result, err := r.db.Exec(
		`INSERT INTO tags (name, icon, family_id, is_active, sort_order, created_at)
		 VALUES (?, ?, ?, 1, ?, ?)`,
		tag.Name, tag.Icon, tag.FamilyID, tag.SortOrder, time.Now(),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = id
	return nil
}
