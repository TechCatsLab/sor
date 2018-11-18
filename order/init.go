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
	r.Post("/api/v1/order/info", c.OrderInfoByOrderID)
	//fix to more info
	r.Post("/api/v1/order/user", c.LisitOrderByUserID)
	r.Post("/api/v1/order/id", c.OrderIDByOrderCode)
	r.Post("/api/v1/order/pay", c.Pay)
	r.Post("/api/v1/order/consign", c.Consign)
	r.Post("/api/v1/order/success", c.Success)
	r.Post("/api/v1/order/cancel", c.Cancel)

	return nil
}
