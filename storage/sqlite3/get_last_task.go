package sqlite3

import (
	"database/sql"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) GetLastTask() (*types.Task, error) {
	t := &types.Task{}
	err := s.tx.QueryRow("SELECT id, start, stop FROM tasks ORDER BY stop DESC LIMIT 1").Scan(&t.ID, &t.Start, &t.Stop)
	if err == sql.ErrNoRows {
		// This case should never matter, so return an empty task
		return t, nil
	} else if err != nil {
		return nil, err
	}

	return t, nil
}
