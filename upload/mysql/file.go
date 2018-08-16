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
	mysqlFileQueryByTime
	mysqlFileQueryByUserID
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

	sqlString = []string{
		`CREATE TABLE IF NOT EXISTS files (
			user_id 	INTEGER UNSIGNED NOT NULL,
			md5 		VARCHAR(512) NOT NULL DEFAULT ' ',
			path 		VARCHAR(512) NOT NULL DEFAULT ' ',
			created_at 	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (md5)
		) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;`,
		`INSERT INTO upload.files(user_id,md5,path,created_at) VALUES (?,?,?,?)`,
		`SELECT * FROM upload.files WHERE created_at < ?`,
		`SELECT md5,path,created_at FROM upload.files WHERE user_id = ?`,
		`SELECT path FROM upload.files WHERE md5 = ?`,
	}
)

// Create Files table.
func Create(db *sql.DB) error {
	_, err := db.Exec(sqlString[mysqlFileCreateTable])

	return err
}

// Insert a file
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

func QueryByTime(db *sql.DB, t string) (*[]files, error) {
	var (
		userID uint
		md5 string
		path string
		createdAt time.Time
		result []files
	)

	rows, err := db.Query(sqlString[mysqlFileQueryByTime], t)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userID, &md5, &path, &createdAt)
		if err != nil {
			return nil, err
		}

		content := files{
			UserID:    userID,
			Path:      path,
			MD5:       md5,
			CreatedAt: createdAt,
		}

		result = append(result, content)
	}

	return &result, nil
}

func QueryByMD5(db *sql.DB, md5 string) (string, error) {
	var (
		path string
	)

	err := db.QueryRow(sqlString[mysqlFileQueryByMD5], md5).Scan(&path)
	if err != nil {
		return path, err
	}

	return path, nil
}

func QueryByUserID(db *sql.DB, userID uint) (*[]files, error) {
	var (
		md5 string
		path string
		createdAt time.Time
		result []files
	)

	rows, err := db.Query(sqlString[mysqlFileQueryByUserID], userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&path, &md5, &createdAt)
		if err != nil {
			return nil, err
		}

		content := files{
			Path:      path,
			MD5:       md5,
			CreatedAt: createdAt,
		}

		result = append(result, content)
	}

	return &result, nil
}
