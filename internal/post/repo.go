package post

import (
	"gorm.io/gorm"

	"github.com/imarrche/nix-ed/internal/model"
)

// repo is post repository implementation.
type repo struct {
	db *gorm.DB
}

// NewRepo creates and returns a new Repo instaace.
func NewRepo(db *gorm.DB) Repo {
	return &repo{db}
}

// GetAll gets and returns all posts.
func (r *repo) GetAll() (ps []model.Post, err error) {
	r.db.Find(&ps)

	return
}

// Create creates a post and returns it.
func (r *repo) Create(p model.Post) (model.Post, error) {
	res := r.db.Create(&p)

	return p, res.Error
}

// GetByID gets and returns the post with specifid ID.
func (r *repo) GetByID(id int) (p model.Post, err error) {
	r.db.First(&p, id)
	if p.ID != id {
		return p, ErrNotFound
	}

	return p, nil
}

// Update updates the post and returns it.
func (r *repo) Update(p model.Post) (model.Post, error) {
	r.db.Save(&p)

	return p, nil
}

// DeleteByID deletes the post with specific ID.
func (r *repo) DeleteByID(id int) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	r.db.Where("id = ?", id).Delete(model.Post{})

	return nil
}
