package comment

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/imarrche/nix-ed/internal/auth"
	"github.com/imarrche/nix-ed/internal/model"
)

type key int

const (
	uEmailKey key = iota
)

type errResponse struct {
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
	UserID string `json:"userId" xml:"userId"`
}

// Handler is http handler for comment resource.
type Handler struct {
	cs Service
	as auth.Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(cs Service, as auth.Service) *Handler {
	return &Handler{cs: cs, as: as}
}

// respond responds to request with XML or JSON.
func respond(c echo.Context, code int, data interface{}) error {
	if c.Request().Header.Get("Accept-Encoding") == "text/xml" {
		return c.XML(code, data)
	}

	return c.JSON(code, data)
}

type authResponse struct {
	Error map[string]interface{} `json:"error"`
	ID    string                 `json:"id"`
	Email string                 `json:"email"`
}

// Auth is middleware for user authentication.
func (h *Handler) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data, err := h.as.GetUserInfo(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/auth/google/sign-in")
		}
		udata := &authResponse{}
		if err := json.Unmarshal(data, &udata); err != nil || udata.Error != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		r := c.Request()
		r = r.WithContext(context.WithValue(r.Context(), uEmailKey, udata.Email))
		c.SetRequest(r)
		return next(c)
	}
}

// CommentAuthor is middleware that ensures that comment's author made a request.
func (h *Handler) CommentAuthor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uEmail, ok := c.Request().Context().Value(uEmailKey).(string)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		cID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return respond(c, http.StatusBadRequest, err)
		}

		cm, err := h.cs.GetByID(cID)
		if err == ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		} else if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if cm.Email != uEmail {
			return c.NoContent(http.StatusForbidden)
		}

		return next(c)
	}
}

// GetAll returns comment list.
// @Summary Show all comments
// @Descriptions show all comments
// @Tags comments
// @ID comment-list
// @Accept json
// @Produce json,xml
// @Success 200 {array} model.Comment
// @Failure 500 ""
// @Router /comments [get]
func (h *Handler) GetAll(c echo.Context) error {
	cs, err := h.cs.GetAll()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return respond(c, http.StatusOK, cs)

}

// Create creates a comment.
// @Summary Create a comment
// @Descriptions create a comment
// @Tags comments
// @ID comment-create
// @Accept json
// @Produce json,xml
// @Param input body model.Comment true "comment data"
// @Success 201 {object} model.Comment
// @Failure 400 {object} errResponse
// @Router /comments [post]
func (h *Handler) Create(c echo.Context) error {
	email, ok := c.Request().Context().Value(uEmailKey).(string)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	cm := model.Comment{}
	if err := c.Bind(&cm); err != nil {
		return respond(c, http.StatusBadRequest, err)
	}
	cm.Email = email

	cm, err := h.cs.Create(cm)
	if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusCreated, cm)
}

// GetByID returns comment detail.
// @Summary Comment detail
// @Descriptions comment detail
// @Tags comments
// @ID comment-detail
// @Accept json
// @Produce json,xml
// @Param id path int true "comment id"
// @Success 200 {object} model.Comment
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /comments/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	cm, err := h.cs.GetByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusOK, cm)
}

// Update updates a comment.
// @Summary Comment update
// @Descriptions comment update
// @Tags comments
// @ID comment-update
// @Accept json
// @Produce xml
// @Param id path int true "comment id"
// @Param input body model.Comment true "comment data"
// @Success 200 {object} model.Comment
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /comments/{id} [patch]
func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	cm := model.Comment{}
	if err := c.Bind(&cm); err != nil {
		return respond(c, http.StatusBadRequest, err)
	}
	cm.ID = id

	cm, err = h.cs.Update(cm)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusOK, cm)
}

// DeleteByID deletes a comment.
// @Summary Comment delete
// @Descriptions comment delete
// @Tags comments
// @ID comment-delete
// @Accept json
// @Produce json,xml
// @Param id path int true "comment id"
// @Success 204 ""
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /comments/{id} [delete]
func (h *Handler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	err = h.cs.DeleteByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
