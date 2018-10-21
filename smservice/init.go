package smservice

import (
	"database/sql"

	"github.com/TechCatsLab/logging/logrus"

	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/sor/smservice/config"
	"github.com/TechCatsLab/sor/smservice/controller"
)

func Register(r *server.Router, db *sql.DB, c *config.Config) error {
	if r == nil {
		logrus.Fatal("[InitRouter]: server is nil")
	}

	con := controller.New(db, c)

	if err := con.CreateDB(); err != nil {
		return err
	}

	if err := con.CreateTable(); err != nil {
		return err
	}

	r.Post("/api/v1/sms/send", con.Send)
	r.Post("/api/v1/sms/check", con.Check)

	return nil
}
