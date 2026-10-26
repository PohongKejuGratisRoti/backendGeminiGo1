package main

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

// sqlite declaration
var db *sql.DB

func initDB() error {
	var DbErr error

	db, DbErr = sql.Open("sqlite", "notsolusindoiterasi3.db")
	if DbErr != nil {
		return DbErr
	}
	if DbErr = db.Ping(); DbErr != nil {
		db.Close()
		return DbErr
	}

	_, DbErr = db.Exec(`
		CREATE TABLE IF NOT EXISTS GSusers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username text not null,
			passsword text not null
		)
	`)
	if DbErr != nil {
		db.Close()
		return DbErr
	}
	return nil
}

func checkUser(username string, passsword string) bool {
	var savedPassword string
	DbErr := db.QueryRow("select passsword from gsusers where username = ?", username).Scan(&savedPassword)
	if DbErr != nil {
		return false
	}
	DbErr = bcrypt.CompareHashAndPassword(
		[]byte(savedPassword),
		[]byte(passsword),
	)

	return DbErr == nil
}

func checkID(username string) (int, error) {
	var savedID int
	err := db.QueryRow("select id from gsusers where username = ?", username).Scan(&savedID)
	if err != nil {
		return 0, err
	}
	return savedID, nil
}
