package main

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
)

// User - модель описывающая пользователя
type User struct {
	ID       int64  `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Surname  string `json:"surname" form:"surname"`
	Address  string `json:"address" form:"address"`
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
}

// Admin - модель описывающая администратора
type Admin struct {
	ID       int64  `json:"id" form:"id"`
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
}

// Reading - модель описывающая показания счетчика
type Reading struct {
	ID       int64  `json:"id" form:"id"`
	Month    string `json:"month" form:"month"`
	Quantity int64  `json:"quantity" form:"quantity"`
	UserID   int64  `json:"userid" form:"userid"`
	User     *User
	Type     string `json:"type" form:"type"`
}

// Функция предоставляющая строковое представление структуры пользователя
func (u *User) String() string {
	return fmt.Sprintf("User<%d %s %s %s %s %s>", u.ID, u.Name, u.Surname, u.Address, u.Login, u.Password)
}

// Функция предоставляющая строковое представление структуры показания
func (r *Reading) String() string {
	return fmt.Sprintf("Reading<%d %s %d %s %s>", r.ID, r.Month, r.Quantity, r.User, r.Type)
}

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

// ExampleDBModel - Проверка интеграции и работы ORM
func ExampleDBModel() string {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "1234",
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	user1 := &User{
		Name:     "first",
		Surname:  "second",
		Address:  "address",
		Login:    "login",
		Password: "password",
	}
	err = db.Insert(user1)
	if err != nil {
		panic(err)
	}

	err = db.Insert(&Admin{
		Login:    "admin",
		Password: "admin",
	})
	if err != nil {
		panic(err)
	}

	reading1 := &Reading{
		Month:    "month",
		Quantity: 100,
		UserID:   user1.ID,
		Type:     "cold",
	}
	err = db.Insert(reading1)
	if err != nil {
		panic(err)
	}

	// Получаем юзера по первичному ключу
	user := &User{ID: user1.ID}
	err = db.Select(user)
	if err != nil {
		panic(err)
	}

	// Получаем показание и пользователя в одном запросе
	reading := new(Reading)
	err = db.Model(reading).
		Relation("User").
		Where("reading.id = ?", reading1.ID).
		Select()
	if err != nil {
		panic(err)
	}

	return reading.String()
}

func main() {

	str := ExampleDBModel()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, str)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
