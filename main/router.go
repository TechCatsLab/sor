/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package main

import (
	"database/sql"

	"github.com/TechCatsLab/sor/smservice/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/sor/smservice"
)

var (
	router *server.Router
)

type funcv struct{}

func (v funcv) OnVerifySucceed(a, b string) {}
func (v funcv) OnVerifyFailed(a, b string)  {}

func init() {
	router = server.NewRouter()

	/*uploadDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/upload?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal(err)
	}
	upload.InitRouter(router, uploadDB, "http://127.0.0.1:9573", "UserTokenKey")

	adminDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal(err)
	}
	admin.InitAdminRouter(router, adminDB, "AdminTokenKey")
	*/
	smserviceDB, err := sql.Open("mysql", "root:yhyddr119216@tcp(127.0.0.1:3306)/?charset=utf8mb4")
	if err != nil {
		logrus.Fatal(err)
	}
	var v funcv
	c := &config.Config{
		Host:           "https://fesms.market.alicloudapi.com/sms/",
		Appcode:        "6f37345cad574f408bff3ede627f7014",
		Digits:         6,
		ResendInterval: 60,
		OnCheck:        v,
		DB:             smserviceDB,
	}
	smservice.Register(router, smserviceDB, c)
}
