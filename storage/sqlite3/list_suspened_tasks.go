package sqlite3

import (
	"strings"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) ListSuspendedTasks() ([]*types.SuspensionRecord, error) {
	rows, err := s.tx.Query("SELECT tags, created_at FROM suspension_records ORDER BY created_at DESC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suspensions []*types.SuspensionRecord
	for rows.Next() {
		sr := &types.SuspensionRecord{}
		var tags string
		err := rows.Scan(&tags, &sr.Time)
		if err != nil {
			return nil, err
		}
		sr.Tags = strings.Split(tags, tagBundleSep)
		suspensions = append(suspensions, sr)
	}

	return suspensions, nil
}
