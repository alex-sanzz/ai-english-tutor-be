package models

import "time"

type User struct {
	ID         string
	GoogleSub  string
	Email      string
	Name       string
	CreatedAt time.Time
	LastLoginAt time.Time
}