package mysql

import (
	"database/sql"
	"errors"
	"time"
)

var (
	errInvalidInsert       = errors.New("insert comment: insert affected 0 rows")
	errInvalidChangeTime   = errors.New("change time: affected 0 rows")
	errInvalidAddShip      = errors.New("update time: affected 0 rows")
	errInvalidChangeStatus = errors.New("change status: affected 0 rows")
)

type Order struct {
	ID          uint32 // todo: fix
	OrderId     string
	Payment     string // uint32
	PaymentType int
	Promotion   string // uint32
	PostFee     string
	Status      int
	CreateTime  time.Time
	PayTime     time.Time
	ConsignTime time.Time
	CloseTime   time.Time
	EndTime     time.Time
	ShipName    string // int32, uint32
	ShipCode    string
	UserID      int64
}

type Item struct {
	ID      uint32 // todo: fix
	ItemId  uint
	OrderId string
	Num     int
	Price   string
	Total   string
}

func CreateDB(db *sql.DB, createdb string) error {
	_, err := db.Exec(createdb)
	return err
}

func CreateTable(db *sql.DB, createtable string) error {
	_, err := db.Exec(createtable)
	return err
}

func SelectByOrderKey(db *sql.DB, query string, orderid string) (*Order, error) {
	var order Order

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&order.OrderId, &order.Payment, &order.PaymentType, &order.Promotion, &order.PostFee, &order.Status, &order.CreateTime, &order.PayTime, &order.ConsignTime, &order.CloseTime, &order.EndTime, &order.ShipName, &order.ShipCode, &order.UserID); err != nil {

			return nil, err
		}
	}

	return &order, nil
}

/* *
 * 初始阶段：1、未付款、未发货；初始化所有数据
 * 付款阶段：2、已付款；更改付款时间
 * 发货阶段：3、已发货；更改发货时间、物流名称、物流单号
 * 成功阶段：4、已成功；更改交易结束时间。
 * 关闭阶段：5、关闭：  更改交易关闭时间。
 * */
func UpdateTimeByOrderKey(db *sql.DB, update string, orderid string, time time.Time) error {
	result, err := db.Exec(update, time, orderid)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeTime
	}

	return nil
}

func UpdatePaymentTypeByOrderKey(db *sql.DB, update string, orderid string, paymenttype int) error {
	result, err := db.Exec(update, paymenttype, orderid)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeStatus
	}
	return nil
}
func UpdateShipByOrderKey(db *sql.DB, update string, orderid string, shipname string, shipcode string) error {
	result, err := db.Exec(update, shipname, shipcode, orderid)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidAddShip
	}
	return nil
}
func UpdateStatusByOrderKey(db *sql.DB, update string, orderid string, status int) error {
	result, err := db.Exec(update, status, orderid)
	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return errInvalidChangeStatus
	}
	return nil
}

func LisitOrderByUserId(db *sql.DB, query string, userid int64) ([]*Order, error) {
	var (
		orders []*Order

		orderId     string
		payment     string
		paymentType int
		promotion   string
		postFee     string
		status      int
		createTime  time.Time
		payTime     time.Time
		consignTime time.Time
		closeTime   time.Time
		endTime     time.Time
		shipName    string
		shipCode    string
		userID      int64
	)

	rows, err := db.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&orderId, &payment, &paymentType, &promotion, &postFee, &status, &createTime, &payTime, &consignTime, &closeTime, &endTime, &shipName, &shipCode, &userID); err != nil {
			return nil, err
		}

		order := &Order{
			OrderId:     orderId,
			Payment:     payment,
			PaymentType: paymentType,
			Promotion:   promotion,
			PostFee:     postFee,
			Status:      status,
			CreateTime:  createTime,
			PayTime:     payTime,
			ConsignTime: consignTime,
			CloseTime:   closeTime,
			EndTime:     endTime,
			ShipName:    shipName,
			ShipCode:    shipCode,
			UserID:      userID,
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func LisitItemByOrderId(db *sql.DB, query string, orderid string) ([]*Item, error) {
	var (
		itemid  uint
		orderId string
		num     int
		price   string
		total   string

		items []*Item
	)
	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&itemid, &orderId, &num, &price, &total); err != nil {
			return nil, err
		}
		item := &Item{
			ItemId:  itemid,
			OrderId: orderId,
			Num:     num,
			Price:   price,
			Total:   total,
		}
		items = append(items, item)
	}
	return items, nil
}

func DeleteItemByOrderId(db *sql.DB, delete string, orderId string) error {
	_, err := db.Exec(delete, orderId)
	return err
}

func DeleteOrderByOrderId(db *sql.DB, delete string, orderId string) error {
	_, err := db.Exec(delete, orderId)
	return err
}
