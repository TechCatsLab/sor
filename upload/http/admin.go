/*
 * Revision History:
 *     Initial: 2018/08/10        Shi Ruitao
 */

package http

import (
	"net/http"

	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/upload/mysql"
	"github.com/TechCatsLab/sor/base/filter"
	"github.com/TechCatsLab/sor/base/constants"
)

func (u *UploadController) QueryByTime(c *server.Context) error {
	var t struct {
		Time string `json:"time"`
	}

	isAdmin, err := filter.GetAdmin(c, u.TokenKey())
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	if !isAdmin {
		log.Warn("Not an administrator")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	err = c.JSONBody(&t)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	files, err := mysql.QueryByTime(u.SQLStore(), t.Time)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	return c.ServeJSON(respStatusAndData(http.StatusOK, files))
}

func (u *UploadController) QueryByUserID(c *server.Context) error {
	var user struct {
		UserID uint `json:"user_id"`
	}

	isAdmin, err := filter.GetAdmin(c, u.TokenKey())
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	if !isAdmin {
		log.Warn("Not an administrator")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	err = c.JSONBody(&user)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	if user.UserID == constants.InvalidUID {
		log.Error("userID invalid")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	files, err := mysql.QueryByUserID(u.SQLStore(), user.UserID)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	return c.ServeJSON(respStatusAndData(http.StatusOK, files))
}
