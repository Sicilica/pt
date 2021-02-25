package sqlite3

import (
	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) RemoveTagFromTask(t *types.Task, tag string) error {
	stmt, err := s.tx.Prepare("DELETE FROM task_tags WHERE task_id=? AND tag=?")
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
