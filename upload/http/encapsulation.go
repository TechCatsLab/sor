/*
 * Revision History:
 *     Initial: 2018/08/14        Shi Ruitao
 */

package http

import (
	"github.com/TechCatsLab/sor/base"
)

var(
	fileMap = map[string]string{}
	image = []string{".jpg", ".png", ".jpeg", ".gif", ".bmp"}
	video = []string{".avi", ".wmv", ".mpg", ".mpeg", ".mpe", ".mov", ".rm", ".ram", ".swf", ".mp4", ".rmvb", ".asf", ".divx", ".vob"}
	fileDir = filePath()
)

func filePath() map[string]string {
	for _, suffix := range image {
		fileMap[suffix] = "/image/"
	}
	for _, suffix := range video {
		fileMap[suffix] = "/video/"
	}
	return fileMap
}

func respStatusAndData(statusCode int, data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{base.RespKeyStatus: statusCode}
	}
	return map[string]interface{}{base.RespKeyStatus: statusCode, base.RespKeyData: data}
}

func classifyBySuffix(suffix string) string {

	if dir := fileDir[suffix]; dir != "" {
		return dir
	}
	return otherFile
}
