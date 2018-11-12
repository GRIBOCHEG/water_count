package main

import "fmt"

type Data struct {
	Reading *Reading
	User    *User
}

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

// Функция предоставляющая строковое представление структуры пользователя
func (u *User) String() string {
	return fmt.Sprintf("User<%d %s %s %s %s %s>", u.ID, u.Name, u.Surname, u.Address, u.Login, u.Password)
}

// Функция предоставляющая строковое представление структуры показания
func (r *Reading) String() string {
	return fmt.Sprintf("Reading<%d %s %d %d %s>", r.ID, r.Month, r.Quantity, r.UserID, r.Water)
}
