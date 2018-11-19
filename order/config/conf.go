package config

import (
	"database/sql"
)

type Stocker interface {
	ModifyProductStock(tx *sql.Tx, targetID uint32, num int) error
}

type UserChecker interface {
	UserCheck(tx *sql.Tx, userid uint64, productID uint32) error
}

type Config struct {
	OrderDB        string
	OrderTable     string
	ItemTable      string
	ClosedInterval int

	Stock Stocker
	User  UserChecker
}
