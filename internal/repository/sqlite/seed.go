package sqlite

import "database/sql"

// seedCategory is one default category row. IDs are fixed so seeding stays
// idempotent and posts keep pointing at the same category across restarts.
type seedCategory struct {
	id    string
	title string
	icon  string
}

var defaultCategories = []seedCategory{
	{
		id:    "6f1c9a10-0000-4000-8000-000000000001",
		title: "Technology",
		icon:  `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>`,
	},
	{
		id:    "6f1c9a10-0000-4000-8000-000000000002",
		title: "Programming",
		icon:  `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>`,
	},
	{
		id:    "6f1c9a10-0000-4000-8000-000000000003",
		title: "Gaming",
		icon:  `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="6" y1="11" x2="10" y2="11"></line><line x1="8" y1="9" x2="8" y2="13"></line><line x1="15" y1="12" x2="15.01" y2="12"></line><line x1="18" y1="10" x2="18.01" y2="10"></line><path d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.545-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z"></path></svg>`,
	},
	{
		id:    "6f1c9a10-0000-4000-8000-000000000004",
		title: "News",
		icon:  `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-2 2Zm0 0a2 2 0 0 1-2-2v-9h2"></path><path d="M18 14h-8"></path><path d="M15 18h-5"></path><line x1="10" y1="6" x2="8" y2="6"></line><line x1="14" y1="6" x2="12" y2="6"></line><line x1="18" y1="6" x2="16" y2="6"></line><line x1="10" y1="10" x2="8" y2="10"></line><line x1="14" y1="10" x2="12" y2="10"></line><line x1="18" y1="10" x2="16" y2="10"></line></svg>`,
	},
}

const seedCategoryQuery = `INSERT OR IGNORE INTO categories (id, title, icon) VALUES(?, ?, ?)`

// seedCategories inserts the default categories. Existing rows are left as they
// are, so it is safe to run on every start.
func seedCategories(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, category := range defaultCategories {
		if _, err := tx.Exec(seedCategoryQuery, category.id, category.title, category.icon); err != nil {
			return err
		}
	}

	return tx.Commit()
}
