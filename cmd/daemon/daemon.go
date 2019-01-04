package daemon

import (
	"log"

	"github.com/opensteel/authserver/cmd/database"
	"github.com/opensteel/authserver/cmd/handler"
	"github.com/opensteel/authserver/cmd/middleware"
	"github.com/opensteel/authserver/pkg/model"
)

// Config :
type Config struct {
	ListenSpec string

	Db database.Config
}

// Start :
func Start(cfg *Config) error {
	log.Printf("Autorization server on port %s started\n", cfg.ListenSpec)

	dab, err := database.InitDb(cfg.Db)
	if err != nil {
		log.Fatal("Fatal error with conecting or preparing databse")
		return err
	}
	log.Println("Database connected successful")

	m := model.New(dab)

	p := middleware.InitPrometheus()

	handler.Start(m, cfg.ListenSpec, p)

	return nil

}
