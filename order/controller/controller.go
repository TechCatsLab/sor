package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/TechCatsLab/sor/base/constants"

	"github.com/TechCatsLab/sor/order/model/mysql"

	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/base"

	"github.com/TechCatsLab/apix/http/server"

	"github.com/TechCatsLab/sor/order/config"
	"github.com/TechCatsLab/sor/order/service"
)

type Controller struct {
	service *service.OrderService
}

func New(db *sql.DB, cnf *config.Config) *Controller {
	return &Controller{
		service: service.NewOrderService(cnf, db),
	}
}

func (ctl *Controller) CreateDB() error {
	return ctl.service.CreateDB()
}

func (ctl *Controller) CreateOrderTable() error {
	return ctl.service.CreateOrderTable()
}
func (ctl *Controller) CreateItemTable() error {
	return ctl.service.CreateItemTable()
}

func (ctl *Controller) Insert(c *server.Context) error {
	var (
		req struct {
			UserID     uint64 `json:"userid"`
			AddressID  string `json:"addressid"`
			TotalPrice uint32 `json:"totalprice"`
			Promotion  uint32 `json:"promotion"`
			Freight    uint32 `json:"freight"`

			Items []mysql.Item `json:"items"`
		}
		rep struct {
			ordercode string
			orderid   uint32
		}
		err error
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	times := time.Now()

	rep.ordercode = strconv.Itoa(times.Year()) + strconv.Itoa(int(times.Month())) + strconv.Itoa(times.Day()) + strconv.Itoa(times.Hour()) + strconv.Itoa(times.Minute()) + strconv.Itoa(times.Second()) + strconv.Itoa(int(req.UserID))
	order := mysql.Order{
		OrderCode:  rep.ordercode,
		UserID:     req.UserID,
		AddressID:  req.AddressID,
		TotalPrice: req.TotalPrice,
		Freight:    req.Freight,
		Created:    times,
	}

	rep.orderid, err = ctl.service.Insert(order, req.Items)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrCreateInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndIDCODEData(constants.ErrSucceed, rep.orderid, rep.ordercode))
}

func (ctl *Controller) OrderInfoByOrderID(c *server.Context) error {
	var (
		req struct {
			OrderId uint32 `json:"orderid"`
		}
		rep struct {
			order *mysql.Order
			items []*mysql.Item
		}
		err error
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	rep.order, err = ctl.service.OrderInfo(req.OrderId)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	rep.items, err = ctl.service.LisitItemByOrderId(rep.order.ID)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndTwoData(constants.ErrSucceed, rep.order, rep.items))

}

func (ctl *Controller) LisitOrderByUserID(c *server.Context) error {
	var (
		req struct {
			Userid uint64 `json:"userid"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	orders, err := ctl.service.LisitOrderByUserId(req.Userid)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, orders))
}

func (ctl *Controller) OrderIDByOrderCode(c *server.Context) error {
	var (
		req struct {
			Ordercode string `json:"ordercode"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	id, err := ctl.service.OrderIDByOrderCode(req.Ordercode)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))
}
//to do fix
func (ctl *Controller) Pay(c *server.Context) error {
	var (
		req struct {
			OrderId uint32 `json:"orderid"`
			PayWay  uint8  `json:"payway"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if _, err := ctl.service.UpdateTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	if _, err := ctl.service.UpdatePayWay(req.OrderId, req.PayWay); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	id, err := ctl.service.UpdateStatus(req.OrderId, constants.OrderPaid)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))
}

func (ctl *Controller) Consign(c *server.Context) error {
	var (
		req struct {
			OrderId      uint32 `json:"orderid"`
			ShippingCode string `json:"shippingcode"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if _, err := ctl.service.UpdateTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}
	if _, err := ctl.service.UpdateShip(req.OrderId, req.ShippingCode); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	id, err := ctl.service.UpdateStatus(req.OrderId, constants.OrderConsign)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))

}

func (ctl *Controller) Success(c *server.Context) error {
	var (
		req struct {
			OrderId uint32 `json:"orderid"`
		}
	)
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if _, err := ctl.service.UpdateTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))

	}

	id, err := ctl.service.UpdateStatus(req.OrderId, constants.OrderFinished)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))
}

func (ctl *Controller) Cancel(c *server.Context) error {
	var (
		req struct {
			OrderId uint32 `json:"orderid"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if _, err := ctl.service.UpdateTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	id, err := ctl.service.UpdateStatus(req.OrderId, constants.OrderCanceled)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, id))
}
