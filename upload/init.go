/*
 * Revision History:
 *     Initial: 2018/08/10        Shi Ruitao
 */

package upload

import (
	"os"
	"database/sql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/upload/http"
	"github.com/TechCatsLab/sor/base"
	"github.com/TechCatsLab/sor/upload/mysql"
	log "github.com/TechCatsLab/logging/logrus"
)

const (
	imageFilePath = "./files/image/"
	videoFilePath = "./files/video/"
	otherFilePath = "./files/other/"
)

func checkDir(path... string) error {
	for _, name := range path{
		_, err := os.Stat(name)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(name, 0777)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func InitRouter(r *server.Router, db *sql.DB, baseUrl string) {
	if r == nil {
		panic("[InitRouter]: server is nil")
	}

	err := mysql.Create(db)
	if err != nil {
		log.Warn("Mysql Error:", err)
	}
	err = checkDir(imageFilePath, videoFilePath, otherFilePath)
	if err != nil {
		log.Warn("mkdir error:")
	}

	c := &http.UploadController{
		base.New(db, baseUrl),
	}

	r.Post("/api/v1/upload", c.Upload)
}
