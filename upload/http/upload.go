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

	_ "github.com/go-sql-driver/mysql"

	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/base"
	"github.com/TechCatsLab/sor/base/constants"
	"github.com/TechCatsLab/sor/upload/mysql"

	"github.com/TechCatsLab/sor/base/filter"
)

type UploadController struct {
	*base.Controller
}

// UploadOne single file upload
func (u *UploadController) Upload(c *server.Context) error {
	if c.Request().Method != "POST" {
		log.Error("Request is not post method")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	userID, err := filter.GetUID(c, u.TokenKey())
	if userID == constants.InvalidUID {
		log.Error("userID invalid")
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	file, header, err := c.Request().FormFile(constants.FileKey)
	defer func() {
		file.Close()
		c.Request().MultipartForm.RemoveAll()
	}()
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	md5 := md5.New()
	_, err = io.Copy(md5, file)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	MD5Str := hex.EncodeToString(md5.Sum(nil))

	fileSuffix := path.Ext(header.Filename)
	filePath := constants.FileUploadDir + "/" + classifyBySuffix(fileSuffix) + "/" + MD5Str + fileSuffix
	cur, err := os.Create(filePath)
	defer cur.Close()
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	_, err = io.Copy(cur, file)
	if err != nil {
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}

	err = mysql.Insert(u.SQLStore(), userID, filePath, MD5Str)
	if err != nil {
		if err.Error() == "Error 1062: "+"Duplicate entry "+"'"+MD5Str+"'"+" for key 'PRIMARY'" {
			filePath, err := mysql.QueryByMD5(u.SQLStore(), MD5Str)
			if err != nil {
				log.Error(err)
				return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
			}
			return c.ServeJSON(respStatusAndData(http.StatusOK, u.BaseUrl()+filePath))
		}
		log.Error(err)
		return c.ServeJSON(respStatusAndData(http.StatusBadRequest, nil))
	}
	return c.ServeJSON(respStatusAndData(http.StatusOK, u.BaseUrl()+filePath))
}
