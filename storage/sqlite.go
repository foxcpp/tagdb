package storage

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

type S struct {
	DB *sql.DB

	// Prepared statements.
	addTag    *sql.Stmt
	remTag    *sql.Stmt
	forgetTag *sql.Stmt
	fileLst   *sql.Stmt
}

func Open(path string) (s *S, err error) {
	s = new(S)
	s.DB, err = sql.Open("sqlite3", "file:"+path)
	if err != nil {
		return nil, err
	}

	err = s.DB.Ping()
	if err != nil {
		return nil, err
	}

	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS map (
			tag TEXT NOT NULL,
			path TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = s.DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS tagpath_indx ON map (tag, path)")
	if err != nil {
		return nil, err
	}

	s.addTag, err = s.DB.Prepare("INSERT OR IGNORE INTO map VALUES ($2, $1)")
	if err != nil {
		return
	}
	s.remTag, err = s.DB.Prepare("DELETE FROM map WHERE path = $1 AND tag = $2")
	if err != nil {
		return nil, err
	}
	s.forgetTag, err = s.DB.Prepare("DELETE FROM map WHERE tag = $1")
	if err != nil {
		return nil, err
	}
	s.fileLst, err = s.DB.Prepare("SELECT path FROM map WHERE tag = $1")
	if err != nil {
		return nil, err
	}

	return
}

func (s *S) Close() {
	s.DB.Close()
}

func (s *S) FileList(tag string) (res []string, err error) {
	rows, err := s.fileLst.Query(tag)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		res = append(res, path)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *S) AddTag(tag, path string) error {
	_, err := s.addTag.Exec(tag, path)
	return err
}

func (s *S) RemoveTag(tag, path string) error {
	_, err := s.remTag.Exec(tag, path)
	return err
}

func (s *S) ForgetTag(tag string) error {
	_, err := s.forgetTag.Exec(tag)
	return err
}

func (s *S) CheckTag(tag, path string) (bool, error) {
	row := s.DB.QueryRow("SELECT path FROM map WHERE path = $1 AND tag = $2")
	var dummy string
	if err := row.Scan(&dummy); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
