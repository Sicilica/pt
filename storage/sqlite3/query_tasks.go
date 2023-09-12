package sqlite3

import (
	"database/sql"
	"time"

	"github.com/sicilica/pt/types"
)

func (s sqlite3Session) QueryTasks(start time.Time, end time.Time, q *types.Query) ([]*types.Task, error) {
	var rows *sql.Rows

	if q.Tag == "" {
		stmt, err := s.tx.Prepare("SELECT tasks.id, start, stop FROM tasks WHERE stop IS NOT NULL AND NOT (stop<? OR start>?) ORDER BY start ASC")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(dbtime(start), dbtime(end))
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err := s.tx.Prepare(`
			WITH tags AS (
				SELECT tag FROM tag_parents_lookup WHERE ancestor=?
				UNION ALL
				SELECT ? tag
			)
			SELECT DISTINCT tasks.id, start, stop
			FROM tasks
			INNER JOIN task_tags ON tasks.id=task_tags.task_id
			INNER JOIN tags ON task_tags.tag = tags.tag
			WHERE
				stop IS NOT NULL AND
				NOT (stop<? OR start>?)
			ORDER BY start ASC
		`)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(q.Tag, q.Tag, start, end)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var tasks []*types.Task
	for rows.Next() {
		task := &types.Task{}
		if err := rows.Scan(&task.ID, &task.Start, &task.Stop); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
