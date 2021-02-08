package post

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/imarrche/nix-ed/internal/model"
)

type errResponse struct {
	Title  string `json:"title" xml:"title"`
	Body   string `json:"body" xml:"body"`
	UserID string `json:"userId" xml:"userId"`
}

// Handler is http handler for post resource.
type Handler struct {
	s Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(s Service) *Handler {
	return &Handler{s}
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
	ps, err := h.s.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, ps)

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
	p := model.Post{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	p, err := h.s.Create(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, p)
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
		return c.JSON(http.StatusBadRequest, err)
	}

	p, err := h.s.GetByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, p)
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
		return c.JSON(http.StatusBadRequest, err)
	}

	p := model.Post{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	p.ID = id

	p, err = h.s.Update(p)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, p)
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
		return c.JSON(http.StatusBadRequest, err)
	}

	err = h.s.DeleteByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.NoContent(http.StatusNoContent)
}
