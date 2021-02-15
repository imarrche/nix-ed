package post

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/imarrche/nix-ed/internal/model"
	mockpost "github.com/imarrche/nix-ed/internal/post/mock"
)

func TestPostService_GetAll(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockRepo, []model.Post)
		posts    []model.Post
		expPosts []model.Post
		expError error
	}{
		{
			name: "posts are retrieved",
			mock: func(r *mockpost.MockRepo, ps []model.Post) {
				r.EXPECT().GetAll().Return(ps, nil)

			},
			posts: []model.Post{
				{Title: "Title 1"}, {Title: "Title 2"},
			},
			expPosts: []model.Post{
				{Title: "Title 1"}, {Title: "Title 2"},
			},
			expError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockpost.NewMockRepo(c)
			tc.mock(repo, tc.posts)
			s := NewService(repo)

			ps, err := s.GetAll()

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expPosts, ps)
		})
	}
}

func TestPostService_Create(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockRepo, model.Post)
		post     model.Post
		expPost  model.Post
		expError error
	}{
		{
			name: "posts is created",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().Create(p).Return(p, nil)
			},
			post:     model.Post{Title: "Title 1", Body: "Body.", UserID: "1"},
			expPost:  model.Post{Title: "Title 1", Body: "Body.", UserID: "1"},
			expError: nil,
		},
		{
			name:     "validation errors",
			mock:     func(_ *mockpost.MockRepo, _ model.Post) {},
			post:     model.Post{Title: "Title 1", UserID: "1"},
			expError: validation.Errors{"body": errors.New("cannot be blank")},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockpost.NewMockRepo(c)
			tc.mock(repo, tc.post)
			s := NewService(repo)

			p, err := s.Create(tc.post)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expPost, p)
		})
	}
}

func TestPostService_GetByID(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockRepo, model.Post)
		post     model.Post
		expPost  model.Post
		expError error
	}{
		{
			name: "post is retrieved by ID",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().GetByID(p.ID).Return(p, nil)
			},
			post:     model.Post{Title: "Title 1"},
			expPost:  model.Post{Title: "Title 1"},
			expError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockpost.NewMockRepo(c)
			tc.mock(repo, tc.post)
			s := NewService(repo)

			p, err := s.GetByID(tc.post.ID)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expPost, p)
		})
	}
}

func TestPostService_Update(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockRepo, model.Post)
		post     model.Post
		expPost  model.Post
		expError error
	}{
		{
			name: "post is updated",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().GetByID(p.ID).Return(p, nil)
				r.EXPECT().Update(p).Return(p, nil)
			},
			post:     model.Post{Title: "Updated title 1", Body: "Body.", UserID: "1"},
			expPost:  model.Post{Title: "Updated title 1", Body: "Body.", UserID: "1"},
			expError: nil,
		},
		{
			name: "validation errors",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().GetByID(p.ID).Return(p, nil)
			},
			post:     model.Post{Title: "Title 1", UserID: "1"},
			expError: validation.Errors{"body": errors.New("cannot be blank")},
		},
		{
			name: "post not found",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().GetByID(p.ID).Return(model.Post{}, errors.New("not found"))
			},
			post:     model.Post{Title: "Title 1", UserID: "1"},
			expError: errors.New("not found"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockpost.NewMockRepo(c)
			tc.mock(repo, tc.post)
			s := NewService(repo)

			p, err := s.Update(tc.post)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expPost, p)
		})
	}
}

func TestPostService_DeleteByID(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockRepo, model.Post)
		post     model.Post
		expError error
	}{
		{
			name: "post is deleted by ID",
			mock: func(r *mockpost.MockRepo, p model.Post) {
				r.EXPECT().DeleteByID(p.ID).Return(nil)
			},
			post:     model.Post{Title: "Title 1"},
			expError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockpost.NewMockRepo(c)
			tc.mock(repo, tc.post)
			s := NewService(repo)

			err := s.DeleteByID(tc.post.ID)

			assert.Equal(t, tc.expError, err)
		})
	}
}
