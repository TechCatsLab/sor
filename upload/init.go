/*
 * Revision History:
 *     Initial: 2018/08/10        Shi Ruitao
 */

package upload

import (
	"database/sql"
	"os"

	"github.com/TechCatsLab/apix/http/server"
	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/base"
	"github.com/TechCatsLab/sor/base/constants"
	"github.com/TechCatsLab/sor/upload/http"
	"github.com/TechCatsLab/sor/upload/mysql"
	"github.com/TechCatsLab/sor/base/filter"
)

func InitRouter(r *server.Router, db *sql.DB, baseUrl, tokenKey string) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.Create(db)
	if err != nil {
		log.Fatal(err)
	}
	err = checkDir(constants.PictureDir, constants.VideoDir, constants.OtherDir)
	if err != nil {
		log.Fatal(err)
	}

	c := &http.UploadController{
		base.New(db),
		baseUrl,
	}

	jwt := filter.New(tokenKey)

	r.Post("/api/v1/user/upload", c.Upload, jwt.Check)
}

func checkDir(path ...string) error {
	for _, name := range path {
		_, err := os.Stat(constants.FileUploadDir + "/" + name)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(constants.FileUploadDir+"/"+name, 0777)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
