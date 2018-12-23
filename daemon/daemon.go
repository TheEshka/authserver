package daemon

import (
	"fmt"

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
	fmt.Println("Autorization server on port %s started\n", cfg.ListenSpec)

	db, err := db.InitDb(cfg.Db)
	if err != nil {
		fmt.Println("Database connecting error")
		return err
	} else {
		fmt.Println("Database connected successful")
	}

	m := model.New(db)

	handler.Start(m, cfg.ListenSpec)

	return nil

}
