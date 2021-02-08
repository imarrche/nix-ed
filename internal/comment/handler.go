package comment

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/imarrche/nix-ed/internal/model"
)

type errResponse struct {
	Name   string `json:"name" xml:"name"`
	Email  string `json:"email" xml:"email"`
	Body   string `json:"body" xml:"body"`
	UserID string `json:"userId" xml:"userId"`
}

// Handler is http handler for comment resource.
type Handler struct {
	s Service
}

// NewHandler creates and returns a new Handler instacne.
func NewHandler(s Service) *Handler {
	return &Handler{s}
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
	cs, err := h.s.GetAll()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, cs)

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
		return c.JSON(http.StatusBadRequest, err)
	}

	cm, err := h.s.GetByID(id)
	if err == ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, cm)
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

	return c.JSON(http.StatusOK, cm)
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
