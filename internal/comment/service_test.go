package comment

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mockcomment "github.com/imarrche/nix-ed/internal/comment/mock"
	"github.com/imarrche/nix-ed/internal/model"
)

func TestCommentService_GetAll(t *testing.T) {
	testcases := []struct {
		name        string
		mock        func(*mockcomment.MockRepo, []model.Comment)
		comments    []model.Comment
		expComments []model.Comment
		expError    error
	}{
		{
			name: "comments are retrieved",
			mock: func(r *mockcomment.MockRepo, cs []model.Comment) {
				r.EXPECT().GetAll().Return(cs, nil)

			},
			comments:    []model.Comment{{Name: "Comment 1"}, {Name: "Comment 1."}},
			expComments: []model.Comment{{Name: "Comment 1"}, {Name: "Comment 1."}},
			expError:    nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockcomment.NewMockRepo(c)
			tc.mock(repo, tc.comments)
			s := NewService(repo)

			cs, err := s.GetAll()

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expComments, cs)
		})
	}
}

func TestCommentService_Create(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockRepo, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expError   error
	}{
		{
			name: "comment is created",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().Create(cm).Return(cm, nil)
			},
			comment:    model.Comment{Name: "Comment 1", Email: "u@t.com", Body: "Body.", PostID: 1},
			expComment: model.Comment{Name: "Comment 1", Email: "u@t.com", Body: "Body.", PostID: 1},
			expError:   nil,
		},
		{
			name:     "validation errors",
			mock:     func(_ *mockcomment.MockRepo, _ model.Comment) {},
			comment:  model.Comment{Name: "Title 1", Email: "u@t.com", PostID: 1},
			expError: validation.Errors{"body": errors.New("cannot be blank")},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockcomment.NewMockRepo(c)
			tc.mock(repo, tc.comment)
			s := NewService(repo)

			cm, err := s.Create(tc.comment)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expComment, cm)
		})
	}
}

func TestCommentService_GetByID(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockRepo, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expError   error
	}{
		{
			name: "comment is retrieved by ID",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().GetByID(cm.ID).Return(cm, nil)
			},
			comment:    model.Comment{Name: "Comment 1"},
			expComment: model.Comment{Name: "Comment 1"},
			expError:   nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockcomment.NewMockRepo(c)
			tc.mock(repo, tc.comment)
			s := NewService(repo)

			cm, err := s.GetByID(tc.comment.ID)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expComment, cm)
		})
	}
}

func TestCommentService_Update(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockRepo, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expError   error
	}{
		{
			name: "comment is updated",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().GetByID(cm.ID).Return(cm, nil)
				r.EXPECT().Update(cm).Return(cm, nil)
			},
			comment:    model.Comment{Name: "Comment 1", Email: "u@t.com", Body: "Body.", PostID: 1},
			expComment: model.Comment{Name: "Comment 1", Email: "u@t.com", Body: "Body.", PostID: 1},
			expError:   nil,
		},
		{
			name: "validation errors",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().GetByID(cm.ID).Return(cm, nil)
			},
			comment:  model.Comment{Name: "Comment 1", Email: "u@t.com", PostID: 1},
			expError: validation.Errors{"body": errors.New("cannot be blank")},
		},
		{
			name: "comment not found",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().GetByID(cm.ID).Return(model.Comment{}, errors.New("not found"))
			},
			comment:  model.Comment{Name: "Comment 1"},
			expError: errors.New("not found"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockcomment.NewMockRepo(c)
			tc.mock(repo, tc.comment)
			s := NewService(repo)

			cm, err := s.Update(tc.comment)

			assert.Equal(t, tc.expError, err)
			assert.Equal(t, tc.expComment, cm)
		})
	}
}

func TestCommentService_DeleteByID(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockcomment.MockRepo, model.Comment)
		comment  model.Comment
		expError error
	}{
		{
			name: "comment is deleted by ID",
			mock: func(r *mockcomment.MockRepo, cm model.Comment) {
				r.EXPECT().DeleteByID(cm.ID).Return(nil)
			},
			comment:  model.Comment{Name: "Comment 1"},
			expError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mockcomment.NewMockRepo(c)
			tc.mock(repo, tc.comment)
			s := NewService(repo)

			err := s.DeleteByID(tc.comment.ID)

			assert.Equal(t, tc.expError, err)
		})
	}
}
