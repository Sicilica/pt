package sqlite3

import (
	"database/sql"
	"os"
	"path"
	"time"

	// Package for sqlite3 support
	_ "github.com/mattn/go-sqlite3"

	"github.com/sicilica/pt/types"
)

type sqlite3Provider struct {
	db *sql.DB
}

// sqlite3Provider implements StorageProvider.
var _ types.StorageProvider = (*sqlite3Provider)(nil)

// New returns a new SQLite3-powered storage provider. It automatically performs all
// necessary initialization logic to persist in a default location on the filesystem.
func New() (types.StorageProvider, error) {
	userDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := path.Join(userDir, "sicilica", "pt")
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path.Join(appDir, "pt.db"))
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		db.Close()
		return nil, err
	}

	return &sqlite3Provider{
		db: db,
	}, nil
}

func (p *sqlite3Provider) Close() error {
	return p.db.Close()
}

func (p *sqlite3Provider) NewSession() (types.StorageSession, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	return &sqlite3Session{
		tx: tx,
	}, nil
}

type sqlite3Session struct {
	tx *sql.Tx
}

// sqlite3Session implements StorageSession.
var _ types.StorageSession = (*sqlite3Session)(nil)

func (s *sqlite3Session) Abort() error {
	return s.tx.Rollback()
}

func (s *sqlite3Session) Commit() error {
	return s.tx.Commit()
}

func dbtime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.999")
}
