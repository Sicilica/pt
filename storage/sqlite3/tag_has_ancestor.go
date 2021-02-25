package sqlite3

func (s sqlite3Session) TagHasAncestor(child, parent string) (bool, error) {
	stmt, err := s.tx.Prepare("WITH tag_tree AS (SELECT * FROM tag_parents WHERE parent=? UNION ALL SELECT a.* FROM tag_parents a JOIN tag_tree ON a.parent=tag_tree.tag) SELECT * FROM tag_tree WHERE tag=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(parent, child)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
