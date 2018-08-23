/*
 * Revision History:
 *     Initial: 2018/08/10        Shi Ruitao
 */

package mysql

import (
	"database/sql"
	"errors"
	"time"
)

const (
	mysqlFileCreateTable = iota
	mysqlFileInsert
	mysqlFileQueryByMD5
)

type (
	files struct {
		UserID    uint
		Path      string
		MD5       string
		CreatedAt time.Time
	}
)

var (
	errInvalidInsert = errors.New("upload file: insert affected 0 rows")
	ErrNoRows = errors.New("there is no such data in database")

	sqlString = []string{
		`CREATE TABLE IF NOT EXISTS files (
			user_id 	INTEGER UNSIGNED NOT NULL,
			md5 		VARCHAR(512) NOT NULL DEFAULT ' ',
			path 		VARCHAR(512) NOT NULL DEFAULT ' ',
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (md5)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;`,
		`INSERT INTO upload.files(user_id,md5,path,created_at) VALUES (?,?,?,?)`,
		`SELECT path FROM upload.files WHERE md5 = ?`,
	}
)

// Create create files table.
func Create(db *sql.DB) error {
	_, err := db.Exec(sqlString[mysqlFileCreateTable])

	return err
}

// Insert insert a file
func Insert(db *sql.DB, userID uint, path, md5 string) error {
	result, err := db.Exec(sqlString[mysqlFileInsert], userID, md5, path, time.Now())
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

// QueryByMD5 select by MD5
func QueryByMD5(db *sql.DB, md5 string) (string, error) {
	var (
		path string
	)

	err := db.QueryRow(sqlString[mysqlFileQueryByMD5], md5).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			return path, ErrNoRows
		}
		return path, err
	}

	return path, nil
}
