package post

import "github.com/imarrche/nix-ed/internal/model"

// service is post service implementation.
type service struct {
	r Repo
}

// NewService creates and returns a new Service instance.
func NewService(r Repo) Service {
	return &service{r}
}

// GetAll gets and returns all posts.
func (s *service) GetAll() (ps []model.Post, err error) {
	return s.r.GetAll()
}

// Create creates a post and returns it.
func (s *service) Create(p model.Post) (model.Post, error) {
	if err := p.Validate(); err != nil {
		return model.Post{}, err
	}

	return s.r.Create(p)
}

// GetByID gets and returns the post with specific ID.
func (s *service) GetByID(id int) (p model.Post, err error) {
	return s.r.GetByID(id)
}

// Update updates the post and returns it.
func (s *service) Update(p model.Post) (model.Post, error) {
	up, err := s.r.GetByID(p.ID)
	if err != nil {
		return model.Post{}, err
	}

	up.Title = p.Title
	up.Body = p.Body
	if err := up.Validate(); err != nil {
		return model.Post{}, err
	}

	return s.r.Update(up)
}

// DeleteByID deletes the post with specific ID.
func (s *service) DeleteByID(id int) error {
	return s.r.DeleteByID(id)
}
