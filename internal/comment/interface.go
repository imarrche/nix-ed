// Package comment provides all comment domain related logic.
package comment

import "github.com/imarrche/nix-ed/internal/model"

//go:generate mockgen -source=interface.go -destination=mock/mock.go

// Repo is the interface all comment repositories must implement.
type Repo interface {
	GetAll() ([]model.Comment, error)
	Create(model.Comment) (model.Comment, error)
	GetByID(int) (model.Comment, error)
	Update(model.Comment) (model.Comment, error)
	DeleteByID(int) error
}

// Service is the interface all comment services must implement.
type Service interface {
	GetAll() ([]model.Comment, error)
	Create(model.Comment) (model.Comment, error)
	GetByID(int) (model.Comment, error)
	Update(model.Comment) (model.Comment, error)
	DeleteByID(int) error
}
