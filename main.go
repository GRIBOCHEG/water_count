package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	e.File("/index", "index.html")
	e.File("/login", "loginform.html")
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
			return c.String(http.StatusOK, "user")
		}
		return echo.ErrUnauthorized
	})

	u := e.Group("/user")
	u.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte("secret"),
		TokenLookup: "cookie:token",
	}))

	u.POST("/userform", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		usr, err := getUserByName(DB, name)
		if err != nil {
			fmt.Println(err)
		}
		month := c.FormValue("month")
		quantity := c.FormValue("quantity")
		quan, _ := strconv.Atoi(quantity)
		water := c.FormValue("water")
		reading, err := getReadingByMonth(DB, usr.ID, month, water)
		if err != nil {
			fmt.Println(err)
		}
		if reading.Quantity == 0 {
			err = createReadingFromForm(DB, month, water, quan, usr.ID)
			if err != nil {
				fmt.Println(err)
			}
			return c.String(http.StatusOK, "Показание успешно сохранено")
		}
		return c.String(http.StatusOK, "Такое показание уже есть в системе")
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
	a.GET("/userlist", func(c echo.Context) error {
		// Default golang template to show userlist
		return c.String(http.StatusOK, "Here'll be list of users")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
