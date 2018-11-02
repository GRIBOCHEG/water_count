package main

import (
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

// connectDB - Подключаемся к бд
func connectDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "1234",
	})
	defer db.Close()

	// err := createSchema(db)
	// if err != nil {
	// 	return nil
	// }

	return db
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

	return err

	// Получаем юзера по первичному ключу
	// user := &User{ID: user1.ID}
	// err = db.Select(user)
	// if err != nil {
	// 	panic(err)
	// }

	// Получаем показание и пользователя в одном запросе
	// reading := new(Reading)
	// err = db.Model(reading).
	// 	Relation("User").
	// 	Where("reading.id = ?", reading1.ID).
	// 	Select()
	// if err != nil {
	// 	panic(err)
	// }
}

func getUserByLogin(login string, db *pg.DB) *User {
	user := &User{Login: login}
	err := db.Select(user)
	if err != nil {
		panic(err)
	}
	if user.ID != 0 {
		return user
	}
	return nil
}
