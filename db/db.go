package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // go uses under the hood, special import (must be part of the project)
)

// can be used outside of the package
var DB *sql.DB // pointer

// initilize the DB, set up connections
func InitDB() {
	// return a pointer of DB, or error
	var err error
	DB, err = sql.Open("sqlite3", "api.db") // driver package, create file storing data

	if err != nil {
		panic("Could not connect to database.") // crash the program, create a log
	}

	// establish the connect successfully
	DB.SetMaxOpenConns(10) // # of open connection we can have to DB at the same time
	DB.SetMaxIdleConns(5)  // if no one is using. open max 5

	createTables() // create tables
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic("Could not create users table")
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createEventsTable) // execute the create table SQL statement
	if err != nil {
		panic("Could not create events table.")
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER,
		user_id INTEGER,
		FOREIGN KEY(event_id) REFERENCES events(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createRegistrationTable)
	if err != nil {
		panic("Could not create resgitrations table.")
	}

}
