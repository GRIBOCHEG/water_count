package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// validate - валидатор данных аутентификации
func validate(username, password string, c echo.Context) (bool, error) {
	user := getUserByLogin(username, DB)
	if user != nil {
		if user.Password == password {
			return true, nil
		}
	}
	return false, nil
}

// startServer - запускает сервер и регистрирует роуты
func startServer() {
	e := echo.New()

	e.Use(middleware.BasicAuth(validate))
	e.GET("/", index)
	e.GET("/seed", seed)

	e.Logger.Fatal(e.Start(":1323"))
}
