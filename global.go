package main

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v2"
)

var app App
var config Config

func initDB() {
	source, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		fmt.Println(err)
	}
	app.DB = pg.Connect(&pg.Options{
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
		Addr:     config.DB.Addr,
	})
	err = createSchema(app.DB)
	if err != nil {
		fmt.Println(err)
	}
	err = seedDB(app.DB)
	if err != nil {
		fmt.Println(err)
	}
}

func initServer() {
	app.Echo = echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	app.Echo.Renderer = t
}
