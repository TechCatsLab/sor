package controller

import (
	"github.com/TechCatsLab/sor/smservice/config"
	"github.com/TechCatsLab/sor/smservice/constants"
	"github.com/TechCatsLab/sor/smservice/services"

	log "github.com/TechCatsLab/logging/logrus"

	"database/sql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/base"
)

type Controller struct {
	service *services.SMservice
}

func New(db *sql.DB, c *config.Config) *Controller {
	return &Controller{
		service: services.NewService(c, db),
	}
}

func (con *Controller) CreateDB() error {
	return con.service.CreateDB()
}

func (con *Controller) CreateTable() error {
	return con.service.CreateTable()
}

//调度分配出发送短信
func (con *Controller) Send(c *server.Context) error {
	var (
		req struct {
			Mobile string `json:"mobile"`
			Sign   string `json:"sign"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	if err := services.Send(req.Mobile, req.Sign, con.service.Conf, con.service.Db); err != nil {
		log.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrSendinMysql, nil))
	}

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, nil))
}

//调度分配检查验证码
func (con *Controller) Check(c *server.Context) error {
	var (
		req struct {
			Code string `json:"code"`
			Sign string `json:"sign"`
		}
	)

	if err := c.JSONBody(&req); err != nil {
		log.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrInvalidParam, nil))
	}

	var resp struct {
		sign   string
		mobile string
	}
	resp.sign = req.Sign
	resp.mobile, _ = con.service.GetMobile(resp.sign)

	if err := services.Check(req.Code, req.Sign, con.service.Conf, con.service.Db); err != nil {

		con.service.Conf.OnCheck.OnVerifyFailed(resp.sign, resp.mobile)

		log.Error(err)
		return c.ServeJSON(base.RespStatusAndData(constants.ErrCheckCode, nil))
	}

	con.service.Conf.OnCheck.OnVerifySucceed(resp.sign, resp.mobile)

	return c.ServeJSON(base.RespStatusAndData(constants.ErrSucceed, resp))
}
