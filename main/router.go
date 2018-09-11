/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/upload"
	"github.com/TechCatsLab/sor/admin"
)

var (
	router *server.Router
)

func init() {
	router = server.NewRouter()

	uploadDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/upload?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal(err)
	}

	upload.InitRouter(router, uploadDB, "http://127.0.0.1:9573", "UserTokenKey")
	adminDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal(err)
	}
	admin.InitAdminRouter(router, adminDB, "AdminTokenKey")
}
