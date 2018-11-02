package main

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
)

// DB - Глобальная переменная содержащая подключение к БД
var DB *pg.DB

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
	e.GET("/getall", func(c echo.Context) error {
		users, err := getUsers(DB)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, "here "+fmt.Sprintln(users))
	})
	e.File("/index", "index.html")
	e.File("/login", "loginform.html")
	e.File("/userform", "userform.html")
	e.POST("/login", func(c echo.Context) error {
		login := c.FormValue("login")
		password := c.FormValue("password")
		users, err := getUsers(DB)
		if err != nil {
			fmt.Println(err)
		}
		for _, usr := range users {
			if login == usr.Login && password == usr.Password {
				return c.Redirect(http.StatusOK, "/userform")
			}
		}
		return c.Redirect(http.StatusOK, "/index")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
