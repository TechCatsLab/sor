/*
 * Revision History:
 *     Initial: 2018/08/18        Feng Yifei
 */

package base

import (
	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/base/constants"
)

const (
	CtxKeyUID = "uid"
)

type Context struct {
	*server.Context
}

func (c *Context) SetUID(id uint) {
	c.Set(CtxKeyUID, id)
}

func (c *Context) UID() uint {
	v := c.Get(CtxKeyUID)
	if v == nil {
		return constants.InvalidUID
	}

	if i, ok := v.(uint); ok {
		return i
	}

	return constants.InvalidUID
}
