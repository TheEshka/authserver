package daemon

import (
	"log"

	"github.com/me/authserver/db"
	"github.com/me/authserver/handler"
	"github.com/me/authserver/model"
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

	handler.Start(m, cfg.ListenSpec)

	return nil

}
