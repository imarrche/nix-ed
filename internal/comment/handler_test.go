package comment

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
	mockcomment "github.com/imarrche/nix-ed/internal/comment/mock"
)

func TestHandler_GetAll(t *testing.T) {
	testcases := []struct {
		name        string
		mock        func(*mockcomment.MockService, []model.Comment)
		comments    []model.Comment
		expComments []model.Comment
		expCode     int
	}{
		{
			name: "comment are retrieved",
			mock: func(s *mockcomment.MockService, cs []model.Comment) {
				s.EXPECT().GetAll().Return(cs, nil)
			},
			comments:    []model.Comment{{Body: "Comment 1"}, {Body: "Comment 2"}},
			expComments: []model.Comment{{Body: "Comment 1"}, {Body: "Comment 2"}},
			expCode:     http.StatusOK,
		},
		{
			name: "internal error",
			mock: func(s *mockcomment.MockService, cs []model.Comment) {
				s.EXPECT().GetAll().Return(nil, errors.New("internal error"))
			},
			comments: []model.Comment{{Body: "Comment 1"}, {Body: "Comment 2"}},
			expCode:  http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		cs := mockcomment.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(cs, tc.comments)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/posts", nil)

		ctx := echo.New().NewContext(r, w)

		NewHandler(cs, as).GetAll(ctx)

		var comments []model.Comment
		json.NewDecoder(w.Body).Decode(&comments)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expComments, comments)
	}
}

func TestHandler_Create(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockService, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expCode    int
	}{
		{
			name: "comment is created",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().Create(cm).Return(cm, nil)
			},
			comment:    model.Comment{Email: "u@t.com", Body: "Comment 1", PostID: 1},
			expComment: model.Comment{Email: "u@t.com", Body: "Comment 1", PostID: 1},
			expCode:    http.StatusCreated,
		},
		{
			name: "comment creating error",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().Create(cm).Return(model.Comment{}, errors.New("internal error"))
			},
			comment: model.Comment{Email: "u@t.com", Body: "Comment 1", PostID: 1},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockcomment.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.comment)
		w := httptest.NewRecorder()

		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.comment)
		r := httptest.NewRequest(http.MethodPost, "/comments", b)
		r = r.WithContext(context.WithValue(r.Context(), uEmailKey, tc.comment.Email))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := echo.New().NewContext(r, w)

		NewHandler(ps, as).Create(ctx)

		var cm model.Comment
		json.NewDecoder(w.Body).Decode(&cm)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expComment, cm)
	}
}

func TestHandler_GetByID(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockService, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expCode    int
	}{
		{
			name: "comment is retrieved",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().GetByID(cm.ID).Return(cm, nil)
			},
			comment:    model.Comment{ID: 1, Body: "Comment 1"},
			expComment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode:    http.StatusOK,
		},
		{
			name: "comment is not found",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().GetByID(cm.ID).Return(model.Comment{}, ErrNotFound)
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().GetByID(cm.ID).Return(model.Comment{}, errors.New("internal error"))
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockcomment.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.comment)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/comments/1", nil)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).GetByID(ctx)

		var cm model.Comment
		json.NewDecoder(w.Body).Decode(&cm)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expComment, cm)
	}
}

func TestHandler_Update(t *testing.T) {
	testcases := []struct {
		name       string
		mock       func(*mockcomment.MockService, model.Comment)
		comment    model.Comment
		expComment model.Comment
		expCode    int
	}{
		{
			name: "comment is updated",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().Update(cm).Return(cm, nil)
			},
			comment:    model.Comment{ID: 1, Body: "Comment 1"},
			expComment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode:    http.StatusOK,
		},
		{
			name: "comment is not found",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().Update(cm).Return(model.Comment{}, ErrNotFound)
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().Update(cm).Return(model.Comment{}, errors.New("internal error"))
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockcomment.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.comment)
		w := httptest.NewRecorder()
		b := &bytes.Buffer{}
		json.NewEncoder(b).Encode(tc.comment)
		r := httptest.NewRequest(http.MethodPatch, "/comment/1", b)
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).Update(ctx)

		var cm model.Comment
		json.NewDecoder(w.Body).Decode(&cm)

		assert.Equal(t, tc.expCode, w.Code)
		assert.Equal(t, tc.expComment, cm)
	}
}

func TestHandler_DeleteyByID(t *testing.T) {
	testcases := []struct {
		name    string
		mock    func(*mockcomment.MockService, model.Comment)
		comment model.Comment
		expCode int
	}{
		{
			name: "comment is deleted",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().DeleteByID(cm.ID).Return(nil)
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusNoContent,
		},
		{
			name: "comment is not found",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().DeleteByID(cm.ID).Return(ErrNotFound)
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mock: func(s *mockcomment.MockService, cm model.Comment) {
				s.EXPECT().DeleteByID(cm.ID).Return(errors.New("internal error"))
			},
			comment: model.Comment{ID: 1, Body: "Comment 1"},
			expCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		c := gomock.NewController(t)
		defer c.Finish()
		ps := mockcomment.NewMockService(c)
		as := mockauth.NewMockService(c)
		tc.mock(ps, tc.comment)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/comment/1", nil)

		ctx := echo.New().NewContext(r, w)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")

		NewHandler(ps, as).DeleteByID(ctx)

		assert.Equal(t, tc.expCode, w.Code)
	}
}
