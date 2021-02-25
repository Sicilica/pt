package sqlite3

import (
	"time"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) SetStopTime(t *types.Task, stop time.Time) error {
	stmt, err := s.tx.Prepare("UPDATE tasks SET stop=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(stop, t.ID)
	if err != nil {
		return err
	}

	t.Stop = stop

	return nil
}
