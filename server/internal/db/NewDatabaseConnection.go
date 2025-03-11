package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pelletier/go-toml"
)

func NewDatabaseConnection() (*sqlx.DB, error) {
	var docker bool
	config, err := os.ReadFile("../config.toml")
	// if err != nil {
	//    config, err1 := os.ReadFile("config.toml")
	// }

	tree, err := toml.Load(string(config))
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	docker = tree.Get("deployment.docker").(bool)
  path := "../server/data/feederizer.db"


	//path := "data/feederizer.db" Uncomment if running server in seperate instance



	if docker {
		path = "/data/feederizer.db"
	}
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.Exec("PRAGMA foreign_keys = ON;")

	return db, nil
}
