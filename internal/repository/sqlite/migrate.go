package sqlite

import (
	"database/sql"
	"fmt"
)

// migrate reconciles databases created before post_categories existed. The
// schema file only uses "IF NOT EXISTS", so it never rewrites an object that is
// already present — anything that changed shape has to be handled here.
func migrate(db *sql.DB) error {
	if err := fixCascadeDeleteTrigger(db); err != nil {
		return err
	}
	return dropLegacyPostCategory(db)
}

// fixCascadeDeleteTrigger replaces cascade_delete_poosts. Older databases hold
// a version referencing the non-existent table "reaction", which makes both
// post deletion and any ALTER TABLE on posts fail.
const recreateCascadeTrigger = `
DROP TRIGGER IF EXISTS cascade_delete_poosts;
CREATE TRIGGER cascade_delete_poosts
AFTER DELETE ON posts
BEGIN
	DELETE FROM comments WHERE old.id = parent_id;
	DELETE FROM reactions WHERE old.id = parent_id;
END;`

func fixCascadeDeleteTrigger(db *sql.DB) error {
	_, err := db.Exec(recreateCascadeTrigger)
	return err
}

// dropLegacyPostCategory moves the old single-category column into the
// post_categories join table and then removes it. Categories now live in
// post_categories only, so leaving the column would give two sources of truth.
func dropLegacyPostCategory(db *sql.DB) error {
	legacy, err := hasColumn(db, "posts", "category")
	if err != nil {
		return err
	}
	if !legacy {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Only carry over rows whose category still resolves, otherwise the
	// post_categories foreign key would reject the backfill.
	if _, err := tx.Exec(`
		INSERT OR IGNORE INTO post_categories (post_id, category_id)
		SELECT P.id, P.category
		FROM posts P
		JOIN categories C ON C.id = P.category
		WHERE P.category IS NOT NULL`); err != nil {
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE posts DROP COLUMN category`); err != nil {
		return fmt.Errorf("dropping legacy posts.category: %w", err)
	}

	return tx.Commit()
}

func hasColumn(db *sql.DB, table, column string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid, notNull, primaryKey int
			name, columnType         string
			defaultValue             sql.NullString
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}
