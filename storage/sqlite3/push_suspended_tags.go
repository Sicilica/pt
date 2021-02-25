package sqlite3

import (
	"strings"

	"github.com/pkg/errors"
)

const tagBundleSep = "|"

func (s sqlite3Session) PushSuspendedTags(tags []string) error {
	for _, tag := range tags {
		if strings.Contains(tag, tagBundleSep) {
			return errors.New("invalid character in tag")
		}
	}

	stmt, err := s.tx.Prepare("INSERT INTO suspension_records (tags) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(strings.Join(tags, tagBundleSep))
	if err != nil {
		return err
	}

	return nil
}
