// Package model keep all project related business models.
package model

// Post model represents a post.
type Post struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}
