package models

import "errors"

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

var Error bool

type User struct {
	ID       int
	Email    string
	Username string
	Password string
}

type Post struct {
	ID      int
	User_id int
	Title   string
	Content string
}
