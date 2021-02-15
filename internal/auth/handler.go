package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/imarrche/nix-ed/internal/config"
)

// Handler is http handler for authorization/authentication.
type Handler struct {
	s Service
}

// NewHandler creates and returns a new Handler instance.
func NewHandler(s Service) *Handler {
	return &Handler{s}
}

// respond responds to request with XML or JSON.
func respond(c echo.Context, code int, data interface{}) error {
	if c.Request().Header.Get("Accept-Encoding") == "text/xml" {
		return c.XML(code, data)
	}

	return c.JSON(code, data)
}

// GoogleSignIn is google sign in handler.
func (h *Handler) GoogleSignIn(c echo.Context) error {
	url := h.s.AuthCodeURL(config.Get().AuthCodeURL)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles redirect after signing in and returns an access token.
func (h *Handler) GoogleCallback(c echo.Context) error {
	token, err := h.s.GetAccessToken(c.FormValue("code"))
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}
