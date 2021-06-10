package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"namespace.com/dupfind/file"
)

const InsertStatement = "INSERT INTO files (filename, checksum) VALUES (?, ?)"
const UpsertStatement = InsertStatement + " ON CONFLICT(filename) DO UPDATE SET checksum=excluded.checksum"
const SelectStatement = "SELECT filename, checksum FROM files WHERE filename = ?"

func Create(db *sql.DB) {
	prepareStmtAndExec(db, "CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, filename TEXT, checksum TEXT)")
	prepareStmtAndExec(db, "CREATE UNIQUE INDEX IF NOT EXISTS idx_files_file ON files (filename)")
}

func prepareStmtAndExec(db *sql.DB, query string, args ...interface{}) sql.Result {
	statement, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	result, err := statement.Exec(args...)
	if err != nil {
		panic(err)
	}
	return result
}

func SelectOne(db *sql.DB, file *file.File) *file.File {
	// Use QueryRow because filename is unique, there will be at most one match
	row := db.QueryRow(SelectStatement, file.Name)
	var checksum string
	err := row.Scan(&file.Name, &checksum)
	if err == nil {
		if file.Checksum != nil {
			panic(fmt.Errorf("checksum for %s was already not nil (%s)", file.Name, *file.Checksum))
		}
		file.Checksum = &checksum // got a row, update the checksum
	} else if err != sql.ErrNoRows {
		panic(err)
	}
	return file
}

func Import(db *sql.DB, files []*file.File) bool {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(UpsertStatement)
	if err != nil {
		panic(err)
	}
	// This used to have issues, but no more? http://go-database-sql.org/prepared.html
	defer stmt.Close()

	for _, f := range files {
		_, err := stmt.Exec(f.Name, f.Checksum)
		if err != nil {
			panic(err)
		}
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return true
}
