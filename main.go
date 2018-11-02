package main

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
)

// DB - Глобальная переменная содержащая подключение к БД
var DB *pg.DB

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

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Henlo!")
}

func main() {
	DB := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "1234",
	})
	defer DB.Close()
	err := createSchema(DB)
	if err != nil {
		fmt.Println(err)
	}
	err = seedDB(DB)
	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()

	e.GET("/", index)
	e.GET("/get", func(c echo.Context) error {
		users, err := getUsers(DB)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, "here "+fmt.Sprintln(users))
	})

	e.Logger.Fatal(e.Start(":1323"))
}
