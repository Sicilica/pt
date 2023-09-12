package sqlite3

func (s sqlite3Session) TagHasAncestor(child, ancestor string) (bool, error) {
	stmt, err := s.tx.Prepare("SELECT 1 FROM tag_parents_lookup WHERE tag=? AND ancestor=?")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(child, ancestor)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
