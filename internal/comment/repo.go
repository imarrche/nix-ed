package comment

import (
	"gorm.io/gorm"

	"github.com/imarrche/nix-ed/internal/model"
)

// repo is comment repository implementation.
type repo struct {
	db *gorm.DB
}

// NewRepo creates and returns a new Repo instaace.
func NewRepo(db *gorm.DB) Repo {
	return &repo{db}
}

// GetAll gets and returns all comments.
func (r *repo) GetAll() (cs []model.Comment, err error) {
	r.db.Find(&cs)

	return
}

// Create creates a comment and returns it.
func (r *repo) Create(c model.Comment) (model.Comment, error) {
	res := r.db.Create(&c)

	return c, res.Error
}

// GetByID gets and returns the comment with specifid ID.
func (r *repo) GetByID(id int) (c model.Comment, err error) {
	r.db.First(&c, id)
	if c.ID != id {
		return c, ErrNotFound
	}

	return c, nil
}

// Update updates the comment and returns it.
func (r *repo) Update(c model.Comment) (model.Comment, error) {
	r.db.Save(&c)

	return c, nil
}

// DeleteByID deletes the comment with specific ID.
func (r *repo) DeleteByID(id int) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}
	r.db.Where("id = ?", id).Delete(model.Comment{})

	return nil
}
