package mysql

import (
	"database/sql"
	"errors"
	"time"
)

//can update:shipname
type Order struct {
	ID         uint32
	OrderCode  string    `json:"ordercode"`
	ShipCode   string    `json:"shipcode"`
	UserID     uint64    `json:"userid"`
	AddressID  string    `json:"addressid"`
	TotalPrice uint32    `json:"totalprice"`
	PayWay     uint8     `json:"payway"`
	Promotion  uint32    `json:"promotion"`
	Freight    uint32    `json:"freight"`
	Status     uint8     `json:"status"`
	Created    time.Time `json:"created"`
	Closed     time.Time `json:"closed"`
	Updated    time.Time `json:"updated"`
}

type Item struct {
	ID        uint32
	ProductId uint32 `json:"productid"`
	OrderID   uint32 `json:"orderid"`
	Count     uint32 `json:"count"`
	Price     uint32 `json:"price"`
	Discount  uint8  `json:"discount"`
	Size      string `json:"size"`//to do fix
	Color     string `json:"color"`
}

func CreateDB(db *sql.DB, createDB string) error {
	_, err := db.Exec(createDB)
	return err
}

func CreateTable(db *sql.DB, createTable string) error {
	_, err := db.Exec(createTable)
	return err
}

func UpdateTimeByOrderKey(db *sql.DB, update string, orderid uint32, time time.Time) (uint32, error) {
	result, err := db.Exec(update, time, orderid)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not affected time")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func UpdatePayWayByOrderKey(db *sql.DB, update string, orderid uint32, payway uint8) (uint32, error) {
	result, err := db.Exec(update, payway, orderid)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] ; not affected paywqy")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func UpdateShipByOrderKey(db *sql.DB, update string, orderid uint32, shipcode string) (uint32, error) {
	result, err := db.Exec(update, shipcode, orderid)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not affected shiporder")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func UpdateStatusByOrderKey(db *sql.DB, update string, orderid uint32, status uint8) (uint32, error) {
	result, err := db.Exec(update, status, orderid)
	if err != nil {
		return 0, err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return 0, errors.New("[change error] : not affected status")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

func OrderIDByOrderCode(db *sql.DB, query string, ordercode string) (uint32, error) {
	var (
		orderid uint32
		err     error
	)

	rows, err := db.Query(query, ordercode)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&orderid); err != nil {
			return 0, err
		}
	}

	return orderid, nil
}

func SelectByOrderKey(db *sql.DB, query string, orderid uint32) (*Order, error) {
	var order Order

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&order.ID, &order.OrderCode, &order.ShipCode, &order.UserID, &order.AddressID, &order.TotalPrice, &order.PayWay, &order.Promotion, &order.Freight, &order.Status, &order.Created, &order.Closed, &order.Updated); err != nil {
			return nil, err
		}
	}

	return &order, nil
}

func LisitOrderByUserId(db *sql.DB, query string, userid uint64) ([]*Order, error) {
	var (
		orders []*Order

		ID         uint32
		OrderCode  string
		ShipCode   string
		UserID     uint64
		AddressID  string
		TotalPrice uint32
		PayWay     uint8
		Promotion  uint32
		Freight    uint32
		Status     uint8
		Created    time.Time
		Closed     time.Time
		Updated    time.Time
	)

	rows, err := db.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ID, &OrderCode, &ShipCode, &UserID, &AddressID, &TotalPrice, &PayWay, &Promotion, &Freight, &Status, &Created, &Closed, &Updated); err != nil {
			return nil, err
		}

		order := &Order{
			ID:         ID,
			OrderCode:  OrderCode,
			ShipCode:   ShipCode,
			UserID:     UserID,
			AddressID:  AddressID,
			TotalPrice: TotalPrice,
			PayWay:     PayWay,
			Promotion:  Promotion,
			Freight:    Freight,
			Status:     Status,
			Created:    Created,
			Closed:     Closed,
			Updated:    Updated,
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func LisitItemByOrderId(db *sql.DB, query string, orderid uint32) ([]*Item, error) {
	var (
		ID        uint32
		ProductId uint32
		OrderID   uint32
		Count     uint32
		Price     uint32
		Discount  uint8
		Size      string
		Color     string

		items []*Item
	)

	rows, err := db.Query(query, orderid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&ID, &ProductId, &OrderID, &Count, &Price, &Discount, &Size, &Color); err != nil {
			return nil, err
		}

		item := &Item{
			ID:        ID,
			ProductId: ProductId,
			OrderID:   OrderID,
			Count:     Count,
			Price:     Price,
			Discount:  Discount,
			Size:      Size,
			Color:     Color,
		}
		items = append(items, item)
	}

	return items, nil
}
