package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackdzi/feederizer/server/internal/db"
)

func main() {
  fmt.Println(os.Getwd())
  dbConn, err := db.NewDatabaseConnection()
  if err != nil {
    log.Fatalf("Failed to connect to database. Error: %v", err)
  }
  defer dbConn.Close()

  router := db.NewRouter(dbConn, false)
  log.Println("Starting server on :8080")
  if err := router.Run(":8080"); err != nil {
    log.Fatalf("Server failed. Error: %v", err)
  }
}
