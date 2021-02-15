package post

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/imarrche/nix-ed/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	mockauth "github.com/imarrche/nix-ed/internal/auth/mock"
	mockpost "github.com/imarrche/nix-ed/internal/post/mock"
)

func TestHandler_Auth(t *testing.T) {
	next := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	testcases := []struct {
		name    string
		mock    func(*mockauth.MockService)
		expCode int
	}{
		{
			name: "user is authenticated",
			mock: func(s *mockauth.MockService) {
				data := []byte(`{"id":"1","email":"u@t.com"}`)
				s.EXPECT().GetUserInfo("token").Return(data, nil)
			},
			expCode: http.StatusOK,
		},
		{
			name: "get user info error",
			mock: func(s *mockauth.MockService) {
				s.EXPECT().GetUserInfo("token").Return(nil, errors.New("internal error"))
			},
			expCode: http.StatusTemporaryRedirect,
		},
		{
			name: "invalid auth data",
			mock: func(s *mockauth.MockService) {
				data := []byte(`{"id":1,"email":"u@t.com"}`)
				s.EXPECT().GetUserInfo("token").Return(data, nil)
			},
			expCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		as := mockauth.NewMockService(c)
		tc.mock(as)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/posts", nil)
		r.Header.Add("Authorization", "token")

		ctx := echo.New().NewContext(r, w)

		hf := NewHandler(nil, as).Auth(next)
		hf(ctx)

		assert.Equal(t, tc.expCode, w.Code)

	}
}

func TestHandler_PostAuthor(t *testing.T) {
	next := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}

	testcases := []struct {
		name    string
		mock    func(*mockpost.MockService, model.Post)
		post    model.Post
		expCode int
	}{
		{
			name: "user is post author",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(p, nil)
			},
			post:    model.Post{ID: 1, Title: "Post 1", UserID: "1"},
			expCode: http.StatusOK,
		},
		{
			name: "post not found",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(model.Post{}, ErrNotFound)
			},
			post:    model.Post{ID: 1, Title: "Post 1", UserID: "1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(model.Post{}, errors.New("internal error"))
			},
			post:    model.Post{ID: 1, Title: "Post 1", UserID: "1"},
			expCode: http.StatusInternalServerError,
		},
		{
			name: "user is not a post author",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(p, nil)
			},
			post:    model.Post{ID: 1, Title: "Post 1", UserID: "2"},
			expCode: http.StatusForbidden,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		tc.mock(ps, tc.post)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPatch, "/posts/1", nil)
		r = r.WithContext(context.WithValue(r.Context(), uIDkey, "1"))

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		hf := NewHandler(ps, nil).PostAuthor(next)
		hf(ctx)

		assert.Equal(t, tc.expCode, w.Code)

	}
}

func TestHandler_GetAll(t *testing.T) {
	testcases := []struct {
		name     string
		mock     func(*mockpost.MockService, []model.Post)
		posts    []model.Post
		expPosts []model.Post
		expCode  int
	}{
		{
			name: "posts are retrieved",
			mock: func(s *mockpost.MockService, ps []model.Post) {
				s.EXPECT().GetAll().Return(ps, nil)
			},
			posts:    []model.Post{{Title: "Post1"}, {Title: "Post2"}},
			expPosts: []model.Post{{Title: "Post1"}, {Title: "Post2"}},
			expCode:  http.StatusOK,
		},
		{
			name: "internal error",
			mock: func(s *mockpost.MockService, ps []model.Post) {
				s.EXPECT().GetAll().Return(nil, errors.New("internal error"))
			},
			posts:   []model.Post{{Title: "Post1"}, {Title: "Post2"}},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.posts)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/posts", nil)

		ctx := echo.New().NewContext(r, w)

		NewHandler(ps, as).GetAll(ctx)

		var posts []model.Post
		json.NewDecoder(w.Body).Decode(&posts)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expPosts, posts)
	}
}

func TestHandler_Create(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(*mockpost.MockService, model.Post)
		post    model.Post
		expPost model.Post
		expCode int
	}{
		{
			name: "post is created",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().Create(p).Return(p, nil)
			},
			post:    model.Post{Title: "Post 1", UserID: "1"},
			expPost: model.Post{Title: "Post 1", UserID: "1"},
			expCode: http.StatusCreated,
		},
		{
			name: "post creating error",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().Create(p).Return(model.Post{}, errors.New("internal error"))
			},
			post:    model.Post{Title: "Post 1", UserID: "1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.post)
		w := httptest.NewRecorder()

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.post)
		r := httptest.NewRequest(http.MethodPost, "/posts", b)
		r = r.WithContext(context.WithValue(r.Context(), uIDkey, tc.post.UserID))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := echo.New().NewContext(r, w)

		NewHandler(ps, as).Create(ctx)

		var p model.Post
		json.NewDecoder(w.Body).Decode(&p)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expPost, p)
	}
}

func TestHandler_GetByID(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(*mockpost.MockService, model.Post)
		post    model.Post
		expPost model.Post
		expCode int
	}{
		{
			name: "post is retrieved",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(p, nil)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expPost: model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusOK,
		},
		{
			name: "post is not found",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(model.Post{}, ErrNotFound)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().GetByID(p.ID).Return(model.Post{}, errors.New("internal error"))
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.post)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/posts/1", nil)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).GetByID(ctx)

		var post model.Post
		json.NewDecoder(w.Body).Decode(&post)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expPost, post)
	}
}

func TestHandler_Update(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(*mockpost.MockService, model.Post)
		post    model.Post
		expPost model.Post
		expCode int
	}{
		{
			name: "post is updated",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().Update(p).Return(p, nil)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expPost: model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusOK,
		},
		{
			name: "post is not found",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().Update(p).Return(model.Post{}, ErrNotFound)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().Update(p).Return(model.Post{}, errors.New("internal error"))
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.post)
		w := httptest.NewRecorder()
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.post)
		r := httptest.NewRequest(http.MethodPatch, "/posts/1", b)
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).Update(ctx)

		var post model.Post
		json.NewDecoder(w.Body).Decode(&post)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expPost, post)
	}
}

func TestHandler_DeleteyByID(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(*mockpost.MockService, model.Post)
		post    model.Post
		expCode int
	}{
		{
			name: "post is deleted",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().DeleteByID(p.ID).Return(nil)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusNoContent,
		},
		{
			name: "post is not found",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().DeleteByID(p.ID).Return(ErrNotFound)
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockpost.MockService, p model.Post) {
				s.EXPECT().DeleteByID(p.ID).Return(errors.New("internal error"))
			},
			post:    model.Post{ID: 1, Title: "Post1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockpost.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.post)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/posts/1", nil)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).DeleteByID(ctx)

		assert.Equal(t, tc.expCode, w.Code)
	}
}
