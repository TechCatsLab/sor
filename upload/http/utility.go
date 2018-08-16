/*
 * Revision History:
 *     Initial: 2018/08/14        Shi Ruitao
 */

package http

import (
	"github.com/TechCatsLab/sor/base/constants"
)

var (
	fileMap = map[string]string{}
	picture = []string{".jpg", ".png", ".jpeg", ".gif", ".bmp"}
	video   = []string{".avi", ".wmv", ".mpg", ".mpeg", ".mpe", ".mov", ".rm", ".ram", ".swf", ".mp4", ".rmvb", ".asf", ".divx", ".vob"}
	fileDir = filePath()
)

func filePath() map[string]string {
	for _, suffix := range picture {
		fileMap[suffix] = constants.PictureDir
	}
	for _, suffix := range video {
		fileMap[suffix] = constants.VideoDir
	}
	return fileMap
}

func respStatusAndData(statusCode int, data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{constants.RespKeyStatus: statusCode}
	}
	return map[string]interface{}{constants.RespKeyStatus: statusCode, constants.RespKeyData: data}
}

func classifyBySuffix(suffix string) string {

	if dir := fileDir[suffix]; dir != "" {
		return dir
	}
	return constants.OtherDir
}
