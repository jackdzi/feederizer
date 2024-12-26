package main

import (
  "log"

  "feederizer/server/internal/api"
  "feederizer/server/internal/db"
)

func main() {
  dbConn, err := db.NewDatabaseConnection()
  if err != nil {
    log.Fatalf("Failed to connect to database. Error: %v", err)
  }
  defer dbConn.Close()

  router := api.NewRouter(dbConn)
  log.Println("Starting server on :8080")
  if err := router.Run(":8080"); err != nil {
    log.Fatalf("Server failed. Error: %v", err)
  }
}
