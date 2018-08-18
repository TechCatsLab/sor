/*
 * Revision History:
 *     Initial: 2018/08/10        Shi Ruitao
 */

package http

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/base"
	"github.com/TechCatsLab/sor/base/constants"
	"github.com/TechCatsLab/sor/upload/mysql"
)

type UploadController struct {
	*base.Controller
	BaseURL string
}

// UploadOne single file upload
func (u *UploadController) Upload(c *server.Context) error {
	if c.Request().Method != "POST" {
		log.Error("Request is not post method")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	ctx := &base.Context{c}
	userID := ctx.UID()

	if userID == constants.InvalidUID {
		log.Error("userID invalid")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	file, header, err := ctx.Request().FormFile(constants.FileKey)
	defer func() {
		file.Close()
		ctx.Request().MultipartForm.RemoveAll()
	}()
	if err != nil {
		log.Error(err)
		return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		log.Error(err)
		return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	MD5Str := hex.EncodeToString(md5.Sum(nil))

	fileSuffix := path.Ext(header.Filename)
	filePath := constants.FileUploadDir + "/" + classifyBySuffix(fileSuffix) + "/" + MD5Str + fileSuffix
	cur, err := os.Create(filePath)
	defer cur.Close()
	if err != nil {
		log.Error(err)
		return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	_, err = io.Copy(cur, file)
	if err != nil {
		log.Error(err)
		return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	err = mysql.Insert(u.SQLStore(), userID, filePath, MD5Str)
	if err != nil {
		if err.Error() == "Error 1062: "+"Duplicate entry "+"'"+MD5Str+"'"+" for key 'PRIMARY'" {
			filePath, err := mysql.QueryByMD5(u.SQLStore(), MD5Str)
			if err != nil {
				log.Error(err)
				return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
			}
			return ctx.ServeJSON(respStatusAndData(http.StatusOK, u.BaseURL+filePath))
		}
		log.Error(err)
		return ctx.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	return ctx.ServeJSON(respStatusAndData(http.StatusOK, u.BaseURL+filePath))
}
