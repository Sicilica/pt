package sqlite3

import (
	"database/sql"
	"strings"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) PopSuspendedTags(offset int) ([]string, error) {
	var id string
	var tags string
	stmt, err := s.tx.Prepare("SELECT id, tags FROM suspension_records ORDER BY created_at DESC LIMIT 1 OFFSET ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(offset).Scan(&id, &tags)
	if err == sql.ErrNoRows {
		return nil, types.ErrNoSuspension
	} else if err != nil {
		return nil, err
	}

	stmt, err = s.tx.Prepare("DELETE FROM suspension_records WHERE id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return nil, err
	}

	return strings.Split(tags, tagBundleSep), nil
}
