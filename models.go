package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/go-pg/pg"

	"github.com/labstack/echo"
)

// User - модель описывающая пользователя sql:"unique" sql:"pk"
type User struct {
	ID       int64  `sql:",pk" json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Surname  string `json:"surname" form:"surname"`
	Address  string `json:"address" form:"address"`
	Login    string `sql:",unique" json:"login" form:"login"`
	Password string `json:"password" form:"password"`
	UserType string `json:"usertype" form:"usertype"`
	Init     bool   `json:"init" form:"init"`
}

// Reading - модель описывающая показания счетчика
type Reading struct {
	ID       int64  `sql:",pk" json:"id" form:"id"`
	Month    string `sql:"unique:user_month_water" json:"month" form:"month"`
	Quantity int64  `json:"quantity,,string" form:"quantity"`
	UserID   int64  `sql:"unique:user_month_water,,notnull" json:"userid,,string" form:"userid"`
	Water    string `sql:"unique:user_month_water" json:"water" form:"water"`
}

type Data struct {
	Rdng Reading
	Usr  User
}

type Template struct {
	templates *template.Template
}

type App struct {
	Echo  *echo.Echo
	DB    *pg.DB
	Slice []byte
}

type Config struct {
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Server struct {
	Debug  bool   `yaml:"debug"`
	Port   string `yaml:"port"`
	Secret string `yaml:"secret"`
}

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Addr     string `yaml:"addr"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Функция предоставляющая строковое представление структуры пользователя
func (u *User) String() string {
	return fmt.Sprintf("User<%d %s %s %s %s %s>", u.ID, u.Name, u.Surname, u.Address, u.Login, u.Password)
}

// Функция предоставляющая строковое представление структуры показания
func (r *Reading) String() string {
	return fmt.Sprintf("Reading<%d %s %d %d %s>", r.ID, r.Month, r.Quantity, r.UserID, r.Water)
}
