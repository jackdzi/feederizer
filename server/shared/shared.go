package shared

import (
	"log"

	"github.com/jackdzi/feederizer/server/internal/db"
)

func RunServer() {
	dbConn, err := db.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}
	defer dbConn.Close()

	router := db.NewRouter(dbConn, true)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed. Error: %v", err)
	}
}
