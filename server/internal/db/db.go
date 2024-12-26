package db

import (
  "github.com/jmoiron/sqlx"
  _ "github.com/mattn/go-sqlite3"
)

func NewDatabaseConnection() (*sqlx.DB, error) {
  db, err := sqlx.Open("sqlite3", "../../feederizer.db")
  if err != nil {
    return nil, err
  }
  db.Exec("PRAGMA foreign_keys = ON;")

  return db, nil
}
