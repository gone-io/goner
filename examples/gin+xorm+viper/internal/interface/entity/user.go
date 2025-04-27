package entity

import "time"

type User struct {
	Id       int64  `json:"id" xorm:"pk autoincr"`
	Username string `json:"username"`
	Password string `json:"-"`

	CreatedAt *time.Time `json:"-" xorm:"created"`
	UpdatedAt *time.Time `json:"-" xorm:"updated"`
	DeletedAt *time.Time `json:"-" xorm:"deleted"`
}

type LoginParam struct {
	Username string `json:"username" bind:"required"`
	Password string `json:"password" bind:"required"`
}

type LoginResult struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type RegisterParam struct {
	Username string `json:"username" bind:"required"`
	Password string `json:"password" bind:"required"`
}
