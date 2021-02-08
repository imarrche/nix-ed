package comment

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/imarrche/nix-ed/internal/model"
)

// Handler is http handler for comment resource.
type Handler struct {
	s Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(s Service) *Handler {
	return &Handler{s}
}

// GetAll returns comment list.
func (h *Handler) GetAll(c echo.Context) error {
	cs, err := h.s.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, cs)

}

// Create creates a comment.
func (h *Handler) Create(c echo.Context) error {
	cm := model.Comment{}
	if err := c.Bind(&cm); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cm, err := h.s.Create(cm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, cm)
}

// GetByID returns comment detail.
func (h *Handler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cm, err := h.s.GetByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, cm)
}

// Update updates a comment.
func (h *Handler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cm := model.Comment{}
	if err := c.Bind(&cm); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	cm.ID = id

	cm, err = h.s.Update(cm)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, cm)
}

// DeleteByID deletes a comment.
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
