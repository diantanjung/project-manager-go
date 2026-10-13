package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) register(c *gin.Context) {
	var input service.CreateUserInput
	if !bindJSON(c, &input) {
		return
	}
	user, err := h.service.Register(c.Request.Context(), input)
	respond(c, http.StatusCreated, user, err)
}

func (h *Handler) login(c *gin.Context) {
	var input service.LoginInput
	if !bindJSON(c, &input) {
		return
	}
	result, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	setRefreshCookie(c, h.cfg, result.RefreshToken)
	c.Header("Cache-Control", "no-store")
	respond(c, http.StatusOK, gin.H{"user": result.User, "accessToken": result.AccessToken, "refreshToken": result.RefreshToken}, nil)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, ok := requestRefreshToken(c)
	if !ok {
		return
	}
	if refreshToken == "" {
		respondError(c, domain.NewError(domain.ErrInvalidRefresh, "Refresh token not found"))
		return
	}
	result, err := h.service.RefreshAccessToken(c.Request.Context(), refreshToken)
	if err != nil {
		clearRefreshCookie(c, h.cfg)
		respondError(c, err)
		return
	}
	setRefreshCookie(c, h.cfg, result.RefreshToken)
	c.Header("Cache-Control", "no-store")
	respond(c, http.StatusOK, gin.H{"accessToken": result.AccessToken, "refreshToken": result.RefreshToken}, nil)
}

func (h *Handler) me(c *gin.Context) {
	user, err := h.service.CurrentUser(c.Request.Context(), currentUser(c))
	respond(c, http.StatusOK, user, err)
}

func (h *Handler) logout(c *gin.Context) {
	refreshToken, ok := requestRefreshToken(c)
	if !ok {
		return
	}
	err := h.service.Logout(c.Request.Context(), currentUser(c), refreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	clearRefreshCookie(c, h.cfg)
	c.Status(http.StatusNoContent)
}

func requestRefreshToken(c *gin.Context) (string, bool) {
	if token, err := c.Cookie("refreshToken"); err == nil && token != "" {
		return token, true
	}
	if c.Request.ContentLength == 0 {
		return "", true
	}
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if !bindJSON(c, &body) {
		return "", false
	}
	return body.RefreshToken, true
}
