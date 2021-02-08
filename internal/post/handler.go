package post

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/imarrche/nix-ed/internal/model"
)

// Handler is http handler for post resource.
type Handler struct {
	s Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(s Service) *Handler {
	return &Handler{s}
}

// GetAll returns post list.
func (h *Handler) GetAll(c echo.Context) error {
	ps, err := h.s.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, ps)

}

// Create creates a post.
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

	return c.JSON(http.StatusCreated, p)
}

// Update updates a post.
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

	return c.JSON(http.StatusCreated, p)
}

// DeleteByID deletes a post.
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
