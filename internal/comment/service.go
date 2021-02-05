package comment

import "github.com/imarrche/nix-ed/internal/model"

// service is comment service implementation.
type service struct {
	r Repo
}

// NewService creates and returns a new Service instance.
func NewService(r Repo) Service {
	return &service{r}
}

// GetAll gets and returns all comment.
func (s *service) GetAll() (cs []model.Comment, err error) {
	return s.r.GetAll()
}

// Create creates a comment and returns it.
func (s *service) Create(c model.Comment) (model.Comment, error) {
	if err := c.Validate(); err != nil {
		return model.Comment{}, err
	}

	return s.r.Create(c)
}

// GetByID gets and returns the comment with specific ID.
func (s *service) GetByID(id int) (c model.Comment, err error) {
	return s.r.GetByID(id)
}

// Update updates the comment and returns it.
func (s *service) Update(c model.Comment) (model.Comment, error) {
	uc, err := s.r.GetByID(c.ID)
	if err != nil {
		return model.Comment{}, err
	}

	uc.Name = c.Name
	uc.Body = c.Body
	if err := uc.Validate(); err != nil {
		return model.Comment{}, err
	}

	return s.r.Update(uc)
}

// DeleteByID deletes the comment with specific ID.
func (s *service) DeleteByID(id int) error {
	return s.r.DeleteByID(id)
}
