package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) dashboard(c *gin.Context) {
	summary, err := h.service.DashboardSummary(c.Request.Context(), currentUser(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}
