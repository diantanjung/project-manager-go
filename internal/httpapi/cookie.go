package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/config"
)

func setRefreshCookie(c *gin.Context, cfg config.Config, token string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		"refreshToken",
		token,
		int(cfg.JWTRefreshExpiresIn.Seconds()),
		"/",
		"",
		cfg.NodeEnv == "production",
		true,
	)
}

func clearRefreshCookie(c *gin.Context, cfg config.Config) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refreshToken", "", -1, "/", "", cfg.NodeEnv == "production", true)
}
