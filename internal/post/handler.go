package post

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
	uIDkey key = iota
)

type errResponse struct {
	Title  string `json:"title" xml:"title"`
	Body   string `json:"body" xml:"body"`
	UserID string `json:"userId" xml:"userId"`
}

// Handler is http handler for post resource.
type Handler struct {
	ps Service
	as auth.Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(ps Service, as auth.Service) *Handler {
	return &Handler{ps: ps, as: as}
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
		r = r.WithContext(context.WithValue(r.Context(), uIDkey, udata.ID))
		c.SetRequest(r)
		return next(c)
	}
}

// PostAuthor is middleware that ensures that post's author made a request.
func (h *Handler) PostAuthor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, ok := c.Request().Context().Value(uIDkey).(string)
		if !ok {
			return c.NoContent(http.StatusInternalServerError)
		}
		pID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return respond(c, http.StatusBadRequest, err)
		}

		p, err := h.ps.GetByID(pID)
		if err == ErrNotFound {
			return c.NoContent(http.StatusNotFound)
		} else if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		if p.UserID != uID {
			return c.NoContent(http.StatusForbidden)
		}

		return next(c)
	}
}

// GetAll returns post list.
// @Summary Show all posts
// @Descriptions show all posts
// @Tags posts
// @ID post-list
// @Accept json
// @Produce json,xml
// @Success 200 {array} model.Post
// @Failure 500 ""
// @Router /posts [get]
func (h *Handler) GetAll(c echo.Context) error {
	ps, err := h.ps.GetAll()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return respond(c, http.StatusOK, ps)
}

// Create creates a post.
// @Summary Create a post
// @Descriptions create a post
// @Tags posts
// @ID post-create
// @Accept json
// @Produce json,xml
// @Param input body model.Post true "post data"
// @Success 201 {object} model.Post
// @Failure 400 {object} errResponse
// @Router /posts [post]
func (h *Handler) Create(c echo.Context) error {
	id, ok := c.Request().Context().Value(uIDkey).(string)
	if !ok {
		return c.NoContent(http.StatusInternalServerError)
	}

	p := model.Post{}
	if err := c.Bind(&p); err != nil {
		return respond(c, http.StatusBadRequest, err)
	}
	p.UserID = id

	p, err := h.ps.Create(p)
	if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusCreated, p)
}

// GetByID returns post detail.
// @Summary Post detail
// @Descriptions post detail
// @Tags posts
// @ID post-detail
// @Accept json
// @Produce json,xml
// @Param id path int true "post id"
// @Success 200 {object} model.Post
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /posts/{id} [get]
func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	p, err := h.ps.GetByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusOK, p)
}

// Update updates a post.
// @Summary Post update
// @Descriptions post update
// @Tags posts
// @ID post-update
// @Accept json
// @Produce json,xml
// @Param id path int true "post id"
// @Param input body model.Post true "post data"
// @Success 200 {object} model.Post
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /posts/{id} [patch]
func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	p := model.Post{}
	if err := c.Bind(&p); err != nil {
		return respond(c, http.StatusBadRequest, err)
	}
	p.ID = id

	p, err = h.ps.Update(p)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return respond(c, http.StatusOK, p)
}

// DeleteByID deletes a post.
// @Summary Post delete
// @Descriptions post delete
// @Tags posts
// @ID posts-delete
// @Accept json
// @Produce json,xml
// @Param id path int true "post id"
// @Success 204 ""
// @Failure 400 {object} errResponse
// @Failure 404 ""
// @Router /posts/{id} [delete]
func (h *Handler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return respond(c, http.StatusBadRequest, err)
	}

	err = h.ps.DeleteByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return respond(c, http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
