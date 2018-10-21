package mysql

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Mobile string `db:"mobile"`
	Date   int64  `db:"date"`
	Code   string `db:"code"`
	Sign   string `db:"sign"`
}

var sqlString = []string{
	`CREATE DATABASE IF NOT EXISTS sms`,
	`CREATE TABLE IF NOT EXISTS sms.msg(
		mobile VARCHAR(32) UNIQUE NOT NULL,
		date  INT(11) DEFAULT 0,
		code VARCHAR(32) ,
		sign VARCHAR(32) UNIQUE NOT NULL
	)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
	`INSERT INTO sms.msg(mobile,date,code,sign) VALUES (?,?,?,?)`,
	`SELECT date FROM sms.msg WHERE sign = ?`,
	`DELETE FROM sms.msg WHERE sign = ? LIMIT 1`,
	`SELECT code FROM sms.msg WHERE sign = ?`,
	`SELECT mobile FROM sms.msg WHERE sign = ?`,
}

// CreateDatabase create database sms.
func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(sqlString[0])
	return err
}

// CreateTable create sms.msg table.
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(sqlString[1])
	return err
}

// Insert Insert a new sms.
func Insert(db *sql.DB, mobile string, date int64, code string, sign string) error {

	result, err := db.Exec(sqlString[2], mobile, date, code, sign)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("errInvalidInsert")
	}

	return nil
}

// GetDate return sms date(unixtime) and nil if no err,or (0,err).
func GetDate(db *sql.DB, sign string) (int64, error) {
	var unixtime int64

	err := db.QueryRow(sqlString[3], sign).Scan(&unixtime)
	if err != nil {
		return 0, errors.New("errQueryDate")
	}

	return unixtime, nil
}

// Clear delete a  msg.
func Delete(db *sql.DB, sign string) error {

	_, err := db.Exec(sqlString[4], sign)
	if err != nil {
		return errors.New("errDeleteMysql")
	}

	return nil
}

// GetCode return sms date and nil or "0"and err.
func GetCode(db *sql.DB, sign string) (string, error) {
	var code string

	err := db.QueryRow(sqlString[5], sign).Scan(&code)
	if err != nil {
		return "0", errors.New("errQueryCode")
	}

	return code, nil
}

//GetMobile return User's mobile like ID or "0" and err
func GetMobile(db *sql.DB, sign string) (string, error) {
	var mobile string

	err := db.QueryRow(sqlString[6], mobile).Scan(&mobile)
	if err != nil {
		return "0", err
	}

	return mobile, nil
}

//GetMessage return msg
func GetMessage(db *sql.DB, sign string) *Message {
	var msg Message
	msg.Code, _ = GetCode(db, sign)
	msg.Date, _ = GetDate(db, sign)
	msg.Mobile, _ = GetMobile(db, sign)
	msg.Sign = sign

	return &msg
}
