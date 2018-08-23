/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package main

import (
	"time"
	"fmt"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/apix/http/server/middleware"
	"github.com/TechCatsLab/logging/logrus"
	jwtgo "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	serverConfig := &server.Configuration{
		Address: "127.0.0.1:9573",
	}

	ep := server.NewEntrypoint(serverConfig, nil)
	ep.AttachMiddleware(middleware.NegroniRecoverHandler())
	ep.AttachMiddleware(middleware.NegroniJwtHandler("TokenHMACKey", nil, nil, nil))

	if err := ep.Start(router.Handler()); err != nil {
		logrus.Fatal(err)
	}

	t, err := NewToken(1, true)
	fmt.Println(t, err)

	ep.Wait()
}

func NewToken(userID uint, admin bool) (string, error) {
	claims := make(jwtgo.MapClaims)
	claims["uid"] = userID
	claims["admin"] = admin
	claims["exp"] = time.Now().Add(time.Hour * 480).Unix()
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	return token.SignedString([]byte("TokenHMACKey"))
}
