package database

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin(db *sql.DB, username, password string) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists > 0 {
		log.Printf("Super admin '%s' already exists, skipping seed", username)
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password_hash, display_name, role, family_id, is_active)
		 VALUES (?, ?, ?, 'super_admin', NULL, 1)`,
		username, string(hash), "Super Admin",
	)
	if err != nil {
		return err
	}

	log.Printf("Super admin '%s' created successfully", username)
	return nil
}

func SeedDefaultTags(db *sql.DB) error {
	tags := []struct {
		Name      string
		Icon      string
		SortOrder int
	}{
		{"Groceries", "🛒", 1},
		{"Milk", "🥛", 2},
		{"Vegetables", "🥬", 3},
		{"Fruits", "🍎", 4},
		{"Petrol/Fuel", "⛽", 5},
		{"Electricity Bill", "💡", 6},
		{"Water Bill", "💧", 7},
		{"Gas Bill", "🔥", 8},
		{"Internet/WiFi", "📶", 9},
		{"Mobile Recharge", "📱", 10},
		{"Rent", "🏠", 11},
		{"Medicine", "💊", 12},
		{"Doctor", "🏥", 13},
		{"Clothing", "👕", 14},
		{"Education", "📚", 15},
		{"Dining Out", "🍽️", 16},
		{"Snacks", "🍿", 17},
		{"Tea/Coffee", "☕", 18},
		{"Auto/Cab", "🛺", 19},
		{"Bus/Train", "🚌", 20},
		{"Amazon/Flipkart", "📦", 21},
		{"Subscription", "📺", 22},
		{"House Maintenance", "🔧", 23},
		{"Kitchen Items", "🍳", 24},
		{"Personal Care", "🧴", 25},
		{"Gifts", "🎁", 26},
		{"Donations", "🙏", 27},
		{"EMI/Loan", "🏦", 28},
		{"Insurance", "🛡️", 29},
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tags WHERE family_id IS NULL").Scan(&count)
	if err != nil {
		return err
	}
	if count >= len(tags) {
		log.Println("Default tags already seeded, skipping")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO tags (name, icon, family_id, is_active, sort_order) VALUES (?, ?, NULL, 1, ?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range tags {
		if _, err := stmt.Exec(t.Name, t.Icon, t.SortOrder); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("Seeded %d default tags", len(tags))
	return nil
}
