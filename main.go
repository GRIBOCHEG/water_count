package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	e.GET("/index", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.GET("/login", func(c echo.Context) error {
		return c.File("loginform.html")
	})
	e.POST("/login", func(c echo.Context) error {
		user := new(User)
		err := c.Bind(user)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("try to login with", user.Login, user.Password)
		usr, err := getUserByLogin(DB, user.Login)
		if err != nil {
			fmt.Println(err)
			return echo.ErrUnauthorized
		}

		if user.Password == usr.Password {
			if usr.UserType == "admin" {
				token := jwt.New(jwt.SigningMethodHS256)

				// Set claims
				claims := token.Claims.(jwt.MapClaims)
				claims["name"] = usr.Name
				claims["admin"] = true
				claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

				// Generate encoded token
				t, err := token.SignedString([]byte("secret"))
				if err != nil {
					return err
				}
				cookie := new(http.Cookie)
				cookie.Name = "token"
				cookie.Value = t
				cookie.Expires = time.Now().Add(24 * time.Hour)
				c.SetCookie(cookie)
				return c.String(http.StatusOK, "admin")
			}

			token := jwt.New(jwt.SigningMethodHS256)

			// Set claims
			claims := token.Claims.(jwt.MapClaims)
			claims["name"] = usr.Name
			claims["admin"] = false
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

			// Generate encoded token
			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return err
			}
			cookie := new(http.Cookie)
			cookie.Name = "token"
			cookie.Value = t
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)

			if usr.Init {
				return c.String(http.StatusOK, "pass")
			}

			return c.String(http.StatusOK, "user")
		}
		return echo.ErrUnauthorized
	})

	u := e.Group("/user")
	u.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
	}))

	u.GET("/changepass", func(c echo.Context) error {
		return c.File("changepass.html")
	})
	u.POST("/changepass", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		usr, err := getUserByName(DB, name)
		if err != nil {
			fmt.Println(err)
		}
		user1 := new(User)
		err = c.Bind(user1)
		if err != nil {
			fmt.Println(err)
		}
		usr.Password = user1.Password
		usr.Init = false
		err = updateUser(DB, &usr)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, "done")
	})
	u.GET("/userform", func(c echo.Context) error {
		return c.File("userform.html")
	})
	u.POST("/userform", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		usr, err := getUserByName(DB, name)
		if err != nil {
			fmt.Println(err)
		}

		reading := new(Reading)
		err = c.Bind(reading)
		if err != nil {
			fmt.Println(err)
		}

		reading1, err := getReadingByMonth(DB, usr.ID, reading.Month, reading.Water)
		if err != nil {
			fmt.Println(err)
		}
		if reading1.Quantity == reading.Quantity && reading1.Month == reading.Month && reading1.Water == reading.Water {
			return c.String(http.StatusOK, "bad")
		}

		reading.UserID = usr.ID
		err = createReading(DB, reading)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, "good")

	})

	a := e.Group("/admin")
	a.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
	}))
	a.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			isadmin := claims["admin"].(bool)
			if isadmin {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
	})
	a.GET("/adminform", func(c echo.Context) error {
		return c.File("adminform.html")
	})
	a.GET("/adduser", func(c echo.Context) error {
		return c.File("adduser.html")
	})
	a.POST("/adduser", func(c echo.Context) error {
		user := new(User)
		err := c.Bind(user)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "bad")
		}
		user.Password = user.Name + user.Surname
		user.UserType = "user"
		user.Init = true
		err = createUser(DB, user)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "bad")
		}
		return c.String(http.StatusOK, "good")
	})
	a.GET("/userlist", func(c echo.Context) error {
		// Default golang template to show userlist
		return c.String(http.StatusOK, "Here'll be list of users")
	})
	a.GET("/statistics", func(c echo.Context) error {
		return c.String(http.StatusOK, "Here'll be statistics")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
