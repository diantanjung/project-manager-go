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
	respond(c, http.StatusOK, gin.H{"user": result.User, "accessToken": result.AccessToken}, nil)
}

func (h *Handler) refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
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
	respond(c, http.StatusOK, gin.H{"accessToken": result.AccessToken}, nil)
}

func (h *Handler) me(c *gin.Context) {
	user, err := h.service.CurrentUser(c.Request.Context(), currentUser(c))
	respond(c, http.StatusOK, user, err)
}

func (h *Handler) logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refreshToken")
	err := h.service.Logout(c.Request.Context(), currentUser(c), refreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	clearRefreshCookie(c, h.cfg)
	c.Status(http.StatusNoContent)
}
