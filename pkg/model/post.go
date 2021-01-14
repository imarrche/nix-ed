package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}
