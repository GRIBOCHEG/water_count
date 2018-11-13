package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func login(c echo.Context) error {
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("try to login with", user.Login, user.Password)
	usr, err := getUserByLogin(app.DB, user.Login)
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
			t, err := token.SignedString(app.Slice)
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
		t, err := token.SignedString(app.Slice)
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
}

func setPassword(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	usr, err := getUserByName(app.DB, name)
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
	err = updateUser(app.DB, &usr)
	if err != nil {
		fmt.Println(err)
	}
	return c.String(http.StatusOK, "done")
}

func saveReading(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	usr, err := getUserByName(app.DB, name)
	if err != nil {
		fmt.Println(err)
	}

	reading := new(Reading)
	err = c.Bind(reading)
	if err != nil {
		fmt.Println(err)
	}

	reading.UserID = usr.ID
	err = createReading(app.DB, reading)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "bad")
	}
	return c.String(http.StatusOK, "good")
}

func saveUser(c echo.Context) error {
	user := new(User)
	err := c.Bind(user)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "bad")
	}
	user.Password = user.Name + user.Surname
	user.UserType = "user"
	user.Init = true
	err = createUser(app.DB, user)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "bad")
	}
	return c.String(http.StatusOK, "good")
}

func usersList(c echo.Context) error {
	users, err := getOnlyUsers(app.DB)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	return c.Render(http.StatusOK, "users", users)
}

func readingsList(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	readings, err := getReadingsByUserID(app.DB, int64(userID))
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}
	return c.Render(http.StatusOK, "readings", readings)
}

func topConsumers(c echo.Context) error {
	var data [3]Data
	water := c.Param("water")
	readings, err := getReadingsByTypeAndOrderByQuantity(app.DB, water)
	fmt.Println(readings)
	if err != nil {
		fmt.Println(err)
	}
	for i, reading := range readings {
		user, err := getUserByID(app.DB, int(reading.UserID))
		fmt.Println(user)
		if err != nil {
			fmt.Println(err)
		}
		data[i].Rdng = reading
		data[i].Usr = user
	}
	fmt.Println(data)
	return c.Render(http.StatusOK, "consumers", data)
}
