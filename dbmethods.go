package main

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Функция создания структуры и таблиц в бд
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{(*User)(nil), (*Reading)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// seedDB - Наполняем дб тестовыми данными
func seedDB(db *pg.DB) error {
	if config.Server.Debug {
		user1 := &User{
			Name:     "first",
			Surname:  "second",
			Address:  "address",
			Login:    "login",
			Password: "password",
			UserType: "user",
			Init:     false,
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
			UserType: "user",
			Init:     true,
		}
		err = db.Insert(user2)
		if err != nil {
			fmt.Println(err)
		}
		reading1 := &Reading{
			Month:    "month",
			Quantity: 100,
			UserID:   1,
			Water:    "cold",
		}
		err = db.Insert(reading1)
		if err != nil {
			fmt.Println(err)
		}
	}
	admin := &User{
		Name:     "admin",
		Login:    "admin",
		Password: "admin",
		UserType: "admin",
		Init:     false,
	}
	err := db.Insert(admin)
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

func getOnlyUsers(db *pg.DB) ([]User, error) {
	var users []User
	_, err := db.Query(&users, `SELECT * FROM users WHERE user_type = ?`, "user")
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

func getUserByName(db *pg.DB, name string) (User, error) {
	var user User
	_, err := db.QueryOne(&user, `SELECT * FROM users WHERE name = ?`, name)
	return user, err
}

func getReadings(db *pg.DB) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings`)
	return readings, err
}

func getReadingsByUserID(db *pg.DB, userID int64) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings WHERE user_id = ?`, userID)
	return readings, err
}

func getReadingsByType(db *pg.DB, water string) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings WHERE water = ?`, water)
	return readings, err
}

func getReadingByMonth(db *pg.DB, id int64, month, water string) (Reading, error) {
	var reading Reading
	_, err := db.QueryOne(&reading, `SELECT * FROM readings WHERE userid = ?, month = ?, water = ?`, id, month, water)
	return reading, err
}

func getReadingsByTypeAndOrderByQuantity(db *pg.DB, water string) ([]Reading, error) {
	var readings []Reading
	_, err := db.Query(&readings, `SELECT * FROM readings WHERE water = ? ORDER BY quantity ASC LIMIT 3`, water)
	return readings, err
}

func createUser(db *pg.DB, user *User) error {
	_, err := db.QueryOne(user, `
		INSERT INTO users (name, surname, address, login, password, init, user_type) VALUES (?name, ?surname, ?address, ?login, ?password, ?init, ?user_type)
		RETURNING id
	`, user)
	return err
}

func createReading(db *pg.DB, reading *Reading) error {
	_, err := db.QueryOne(reading, `
		INSERT INTO readings (month, quantity, user_id, water) VALUES (?month, ?quantity, ?user_id, ?water)
		RETURNING id
	`, reading)
	return err
}

func createReadingFromForm(db *pg.DB, month, water string, quantity int, userid int64) error {
	var reading Reading
	_, err := db.QueryOne(reading, `
		INSERT INTO readings (month, quantity, userid, water) VALUES (?month, ?quantity, ?userid, ?water)
		RETURNING id
	`, reading)
	return err
}

func updateUser(db *pg.DB, user *User) error {
	_, err := db.QueryOne(user, `
		UPDATE users SET (name, surname, address, login, password, init) = (?name, ?surname, ?address, ?login, ?password, ?init) WHERE id = (?id)
		RETURNING id
	`, user)
	return err
}

func updateReading(db *pg.DB, reading *Reading) error {
	_, err := db.QueryOne(reading, `
		UPDATE readings SET (month, quantity, userid, water) = (?month, ?quantity, ?userid, ?water) WHERE id = (?id)
		RETURNING id
	`, reading)
	return err
}
