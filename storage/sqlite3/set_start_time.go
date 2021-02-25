package sqlite3

import (
	"time"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) SetStartTime(t *types.Task, start time.Time) error {
	stmt, err := s.tx.Prepare("UPDATE tasks SET start=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(start, t.ID)
	if err != nil {
		return err
	}

	t.Start = start

	return nil
}
