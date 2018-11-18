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

func NewOrderService(c *config.Config, db *sql.DB) *OrderService {
	store := c.OrderDB + "." + c.OrderTable
	os := &OrderService{
		db: db,
		SQLS: []string{
			`CREATE DATABASE IF NOT EXISTS ` + c.OrderDB,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.OrderTable + `(
				id INT UNSIGNED NOT NULL AUTO_INCREMENT ,
				orderCode VARCHAR(50) NOT NULL ,
				shipCode VARCHAR(50) NOT NULL DEFAULT 'wu',
				userID BIGINT UNSIGNED NOT NULL ,
				addressID VARCHAR(20) NOT NULL ,
				totalPrice INT UNSIGNED NOT NULL ,
				payWay TINYINT UNSIGNED DEFAULT '0' ,
				promotion INT UNSIGNED DEFAULT '0',
				freight INT UNSIGNED NOT NULL,
				status TINYINT UNSIGNED DEFAULT '0',
				created DATETIME DEFAULT NOW() ,
				closed DATETIME DEFAULT '8012-12-31 00:00:00' ,
				updated DATETIME DEFAULT NOW() ,
				PRIMARY KEY (id) ,
				UNIQUE KEY orderCode (orderCode) USING BTREE ,
				KEY created (created) ,
				KEY updated (updated) ,
				KEY status (status) ,
				KEY payWay (payWay)
			)ENGINE=InnoDB AUTO_INCREMENT = 10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='order info'`,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.ItemTable + `(
				id INT UNSIGNED NOT NULL AUTO_INCREMENT ,
				productID INT UNSIGNED NOT NULL ,
				orderID VARCHAR(50) NOT NULL ,
				count INT UNSIGNED NOT NULL ,
				price INT UNSIGNED NOT NULL ,
				discount TINYINT UNSIGNED NOT NULL ,
				size VARCHAR(50) NOT NULL ,
				color VARCHAR(50) NOT NULL ,
				PRIMARY KEY (id) ,
				KEY orderID (orderID)
			)ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='orderitem info'`,
			`INSERT INTO ` + store + ` (orderCode,userID,addressID,totalPrice,promotion,freight,closed) VALUES(?,?,?,?,?,?,?)`,
			`INSERT INTO ` + store +  ` (productID,orderID,count,price,discount,size,color) VALUES(?,?,?,?,?,?,?)`,
			`SELECT * FROM ` + store + ` WHERE userID = ? LOCK IN SHARE MODE`,
			`SELECT id FROM ` + store +  ` WHERE orderCode = ? LOCK IN SHARE MODE`,
			`SELECT * FROM `+ store +   `WHERE id = ? LOCK IN SHARE MODE`,
			`SELECT * FROM ` + store + `WHERE orderID = ? LOCK IN SHARE MODE`,
			`UPDATE ` + store + ` SET updated = ? WHERE id = ? LIMIT 1`,
			`UPDATE ` + store + ` SET status = ? WHERE id = ? LIMIT 1 `,
			`UPDATE `+ store + ` SET shipCode = ? WHERE id = ? LIMIT 1 `,
			`UPDATE ` + store +  ` SET payWay = ? WHERE id = ? LIMIT 1 `,
		},
		Cnf: c,
	}
	return os
}

func (os *OrderService) CreateDB() error {
	return mysql.CreateDB(os.db, os.SQLS[0])
}

func (os *OrderService) CreateOrderTable() error {
	return mysql.CreateTable(os.db, os.SQLS[1])
}

func (os *OrderService) CreateItemTable() error {
	return mysql.CreateTable(os.db, os.SQLS[2])
}

func (os *OrderService) Insert(order mysql.Order, items []mysql.Item) (uint32, error) {
	var (
		err error
	)

	err = os.Cnf.User.UserCheck(order.UserID)
	if err != nil {
		return 0, err
	}

	tx, err := os.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	order.Closed = order.Created.Add(time.Duration(os.Cnf.ClosedInterval * int(time.Hour)))

	result, err := tx.Exec(os.SQLS[3], order.OrderCode, order.UserID, order.AddressID, order.TotalPrice, order.Promotion, order.Freight, order.Closed)

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
		result, err := tx.Exec(os.SQLS[4], x.ProductId, order.ID, x.Count, x.Price, x.Discount, x.Size, x.Color)
		if err != nil {
			return 0, err
		}

		if affected, _ := result.RowsAffected(); affected == 0 {
			return 0, errors.New("insert item: insert affected 0 rows")
		}
	}
	for _, x := range items {
		err = os.Cnf.Stock.ModifyProductStock(x.ProductId, int(x.Count))
		if err != nil {
			return 0, err
		}
	}

	return order.ID, nil
}

func (os *OrderService) LisitOrderByUserId(userid uint64) ([]*mysql.Order, error) {
	orders, err := mysql.LisitOrderByUserId(os.db, os.SQLS[5], userid)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *OrderService) OrderIDByOrderCode(ordercode string) (uint32, error) {
	return mysql.OrderIDByOrderCode(os.db, os.SQLS[6], ordercode)
}

func (os *OrderService) OrderInfo(orderid uint32) (*mysql.Order, error) {
	order, err := mysql.SelectByOrderKey(os.db, os.SQLS[7], orderid)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (os *OrderService) LisitItemByOrderId(orderid uint32) ([]*mysql.Item, error) {
	items, err := mysql.LisitItemByOrderId(os.db, os.SQLS[8], orderid)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (os *OrderService) UpdateTime(orderid uint32, updated time.Time) (uint32, error) {
	return mysql.UpdateTimeByOrderKey(os.db, os.SQLS[9], orderid, updated)
}

func (os *OrderService) UpdateStatus(orderid uint32, status uint8) (uint32, error) {
	return mysql.UpdateStatusByOrderKey(os.db, os.SQLS[10], orderid, status)
}

func (os *OrderService) UpdateShip(orderid uint32, shippingcode string) (uint32, error) {
	return mysql.UpdateShipByOrderKey(os.db, os.SQLS[11], orderid, shippingcode)
}

func (os *OrderService) UpdatePayWay(orderid uint32, payWay uint8) (uint32, error) {
	return mysql.UpdatePayWayByOrderKey(os.db, os.SQLS[12], orderid, payWay)
}
