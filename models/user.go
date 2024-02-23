package models

import (
	"errors"

	"example.com/eveny-booking/db"
	"example.com/eveny-booking/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword) // put hashed password into DB
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId() // get the id of user inserted

	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email) // only select ONE row, inject email to ?

	var retrievedPassword string

	// populate u.ID and retrievedPassword var with data (id, hashed password from DB)
	err := row.Scan(&u.ID, &retrievedPassword) // error if no line match the query (did not find anything)
	if err != nil {
		return errors.New("User not found!")
	}

	// find the user, compare the password
	// 1. plain password, 2.hashed password
	passwordIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !passwordIsValid {
		return errors.New("Incorrect Password!")
	}

	// credential is valid
	return nil // everything worked
}
