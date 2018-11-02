package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func index(c echo.Context) error {
	return c.String(http.StatusOK, "Henlo!")
}

func seed(c echo.Context) error {
	err := seedDB(DB)
	if err != nil {
		panic(err)
	}
	return c.String(http.StatusOK, "DB seeded!")
}
