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
	os := &OrderService{
		db: db,
		SQLS: []string{
			`CREATE DATABASE IF NOT EXISTS ` + c.OrderDB,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.OrderTable + `(
				orderId VARCHAR(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '订单id', 
				payment VARCHAR(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '实付金额。精确到2位小数;单位:元。如:200.07，表示:200元7分',
				paymentType INT(2) DEFAULT NULL COMMENT '支付平台',
				promotion VARCHAR(50) DEFAULT NULL COMMENT '促销',
				postFee VARCHAR(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '邮费。精确到2位小数;单位:元。如:200.07，表示:200元7分',  
				status INT(10) DEFAULT '1' COMMENT '状态：1、未付款，2、已付款，3、未发货，4、已发货，5、交易成功，6、交易关闭',
				createTime DATETIME DEFAULT NULL COMMENT '订单创建时间',
				payTime DATETIME DEFAULT NULL COMMENT '付款时间',
				consignTime DATETIME DEFAULT NULL COMMENT '发货时间',
				closeTime DATETIME DEFAULT NULL COMMENT '交易关闭时间',
				endTime DATETIME DEFAULT NULL COMMENT '交易完成时间',
				shippingName VARCHAR(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '物流名称',
				shippingCode VARCHAR(20) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '物流单号',
				userId BIGINT(20) DEFAULT NULL COMMENT '用户id',
				UNIQUE KEY orderId (orderId) USING BTREE,
				KEY createTime (createTime),
				KEY status (status),
				KEY paymentType (paymentType)
			)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='订单基本信息表'`,
			`CREATE TABLE IF NOT EXISTS ` + c.OrderDB + `.` + c.ItemTable + `(
				itemId VARCHAR(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '商品id',
				orderId VARCHAR(50) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '订单id',
				num INT(10) DEFAULT NULL COMMENT '商品购买数量',
				price VARCHAR(50) DEFAULT NULL COMMENT '商品单价',
				total VARCHAR(50) DEFAULT NULL COMMENT '商品总金额',
				KEY orderId (orderId)
			)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='订单商品表'`,
			`INSERT INTO ` + c.OrderDB + `.` + c.OrderTable + ` (orderId,payment,promotion,postFee,createTime,endTime,userId) VALUES(?,?,?,?,?,?,?)`,
			`INSERT INTO ` + c.OrderDB + `.` + c.ItemTable + ` (itemId,orderId,num,price,total) VALUES(?,?,?,?,?)`,
			`SELECT * FROM ` + c.OrderDB + `.` + c.OrderTable + ` WHERE orderId = ? LIMIT 1 `,
			`SELECT * FROM ` + c.OrderDB + `.` + c.OrderTable + ` WHERE orderId = ? `,
			`SELECT * FROM ` + c.OrderDB + `.` + c.ItemTable + ` WHERE orderId = ?  `,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET payTime = ? WHERE orderId = ? `,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET consignTime = ? WHERE orderId = ? `,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET endTime = ? WHERE orderId = ? `,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET status = ? WHERE orderId = ? `,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET shippingName = ? , shippingCode = ? WHERE orderId = ? `,
			`DELETE FROM ` + c.OrderDB + `.` + c.OrderTable + ` WHERE orderId = ? LIMIT 1 LOCK IN SHARE MODE`,
			`DELETE FROM ` + c.OrderDB + `.` + c.ItemTable + ` WHERE orderId = ? LOCK IN SHARE MODE`,
			`UPDATE ` + c.OrderDB + `.` + c.OrderTable + ` SET paymentType = ? WHERE orderId = ? `,
		},
		Cnf: c,
	}
	return os
}

type Item struct {
	ItemId string `json:"itemid"`
	Num    int    `json:"num"`
	Price  string `json:"price"`
	Total  string `json:"total"`
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

func (os *OrderService) Insert(order mysql.Order, items []Item) (string, error) {
	tx, err := os.db.Begin()
	if err != nil {
		return "0", err
	}
	defer tx.Rollback()

	result, err := tx.Exec(os.SQLS[3], order.OrderId, order.Payment, order.Promotion, order.PostFee, order.CreateTime, order.CloseTime, order.UserID)
	if err != nil {
		return "0", err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return "0", errors.New("insert order: insert affected 0 rows")
	}

	for _, x := range items {
		result, err := tx.Exec(os.SQLS[4], x.ItemId, order.OrderId, x.Num, x.Price, x.Total)
		if err != nil {
			return "0", err
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return "0", errors.New("insert item: insert affected 0 rows")
		}
	}

	if err = tx.Commit(); err != nil {
		return "0", err
	}

	return order.OrderId, nil

}

func (os *OrderService) OrderInfo(orderid string) (*mysql.Order, error) {
	order, err := mysql.SelectByOrderKey(os.db, os.SQLS[5], orderid)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (os *OrderService) LisitOrderByUserId(userid int64) ([]*mysql.Order, error) {
	orders, err := mysql.LisitOrderByUserId(os.db, os.SQLS[6], userid)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *OrderService) LisitItemByOrderId(orderid string) ([]*mysql.Item, error) {
	items, err := mysql.LisitItemByOrderId(os.db, os.SQLS[7], orderid)
	if err != nil {
		return nil, err
	}

	return items, nil
}
func (os *OrderService) ChangePayTime(orderid string, paytime time.Time) error {
	return mysql.UpdateTimeByOrderKey(os.db, os.SQLS[8], orderid, paytime)
}

func (os *OrderService) ChangeConsignTime(orderid string, paytime time.Time) error {
	return mysql.UpdateTimeByOrderKey(os.db, os.SQLS[9], orderid, paytime)
}

func (os *OrderService) ChangeEndTime(orderid string, paytime time.Time) error {
	return mysql.UpdateTimeByOrderKey(os.db, os.SQLS[10], orderid, paytime)
}

func (os *OrderService) ChangeStatus(orderid string, status int) error {
	return mysql.UpdateStatusByOrderKey(os.db, os.SQLS[11], orderid, status)
}

func (os *OrderService) ChangeShipp(orderid string, shippingname, shippingcode string) error {
	return mysql.UpdateShipByOrderKey(os.db, os.SQLS[12], orderid, shippingname, shippingcode)
}

func (os *OrderService) ChangePaymentType(orderid string, paymentType int) error {
	return mysql.UpdatePaymentTypeByOrderKey(os.db, os.SQLS[13], orderid, paymentType)
}
