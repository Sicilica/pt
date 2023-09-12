package sqlite3

import (
	"fmt"
	"time"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) CreateTask(start time.Time) (*types.Task, error) {
	stmt, err := s.tx.Prepare("INSERT INTO tasks (start) VALUES (?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	r, err := stmt.Exec(dbtime(start))
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.Task{
		ID:    fmt.Sprint(id),
		Start: start,
	}, nil
}
