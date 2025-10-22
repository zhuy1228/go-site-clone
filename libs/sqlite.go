package libs

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct{}

func (Sqlite) Test() {
	db, err := sql.Open("sqlite3", "./data/example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
  );`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO users (name) VALUES (?)`, "alice")
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = db.QueryRow(`SELECT name FROM users WHERE id = ?`, 1).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("user #1:", name)
}
