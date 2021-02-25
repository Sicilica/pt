package sqlite3

func (s sqlite3Session) ClearTagParent(tag string) error {
	stmt, err := s.tx.Prepare("DELETE FROM tag_parents WHERE tag=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tag)
	if err != nil {
		return err
	}

	return nil
}
