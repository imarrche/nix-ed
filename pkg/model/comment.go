package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
	PostID int    `json:"postId"`
}
