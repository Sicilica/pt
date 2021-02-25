package sqlite3

import (
	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) AddTagToTask(t *types.Task, tag string) error {
	stmt, err := s.tx.Prepare("INSERT INTO task_tags (task_id, tag) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(t.ID, tag)
	if err != nil {
		return err
	}

	return nil
}
