/*
 * Revision History:
 *     Initial: 2018/09/02        Shi Ruitao
 */

package base

import (
	"github.com/TechCatsLab/sor/base/constants"
)

func RespStatusAndData(statusCode int, data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{constants.RespKeyStatus: statusCode}
	}
	return map[string]interface{}{constants.RespKeyStatus: statusCode, constants.RespKeyData: data}
}