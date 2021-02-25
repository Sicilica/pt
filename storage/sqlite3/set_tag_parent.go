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

	return nil
}
