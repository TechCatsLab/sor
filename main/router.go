/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package main

import (
	"database/sql"

	"github.com/TechCatsLab/sor/category"

	"github.com/TechCatsLab/sor/category/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/logging/logrus"
)

var (
	router *server.Router
)

type funcv struct{}

func (v funcv) OnVerifySucceed(a, b string) {}
func (v funcv) OnVerifyFailed(a, b string)  {}

func init() {
	router = server.NewRouter()
	//之前学长自己做的测试，upload和admin文件夹
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

	//这是发短信服务的测试，目前返回值仍没搞懂怎样返回resp（是空的）
	/*smserviceDB, err := sql.Open("mysql", "root:yhyddr119216@tcp(127.0.0.1:3306)/?charset=utf8mb4")
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
	*/
	//这是目录的测试
	CategoryDB, err := sql.Open("mysql", "root:yhyddr119216@tcp(127.0.0.1:3306)/?parseTime=true")
	if err != nil {
		logrus.Fatal(err)
	}

	cc := &config.Config{
		CategoryDB:    "mall",
		CategoryTable: "table",
	}
	category.Register(router, CategoryDB, cc)
}
