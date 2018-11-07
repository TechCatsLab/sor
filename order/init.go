package order

import (
	"database/sql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/order/config"
	"github.com/TechCatsLab/sor/order/controller"
)

func Register(r *server.Router, db *sql.DB, cnf *config.Config) error {
	if r == nil {
		logrus.Fatal("[InitRouter]: server is nil")
	}

	c := controller.New(db, cnf)

	if err := c.CreateDB(); err != nil {
		return err
	}

	if err := c.CreateOrderTable(); err != nil {
		return err
	}

	if err := c.CreateItemTable(); err != nil {
		return err
	}

	r.Post("/api/v1/order/create", c.Insert)
	r.Post("/api/v1/order/info/id", c.OrderInfoByOrderId)
	r.Post("/api/v1/order/lisit/user", c.LisitOrderByUserId)
	r.Post("/api/v1/order/lisit/order", c.LisitItemByOrderId)
	r.Post("/api/v1/order/pay", c.Pay)
	r.Post("/api/v1/order/consign", c.Consign)
	r.Post("/api/v1/order/success", c.Success)

	return nil
}