package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Comment model represents a post's comment.
type Comment struct {
	ID     int    `json:"id" xml:"id" gorm:"primaryKey"`
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
	PostID int    `json:"postId" xml:"postId"`
}

// Validate validates comment's fields.
func (c *Comment) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Body, validation.Required),
		validation.Field(&c.PostID, validation.Required),
	)
}
