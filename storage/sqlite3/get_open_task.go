package sqlite3

import (
	"database/sql"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) GetOpenTask() (*types.Task, error) {
	t := &types.Task{}
	err := s.tx.QueryRow("SELECT id, start FROM tasks WHERE stop IS NULL LIMIT 1").Scan(&t.ID, &t.Start)
	if err == sql.ErrNoRows {
		return nil, types.ErrNoOpenTask
	} else if err != nil {
		return nil, err
	}

	return t, nil
}
