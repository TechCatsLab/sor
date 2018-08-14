/*
 * Revision History:
 *     Initial: 2018/08/13        Shi Ruitao
 */

package base

import (
	"database/sql"
)

type Controller struct {
	db *sql.DB
	baseUrl string
}

func New(db *sql.DB, baseUrl string) *Controller {
	return &Controller{
		db: db,
		baseUrl: baseUrl,
	}
}

func (c *Controller) SQLStore() *sql.DB {
	return c.db
}

func (c *Controller) BaseUrl() string {
	return c.baseUrl
}
