package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Функция создания структуры и таблиц в бд
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Admin)(nil), (*Reading)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// seedDB - Наполняем дб тестовыми данными
func seedDB(db *pg.DB) error {
	user1 := &User{
		Name:     "first",
		Surname:  "second",
		Address:  "address",
		Login:    "login",
		Password: "password",
	}
	err := db.Insert(user1)
	if err != nil {
		fmt.Println(err)
	}

	user2 := &User{
		Name:     "name",
		Surname:  "surname",
		Address:  "address2",
		Login:    "username",
		Password: "namesurname",
	}
	err = db.Insert(user2)
	if err != nil {
		fmt.Println(err)
	}

	err = db.Insert(&Admin{
		Login:    "admin",
		Password: "admin",
	})
	if err != nil {
		fmt.Println(err)
	}

	reading1 := &Reading{
		Month:    "month",
		Quantity: 100,
		UserID:   user1.ID,
		Type:     "cold",
	}
	err = db.Insert(reading1)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func getUsers(db *pg.DB) ([]User, error) {
	var users []User
	_, err := db.Query(&users, `SELECT * FROM users`)
	return users, err
}

func getUserByLogin(db *pg.DB, login string) (User, error) {
	var user User
	_, err := db.QueryOne(&user, `SELECT * FROM users WHERE login = ?`, login)
	return user, err
}

func getUserByID(db *pg.DB, id int) (User, error) {
	var user User
	_, err := db.QueryOne(&user, `SELECT * FROM users WHERE id = ?`, id)
	return user, err
}

func getReadings(db *pg.DB) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings`)
	return readings, err
}

func getReadingsByUserID(db *pg.DB, id int) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings WHERE userid = ?`, id)
	return readings, err
}

func getReadingsByType(db *pg.DB, typ string) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings WHERE type = ?`, typ)
	return readings, err
}

func createUser(db *pg.DB, user *User) error {
	_, err := db.QueryOne(user, `
		INSERT INTO users (name, surname, address, login, password) VALUES (?name, ?surname, ?address, ?login, ?password)
		RETURNING id
	`, user)
	return err
}

func createReading(db *pg.DB, reading *Reading) error {
	_, err := db.QueryOne(reading, `
		INSERT INTO readings (month, quantity, userid, type) VALUES (?month, ?quantity, ?userid, ?type)
		RETURNING id
	`, reading)
	return err
}

func updateUser(db *pg.DB, user *User) error {
	_, err := db.QueryOne(user, `
		UPDATE users SET (name, surname, address, login, password) VALUES (?name, ?surname, ?address, ?login, ?password)
		RETURNING id
	`, user)
	return err
}

func updateReading(db *pg.DB, reading *Reading) error {
	_, err := db.QueryOne(reading, `
		UPDATE readings SET (month, quantity, userid, type) VALUES (?month, ?quantity, ?userid, ?type)
		RETURNING id
	`, reading)
	return err
}
