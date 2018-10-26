/*
 * Revision History:
 *     Initial: 2018/09/02        Shi Ruitao
 */

package base

import (
	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/base/constants"
)

func RespStatusAndData(statusCode int, data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{constants.RespKeyStatus: statusCode}
	}
	return map[string]interface{}{constants.RespKeyStatus: statusCode, constants.RespKeyData: data}
}

//做的resp加进来，因为还不懂学长写的response
func WriteStatusAndIDJSON(ctx *server.Context, status int, id interface{}) error {
	return ctx.ServeJSON(map[string]interface{}{
		constants.RespKeyStatus: status,
		"ID": id,
	})
}

func WriteStatusAndDataJSON(ctx *server.Context, status int, data interface{}) error {
	if data == nil {
		return ctx.ServeJSON(map[string]interface{}{constants.RespKeyStatus: status})
	}

	return ctx.ServeJSON(map[string]interface{}{
		constants.RespKeyStatus: status,
		constants.RespKeyData:   data,
	})
}
