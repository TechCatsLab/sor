package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TechCatsLab/sor/order/model/mysql"

	"github.com/TechCatsLab/sor/order/config"
)

type OrderService struct {
	db   *sql.DB
	SQLS []string
	Cnf  *config.Config
}

const (
	orderDB = iota
	orderTable
	itemTable
	orderInsert
	itemInsert
	orderIdByOrderCode
	orderByOrderID
	itemsByOrderID
	orderListByUserID
	payByOrderID
	consignByOrderID
	statusByOrderID
)

func NewOrderService(c *config.Config, db *sql.DB) *OrderService {
	var (
		ostore = c.OrderDB + "." + c.OrderTable
		istore = c.OrderDB + "." + c.ItemTable
	)

	os := &OrderService{
		db: db,
		SQLS: []string{
			`CREATE DATABASE IF NOT EXISTS ` + c.OrderDB,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.OrderTable + `(
				id INT UNSIGNED NOT NULL AUTO_INCREMENT ,
				orderCode VARCHAR(50) NOT NULL,
				userID BIGINT UNSIGNED NOT NULL,
				shipCode VARCHAR(50) NOT NULL DEFAULT '100000',
				addressID VARCHAR(20) NOT NULL,
				totalPrice INT UNSIGNED NOT NULL,
				payWay TINYINT UNSIGNED DEFAULT '0',
				promotion TINYINT(1) UNSIGNED DEFAULT '0',
				freight INT UNSIGNED NOT NULL,
				status TINYINT UNSIGNED DEFAULT '0',
				created DATETIME DEFAULT NOW(),
				closed DATETIME DEFAULT '8012-12-31 00:00:00',
				updated DATETIME DEFAULT NOW(),
				PRIMARY KEY (id),
				UNIQUE KEY orderCode (orderCode) USING BTREE,
				KEY created (created),
				KEY updated (updated),
				KEY status (status), 
				KEY payWay (payWay)
			)ENGINE=InnoDB AUTO_INCREMENT = 10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='order info'`,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.ItemTable + `(
				productID INT UNSIGNED NOT NULL,
				orderID VARCHAR(50) NOT NULL,
				count INT UNSIGNED NOT NULL,
				price INT UNSIGNED NOT NULL,
				discount TINYINT UNSIGNED NOT NULL,
				KEY orderID (orderID)
			)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='orderitem info'`,
			`INSERT INTO ` + ostore + ` (orderCode,userID,addressID,totalPrice,promotion,freight,closed) VALUES(?,?,?,?,?,?,?)`,
			`INSERT INTO ` + istore + ` (productID,orderID,count,price,discount) VALUES(?,?,?,?,?)`,
			`SELECT id FROM ` + ostore + ` WHERE orderCode = ? LOCK IN SHARE MODE`,
			`SELECT * FROM ` + ostore + ` WHERE id = ? LOCK IN SHARE MODE`,
			`SELECT * FROM ` + istore + ` WHERE orderID = ? LOCK IN SHARE MODE`,
			`SELECT * FROM ` + ostore + ` WHERE userID = ? AND status = ? LOCK IN SHARE MODE`,
			`UPDATE ` + ostore + ` SET payWay = ? , updated = ? , status = 2 WHERE id = ? LIMIT 1 `,
			`UPDATE ` + ostore + ` SET shipCode = ? , updated = ? , status = 3 WHERE id = ? LIMIT 1 `,
			`UPDATE ` + ostore + ` SET status = ? , updated = ? WHERE id = ? LIMIT 1 `,
		},
		Cnf: c,
	}
	return os
}

func (os *OrderService) CreateDB() error {
	return mysql.CreateDB(os.db, os.SQLS[orderDB])
}

func (os *OrderService) CreateOrderTable() error {
	return mysql.CreateTable(os.db, os.SQLS[orderTable])
}

func (os *OrderService) CreateItemTable() error {
	return mysql.CreateTable(os.db, os.SQLS[itemTable])
}

func (os *OrderService) Insert(order mysql.Order, items []mysql.Item) (id uint32, err error) {
	tx, err := os.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	order.Closed = order.Created.Add(time.Duration(os.Cnf.ClosedInterval * int(time.Hour)))

	result, err := tx.Exec(os.SQLS[orderInsert], order.OrderCode, order.UserID, order.AddressID, order.TotalPrice, order.Promotion, order.Freight, order.Closed)

	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[insert order] : insert order affected 0 rows")
	}

	Id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	order.ID = uint32(Id)

	for _, x := range items {
		result, err := tx.Exec(os.SQLS[itemInsert], x.ProductId, order.ID, x.Count, x.Price, x.Discount)
		if err != nil {
			return 0, err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return 0, errors.New("insert item: insert affected 0 rows")
		}

		err = os.Cnf.User.UserCheck(tx, order.UserID, x.ProductId)
		if err != nil {
			return 0, err
		}

		err = os.Cnf.Stock.ModifyProductStock(tx, x.ProductId, int(x.Count))
		if err != nil {
			return 0, err
		}
	}

	return order.ID, err
}

func (os *OrderService) OrderIDByOrderCode(ordercode string) (uint32, error) {
	return mysql.OrderIDByOrderCode(os.db, os.SQLS[orderIdByOrderCode], ordercode)
}

func (os *OrderService) OrderInfoByOrderKey(orderid uint32) (*mysql.OrmOrder, error) {
	order, err := mysql.SelectByOrderKey(os.db, os.SQLS[orderByOrderID], os.SQLS[itemsByOrderID], orderid)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (os *OrderService) LisitOrderByUserId(userid uint64, status uint8) ([]*mysql.OrmOrder, error) {
	orders, err := mysql.LisitOrderByUserId(os.db, os.SQLS[orderListByUserID], os.SQLS[itemsByOrderID], userid, status)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *OrderService) UpdatePayByOrderKey(tx *sql.Tx, orderid uint32, payway uint8, time time.Time) (uint32, error) {
	result, err := tx.Exec(os.SQLS[payByOrderID], payway, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] ; not update pay infomation for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func (os *OrderService) UpdateShipByOrderKey(tx *sql.Tx, orderid uint32, shipcode string, time time.Time) (uint32, error) {
	result, err := tx.Exec(os.SQLS[consignByOrderID], shipcode, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not update ship infomation for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func (os *OrderService) UpdateStatusByOrderKey(tx *sql.Tx, orderid uint32, status uint8, time time.Time) (uint32, error) {
	result, err := tx.Exec(os.SQLS[statusByOrderID], status, time, orderid)
	if err != nil {
		return 0, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not update status  for order module ")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func (os *OrderService) CheckPromotion(tx *sql.Tx, orderid uint32) ([]*mysql.Item, error) {
	order, err := mysql.SelectByOrderKey(os.db, os.SQLS[orderByOrderID], os.SQLS[itemsByOrderID], orderid)
	if err != nil {
		return nil, err
	}
	if order.Promotion {
		return order.Orm, nil
	}

	return nil, nil
}
