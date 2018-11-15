package main

import (
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func startServer() {
	router()
	if !config.Server.Debug {
		port := os.Getenv("PORT")
		if port == "" {
			log.Fatal("$PORT must be set")
		}
		app.Echo.Use(middleware.Recover())
		app.Echo.Logger.Fatal(app.Echo.Start(":" + port))
	} else {
		app.Echo.Use(middleware.Recover())
		app.Echo.Logger.Fatal(app.Echo.Start(":" + config.Server.Port))
	}
}

// middleware-функция, проверяющая привелегии пользователя
func middlewarePrivelegeChecker(admin bool) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			isadmin := claims["admin"].(bool)
			if isadmin == admin {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	}
}

func router() {

	app.Echo.Static("/static", "assets")

	app.Echo.GET("/", func(c echo.Context) error {
		return c.File("views/index.html")
	})
	app.Echo.GET("/login", func(c echo.Context) error {
		return c.File("views/loginform.html")
	})
	app.Echo.POST("/login", login)

	u := app.Echo.Group("/user")
	u.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  app.Slice,
		TokenLookup: "cookie:token",
	}))

	u.Use(middlewarePrivelegeChecker(false))

	u.GET("/changepass", func(c echo.Context) error {
		return c.File("views/changepass.html")
	})
	u.POST("/changepass", setPassword)
	u.GET("/userform", func(c echo.Context) error {
		return c.File("views/userform.html")
	})
	u.POST("/userform", saveReading)

	a := app.Echo.Group("/admin")
	a.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  app.Slice,
		TokenLookup: "cookie:token",
	}))

	a.Use(middlewarePrivelegeChecker(true))

	a.GET("/adminform", func(c echo.Context) error {
		return c.File("views/adminform.html")
	})
	a.GET("/adduser", func(c echo.Context) error {
		return c.File("views/adduser.html")
	})
	a.POST("/adduser", saveUser)
	a.GET("/userlist", usersList)
	a.GET("/readinglist/:id", readingsList)
	a.GET("/statistics", func(c echo.Context) error {
		return c.File("views/statistics.html")
	})
	a.GET("/consumers/:water", topConsumers)
}
