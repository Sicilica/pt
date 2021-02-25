package sqlite3

import (
	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) GetTaskTags(t *types.Task) ([]string, error) {
	stmt, err := s.tx.Prepare("SELECT tag FROM task_tags WHERE task_id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(t.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
