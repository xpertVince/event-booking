package models

import (
	"time"

	"example.com/eveny-booking/db"
)

// shape of event, data makes up an event
type Event struct {
	// would generate error if missing required data
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64     // creator of the event
}

var events = []Event{}

// save event to DB
func (e *Event) Save() error {
	// inserting values by ?, to avoid attack
	query := `
	INSERT INTO events(name, description, location, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query) // store in memory, re-use it (better performance)
	if err != nil {
		return err
	}

	// close the statement
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID) // execute this statement, one for each ?, in the order of column name
	if err != nil {
		return err
	}

	id, err := result.LastInsertId() // id of row just inserted (auto increment)
	e.ID = id

	return err

	// // later: add it to the database
	// events = append(events, e)
}

// not a method, normal func
func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query) // Query: get rows; !!!!!!Exec: change the data (insert, update)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []Event

	// boolean: keep loop runnning all rows
	for rows.Next() {
		var event Event

		// populate each col to the event struct
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location,
			&event.DateTime, &event.UserID) // scan each column, read content from the row

		if err != nil {
			return nil, err
		}

		// add to slice event
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	// against attack
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id) // get a single row, argument is queried id

	var event Event

	// populate the Event struct
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location,
		&event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	// sql result; error
	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)

	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

// register a event for user
func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Event) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}
