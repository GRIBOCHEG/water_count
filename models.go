package main

import "fmt"

// User - модель описывающая пользователя sql:"unique" sql:"pk"
type User struct {
	ID       int64  `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Surname  string `json:"surname" form:"surname"`
	Address  string `json:"address" form:"address"`
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
	UserType string `json:"usertype" form:"usertype"`
	Init     bool   `json:"init" form:"init"`
}

// Reading - модель описывающая показания счетчика
type Reading struct {
	ID       int64  `json:"id" form:"id"`
	Month    string `json:"month" form:"month"`
	Quantity int64  `json:"quantity" form:"quantity"`
	UserID   int64  `json:"userid" form:"userid"`
	Water    string `json:"water" form:"water"`
}

// Функция предоставляющая строковое представление структуры пользователя
func (u *User) String() string {
	return fmt.Sprintf("User<%d %s %s %s %s %s>", u.ID, u.Name, u.Surname, u.Address, u.Login, u.Password)
}

// Функция предоставляющая строковое представление структуры показания
func (r *Reading) String() string {
	return fmt.Sprintf("Reading<%d %s %d %d %s>", r.ID, r.Month, r.Quantity, r.UserID, r.Water)
}
