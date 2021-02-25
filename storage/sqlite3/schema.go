package sqlite3

import (
	"database/sql"
)

var dbMigrations = []string{
	// v1
	"CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY, start DATETIME DEFAULT CURRENT_TIMESTAMP, stop DATETIME)",
	"CREATE TABLE IF NOT EXISTS task_tags (task_id INTEGER, tag TEXT)",
	"CREATE INDEX IF NOT EXISTS idx_task_tags_task_id ON task_tags(task_id)",
	"CREATE INDEX IF NOT EXISTS idx_task_tags_tag ON task_tags(tag)",
	"CREATE UNIQUE INDEX IF NOT EXISTS idx_task_tags_unique ON task_tags(task_id, tag)",
	"CREATE TABLE IF NOT EXISTS suspension_records (id INTEGER PRIMARY KEY, tags TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)",
	"CREATE TABLE IF NOT EXISTS tag_parents (tag TEXT PRIMARY KEY, parent TEXT)",
}

func migrateDB(db *sql.DB) error {
	for _, m := range dbMigrations {
		if _, err := db.Exec(m); err != nil {
			return err
		}
	}
	return nil
}
