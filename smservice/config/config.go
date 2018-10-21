package config

import (
	"database/sql"
)

type SMSVerify interface {
	OnVerifySucceed(targetID, mobile string)
	OnVerifyFailed(targetID, mobile string)
}

type Config struct {
	Host           string
	Appcode        string
	Digits         int
	ResendInterval int
	OnCheck        SMSVerify
	DB             *sql.DB
}
