package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
)

var dbMigrations = [][]string{
	{
		// tasks
		`CREATE TABLE tasks (
	id INTEGER PRIMARY KEY,
	start DATETIME NOT NULL,
	stop DATETIME
	);`,
		"CREATE INDEX tasks_start_idx ON tasks (start ASC);",
		"CREATE INDEX tasks_stop_idx ON tasks (stop DESC);",

		// task_tags
		`CREATE TABLE task_tags (
task_id INTEGER NOT NULL,
tag TEXT NOT NULL,
PRIMARY KEY (task_id, tag)
);`,
		"CREATE INDEX task_tags_task_id_idx ON task_tags (task_id);",
		"CREATE INDEX task_tags_tag_idx ON task_tags (tag);",

		// suspension_records
		`CREATE TABLE suspension_records (
	id INTEGER PRIMARY KEY,
	tags TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`,
		"CREATE INDEX suspension_records_created_at ON suspension_records (created_at DESC);",

		// tag_parents
		`CREATE TABLE tag_parents (
	tag TEXT PRIMARY KEY,
	parent TEXT NOT NULL
	);`,
		"CREATE INDEX tag_parents_parent_idx ON tag_parents (parent);",

		// tag_parents_lookup
		`CREATE TABLE tag_parents_lookup (
	tag TEXT NOT NULL,
	ancestor TEXT NOT NULL,
	depth TINYINT NOT NULL,
	PRIMARY KEY (tag, ancestor)
	);`,
		"CREATE INDEX tag_parents_lookup_tag_idx ON tag_parents_lookup (tag);",
		"CREATE INDEX tag_parents_lookup_ancestor_idx ON tag_parents_lookup (ancestor);",
		"CREATE INDEX tag_parents_lookup_depth_idx ON tag_parents_lookup (depth);",
	},
}

func migrateDB(db *sql.DB) error {
	ok, err := migrationTableExists(db)
	if err != nil {
		return err
	}
	if !ok {
		if _, err := db.Exec("CREATE TABLE ptversion (version INTEGER)"); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO ptversion (version) VALUES (0)"); err != nil {
			return err
		}
	}

	v, err := getMigrationVer(db)
	if err != nil {
		return err
	}

	if v > len(dbMigrations) {
		return errors.New("db version is newer than pt")
	}
	if v == len(dbMigrations) {
		return nil
	}

	for i := v; i < len(dbMigrations); i++ {
		for _, m := range dbMigrations[i] {
			if _, err := db.Exec(m); err != nil {
				return err
			}
		}
	}

	if err := setMigrationVer(db, len(dbMigrations)); err != nil {
		return err
	}

	return nil
}

func getMigrationVer(db *sql.DB) (int, error) {
	row := db.QueryRow("SELECT version FROM ptversion")
	var v int
	if err := row.Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func setMigrationVer(db *sql.DB, v int) error {
	res, err := db.Exec(fmt.Sprintf("UPDATE ptversion SET version=%d", v))
	if err != nil {
		return err
	}
	cnt, err := res.RowsAffected()
	if cnt < 1 || err != nil {
		return errors.New("failed to write migration version")
	}
	return nil
}

func migrationTableExists(db *sql.DB) (bool, error) {
	rows, err := db.Query("SELECT 1 FROM sqlite_master WHERE type='table' AND name='ptversion'")
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
