// Package post provides all post domain related logic.
package post

import "github.com/imarrche/nix-ed/internal/model"

//go:generate mockgen -source=interface.go -destination=mock/mock.go

// Repo is the interface all post repositories must implement.
type Repo interface {
	GetAll() ([]model.Post, error)
	Create(model.Post) (model.Post, error)
	GetByID(int) (model.Post, error)
	Update(model.Post) (model.Post, error)
	DeleteByID(int) error
}

// Service is the interface all post services must implement.
type Service interface {
	GetAll() ([]model.Post, error)
	Create(model.Post) (model.Post, error)
	GetByID(int) (model.Post, error)
	Update(model.Post) (model.Post, error)
	DeleteByID(int) error
}
