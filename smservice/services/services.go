package services

import (
	"database/sql"

	"github.com/TechCatsLab/sor/smservice/config"
	"github.com/TechCatsLab/sor/smservice/model/mysql"
)

type SMservice struct {
	Db   *sql.DB
	Conf *config.Config
}

func NewService(c *config.Config, db *sql.DB) *SMservice {
	sm := &SMservice{
		Db: db,
		Conf: &config.Config{
			Host:           c.Host,
			Appcode:        c.Appcode,
			Digits:         c.Digits,
			ResendInterval: c.ResendInterval,
			OnCheck:        c.OnCheck,
			DB:             c.DB,
		},
	}
	return sm
}

func (sm *SMservice) CreateDB() error {
	return mysql.CreateDatabase(sm.Db)
}

func (sm *SMservice) CreateTable() error {
	return mysql.CreateTable(sm.Db)
}

func (sm *SMservice) Insert(mobile string, date int64, code string, sign string) error {
	return mysql.Insert(sm.Db, mobile, date, code, sign)
}

func (sm *SMservice) GetCode(sign string) (string, error) {
	return mysql.GetCode(sm.Db, sign)
}

func (sm *SMservice) GetMobile(sign string) (string, error) {
	return mysql.GetMobile(sm.Db, sign)
}

func (sm *SMservice) GetDate(sign string) (int64, error) {
	return mysql.GetDate(sm.Db, sign)
}

func (sm *SMservice) GetMsg(sign string) *mysql.Message {
	return mysql.GetMessage(sm.Db, sign)
}
