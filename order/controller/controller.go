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
			Payment   string         `json:"payment"`
			Promotion string         `json:"promotion"`
			PostFee   string         `json:"postfee"`
			UserId    int64          `json:"userid"`
			Items     []service.Item `json:"items"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	times := time.Now()
	orderid := strconv.Itoa(times.Year()) + strconv.Itoa(int(times.Month())) + strconv.Itoa(times.Day()) + strconv.Itoa(int(req.UserId))
	order := &mysql.Order{
		OrderId:    orderid,
		Payment:    req.Payment,
		Promotion:  req.Promotion,
		PostFee:    req.PostFee,
		UserID:     req.UserId,
		CreateTime: times,
		CloseTime:  times.AddDate(0, 0, 1),
	}
	id, err := ctl.service.Insert(*order, req.Items)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrCreateInMysql, nil))
	}

	for _, x := range req.Items {
		ctl.service.Cnf.ModifyStock.ModifyProductStock(x.ItemId, (-1 * x.Num))
	}

	return base.WriteStatusAndIDJSON(c, constants.ErrSucceed, id)

}

func (ctl *Controller) OrderInfoByOrderId(c *server.Context) error {
	var (
		req struct {
			OrderId string `json:"orderid"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	order, err := ctl.service.OrderInfo(req.OrderId)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))

	}
	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, order))
}

func (ctl *Controller) LisitItemByOrderId(c *server.Context) error {
	var (
		req struct {
			OrderId string `json:"orderid"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	items, err := ctl.service.LisitItemByOrderId(req.OrderId)
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, items))
}

func (ctl *Controller) LisitOrderByUserId(c *server.Context) error {
	var (
		req struct {
			Userid string `json:"userid"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	id, err := strconv.Atoi(req.Userid)
	if err != nil {
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	orders, err := ctl.service.LisitOrderByUserId(int64(id))
	if err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, orders))
}

func (ctl *Controller) Pay(c *server.Context) error {
	var (
		req struct {
			OrderId     string `json:"orderid"`
			PaymentType int    `json:"paymenttype"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if err := ctl.service.ChangePayTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	if err := ctl.service.ChangePaymentType(req.OrderId, req.PaymentType); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	if err := ctl.service.ChangeStatus(req.OrderId, 2); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, constants.ErrSucceed))

}

func (ctl *Controller) Consign(c *server.Context) error {
	var (
		req struct {
			OrderId      string `json:"orderid"`
			ShippingName string `json:"shippingname"`
			ShippingCode string `json:"shippingcode"`
		}
	)
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	if err := ctl.service.ChangeConsignTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}
	if err := ctl.service.ChangeShipp(req.OrderId, req.ShippingName, req.ShippingCode); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	if err := ctl.service.ChangeStatus(req.OrderId, 3); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, constants.ErrSucceed))

}

func (ctl *Controller) Success(c *server.Context) error {
	var (
		req struct {
			OrderId string `json:"orderid"`
		}
	)
	if err := c.JSONBody(&req); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}
	if err := ctl.service.ChangeEndTime(req.OrderId, time.Now()); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))

	}
	if err := ctl.service.ChangeStatus(req.OrderId, 4); err != nil {
		logrus.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrOprationInMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, constants.ErrSucceed))

}
