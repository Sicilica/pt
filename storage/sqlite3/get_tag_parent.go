package sqlite3

import (
	"database/sql"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) GetTagParent(tag string) (string, error) {
	stmt, err := s.tx.Prepare("SELECT parent FROM tag_parents WHERE tag=? LIMIT 1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var parent string
	err = stmt.QueryRow(tag).Scan(&parent)
	if err == sql.ErrNoRows {
		return "", types.ErrNoParent
	} else if err != nil {
		return "", err
	}

	return parent, nil
}
