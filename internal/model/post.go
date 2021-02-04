// Package model keep all project related business models.
package model

import validation "github.com/go-ozzo/ozzo-validation"

// Post model represents a post.
type Post struct {
	ID     int    `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

// Validate validates post's fields.
func (p *Post) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Body, validation.Required),
		validation.Field(&p.UserID, validation.Required),
	)
}
