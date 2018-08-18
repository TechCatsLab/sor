/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package main

import (
	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/apix/http/server/middleware"
	"github.com/TechCatsLab/logging/logrus"
)

func main() {
	serverConfig := &server.Configuration{
		Address: "127.0.0.1:9573",
	}

	ep := server.NewEntrypoint(serverConfig, nil)
	ep.AttachMiddleware(middleware.NegroniRecoverHandler())
	ep.AttachMiddleware(middleware.NegroniJwtHandler("TokenHMACKey", nil, nil, nil))

	if err := ep.Start(router.Handler()); err != nil {
		logrus.Error(err)
		return
	}

	ep.Wait()
}
