/*
 * Revision History:
 *     Initial: 2018/08/13        Shi Ruitao
 */

package base

import (
	"database/sql"
)

type Controller struct {
	db       *sql.DB
	baseUrl  string
	tokenKey string
}

func New(db *sql.DB, baseUrl, tokenKey string) *Controller {
	return &Controller{
		db:       db,
		baseUrl:  baseUrl,
		tokenKey: tokenKey,
	}
}

func (c *Controller) SQLStore() *sql.DB {
	return c.db
}

func (c *Controller) BaseUrl() string {
	return c.baseUrl
}

func (c *Controller) TokenKey() string {
	return c.tokenKey
}
