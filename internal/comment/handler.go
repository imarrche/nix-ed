package comment

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/imarrche/nix-ed/internal/model"
)

// Handler is http handler for all comment endpoints.
type Handler struct {
	r *chi.Mux
	s Service
}

// NewHandler creates and returns a new Handler instance.
func NewHandler(s Service) *Handler {
	return &Handler{chi.NewRouter(), s}
}

// Init initializes all handler's routes.
func (h *Handler) Init(r chi.Router) {
	r.Get("/", h.getAll)
	r.Post("/", h.create)
	r.Get("/{id:[0-9]+}", h.detail)
	r.Patch("/{id:[0-9]+}", h.update)
	r.Delete("/{id:[0-9]+}", h.delete)
}

// ServeHTTP makes router to handle requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

// respond is a helper function for encapsulating responding logic.
func (h *Handler) respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)

	type encoder interface {
		Encode(interface{}) error
	}

	var enc encoder
	if r.Header.Get("Accept-Encoding") == "application/xml" {
		w.Header().Set("Content-Type", "application/xml")
		enc = xml.NewEncoder(w)
	} else {
		w.Header().Set("Content-Type", "application/json")
		enc = json.NewEncoder(w)
	}

	if e, ok := data.(validation.Errors); ok {
		enc.Encode(e)
	} else if e, ok := data.(error); ok {
		enc.Encode(map[string]interface{}{"error": e.Error()})
	} else if data != nil {
		enc.Encode(data)
	}
}

// getAll handles retrieving all comments.
func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	cs, err := h.s.GetAll()
	if err != nil {
		h.respond(w, r, http.StatusInternalServerError, nil)
		return
	}

	h.respond(w, r, http.StatusOK, cs)
}

// create handles comment creating.
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	c := model.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}

	c, err := h.s.Create(c)
	if err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}

	h.respond(w, r, http.StatusCreated, c)
}

// detail handles retrieving a single comment.
func (h *Handler) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respond(w, r, http.StatusInternalServerError, nil)
		return
	}

	c, err := h.s.GetByID(id)
	if err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}

	h.respond(w, r, http.StatusOK, c)
}

// update handles updating a comment.
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respond(w, r, http.StatusInternalServerError, nil)
		return
	}

	c := model.Comment{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}
	c.ID = id

	c, err = h.s.Update(c)
	if err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}

	h.respond(w, r, http.StatusOK, c)
}

// delete handles deleting a comment.
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		h.respond(w, r, http.StatusInternalServerError, nil)
		return
	}

	if err = h.s.DeleteByID(id); err != nil {
		h.respond(w, r, http.StatusBadRequest, err)
		return
	}

	h.respond(w, r, http.StatusNoContent, nil)
}
