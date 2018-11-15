package main

import (
	"encoding/binary"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v2"
)

var app App
var config Config

//Считывание конфиг файла, и формирование секрета для дальнейшей генерации токена
func initConfig() error {
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return err
	}

	if !config.Server.Debug {
		stamp := make([]byte, binary.MaxVarintLen64)
		_ = binary.PutVarint(stamp, time.Now().UnixNano())
		app.Slice = append(stamp, []byte(config.Server.Secret)...)
	} else {
		app.Slice = []byte("secret")
	}
	return nil
}

func pingDB() error {
	_, err := app.DB.ExecOne("SELECT 1")
	return err
}

func initDB() error {
	app.DB = pg.Connect(&pg.Options{
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
		Addr:     config.DB.Addr,
	})

	err := pingDB()
	if err != nil {
		return err
	}

	createSchema(app.DB)
	seedDB(app.DB)
	return nil
}

func initServer() error {
	app.Echo = echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	app.Echo.Renderer = t
	return nil
}
