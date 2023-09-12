package sqlite3

func (s sqlite3Session) SetTagParent(child, parent string) error {
	stmt, err := s.tx.Prepare("INSERT INTO tag_parents (tag, parent) VALUES (?,?) ON CONFLICT(tag) DO UPDATE SET parent=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(child, parent, parent)
	if err != nil {
		return err
	}

	return s.rebuildTagParentsLookup()
}

// NOTE: This already will support one tag having multiple parents.
// If we ever want to support that in the future, we only need to update the
// tag_parents schema and SetTagParent (since it relies on key collisions).
// However, having multiple parents would make it much harder to draw a clear
// hierarchical relationship for reporting, so it's not clear this is even a
// desirable feature.
func (s sqlite3Session) rebuildTagParentsLookup() error {
	// NOTE: sqlite doesn't have an explicit TRUNCATE TABLE statement
	_, err := s.tx.Exec("DELETE FROM tag_parents_lookup")
	if err != nil {
		return err
	}

	_, err = s.tx.Exec("INSERT INTO tag_parents_lookup SELECT tag, parent as ancestor, 0 as depth FROM tag_parents")
	if err != nil {
		return err
	}

	stmt, err := s.tx.Prepare(`
	INSERT INTO tag_parents_lookup
	SELECT A.tag, B.parent AS ancestor, A.depth + 1
	FROM tag_parents_lookup A
	INNER JOIN tag_parents B ON A.ancestor = B.tag
	WHERE A.depth = ?
	ON CONFLICT (tag, ancestor) DO NOTHING
`)
	if err != nil {
		return err
	}

	const MAX_DEPTH = 20
	for depth := 0; depth < MAX_DEPTH; depth++ {
		res, err := stmt.Exec(depth)
		if err != nil {
			return err
		}
		cnt, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if cnt == 0 {
			break
		}
	}

	return nil
}
