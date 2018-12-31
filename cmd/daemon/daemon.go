package daemon

import (
	"log"

	"github.com/opensteel/authserver/cmd/db"
	"github.com/opensteel/authserver/cmd/handler"
	"github.com/opensteel/authserver/cmd/middleware"
	"github.com/opensteel/authserver/pkg/model"
)

// Config :
type Config struct {
	ListenSpec string

	Db db.Config
}

// Start :
func Start(cfg *Config) error {
	log.Printf("Autorization server on port %s started\n", cfg.ListenSpec)

	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Fatal("Fatal error with conecting or preparing databse")
		return err
	} else {
		log.Println("Database connected successful")
	}

	m := model.New(db)

	p := middleware.InitPrometheus()

	handler.Start(m, cfg.ListenSpec, p)

	return nil

}
