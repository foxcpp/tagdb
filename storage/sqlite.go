package storage

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type S struct {
	DB *sql.DB

	// Prepared statements.
	tag         *sql.Stmt
	untag       *sql.Stmt
	delTag      *sql.Stmt
	renameTag   *sql.Stmt
	taggedFiles *sql.Stmt
	tagLst      *sql.Stmt
	tagsOnFile  *sql.Stmt
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

	initStmt := func(target **sql.Stmt, stmt string) {
		var err error
		*target, err = s.DB.Prepare(stmt)
		if err != nil {
			panic(err)
		}
	}
	initStmt(&s.tag, "INSERT OR IGNORE INTO map VALUES ($2, $1)")
	initStmt(&s.untag, "DELETE FROM map WHERE path = $1 AND tag = $2")
	initStmt(&s.delTag, "DELETE FROM map WHERE tag = $1")
	initStmt(&s.renameTag, "UPDATE OR REPLACE map SET tag = $2 WHERE tag = $1")
	initStmt(&s.taggedFiles, "SELECT path FROM map WHERE tag = $1")
	initStmt(&s.tagLst, "SELECT tag FROM map GROUP BY tag")
	initStmt(&s.tagsOnFile, "SELECT tag FROM map WHERE path = $1 GROUP BY tag")

	return
}

func (s *S) Close() {
	s.DB.Close()
}

func (s *S) Begin() (*sql.Tx, error) {
	return s.DB.Begin()
}

func (s *S) Tags() (res []string, err error) {
	rows, err := s.tagLst.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		res = append(res, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *S) TagsOnFile(path string) (res []string, err error) {
	if !IsValidPath(path) {
		return nil, ErrInvalidPath
	}
	rows, err := s.tagsOnFile.Query(path)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		res = append(res, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *S) RenameTag(tx *sql.Tx, oldname, newname string) error {
	if !IsValidTag(oldname) {
		return ErrInvalidTag
	}
	if !IsValidTag(newname) {
		return ErrInvalidTag
	}
	// FIXME: Cryptic error forces us to flip arguments in the statement. Further
	// investigation is required.
	_, err := tx.Stmt(s.renameTag).Exec(newname, oldname)
	return err
}

func (s *S) TaggedFiles(tag string) (res []string, err error) {
	if !IsValidTag(tag) {
		return nil, ErrInvalidTag
	}
	rows, err := s.taggedFiles.Query(tag)
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

func (s *S) Tag(tx *sql.Tx, tag, path string) error {
	if !IsValidTag(tag) {
		return ErrInvalidTag
	}
	if !IsValidPath(path) {
		return ErrInvalidPath
	}
	_, err := tx.Stmt(s.tag).Exec(tag, path)
	return err
}

func (s *S) Untag(tx *sql.Tx, tag, path string) error {
	if !IsValidTag(tag) {
		return ErrInvalidTag
	}
	if !IsValidPath(path) {
		return ErrInvalidPath
	}
	_, err := tx.Stmt(s.untag).Exec(tag, path)
	return err
}

func (s *S) DeleteTag(tx *sql.Tx, tag string) error {
	if !IsValidTag(tag) {
		return ErrInvalidTag
	}
	_, err := tx.Stmt(s.delTag).Exec(tag)
	return err
}

func (s *S) CheckTag(tag, path string) (bool, error) {
	if !IsValidTag(tag) {
		return false, ErrInvalidTag
	}
	if !IsValidPath(path) {
		return false, ErrInvalidPath
	}
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

// FileListAllTags is similar to FileList but returns only entries with all
// listed tags attached.
func (s *S) FileListAllTags(tags ...string) (res []string, err error) {
	// Prepare WHERE condition part.
	// We need a `"tag", "tag"` string.
	tagsQuoted := make([]string, len(tags))
	for i, tag := range tags {
		tagsQuoted[i] = `"` + tag + `"`
	}
	tagsBlock := strings.Join(tagsQuoted, ", ")
	stmt := `
		SELECT path FROM map
		WHERE tag IN (` + tagsBlock + `)
		GROUP BY path
		HAVING COUNT(path) = ` + strconv.Itoa(len(tags))
	rows, err := s.DB.Query(stmt)
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

// FileListAnyTag is similar to FileList but returns entries if it does have
// at least one listed tag attached.
func (s *S) FileListAnyTag(tags ...string) (res []string, err error) {
	// Prepare WHERE condition part.
	// We need a `"tag", "tag"` string.
	tagsQuoted := make([]string, len(tags))
	for i, tag := range tags {
		tagsQuoted[i] = `"` + tag + `"`
	}
	tagsBlock := strings.Join(tagsQuoted, ", ")
	stmt := `
		SELECT path FROM map
		WHERE tag IN (` + tagsBlock + `)
		GROUP BY path`
	rows, err := s.DB.Query(stmt)
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
