package models

import "time"

type User struct {
	Id         int64 `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Address string `json:"address" db:"address"`
	Phone   string `json:"phone" db:"phone"`
	CreatedAt time.Time	`json:"create_at" db:"create_at"`
}

func NewUser(name, address, phone string) User {
	return User {
		Name: name,
		Address: address,
		Phone: phone,
		CreatedAt: time.Now(),
	}
}
