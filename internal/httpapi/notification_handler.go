package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) listNotifications(c *gin.Context) {
	items, err := h.service.ListNotifications(c.Request.Context(), currentUser(c).ID)
	respond(c, http.StatusOK, gin.H{"success": true, "count": len(items), "data": items}, err)
}

func (h *Handler) markNotificationRead(c *gin.Context) {
	item, err := h.service.MarkNotificationRead(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"success": true, "message": "Notification marked as read", "data": item}, err)
}

func (h *Handler) markAllNotificationsRead(c *gin.Context) {
	items, err := h.service.MarkAllNotificationsRead(c.Request.Context(), currentUser(c))
	respond(c, http.StatusOK, gin.H{"success": true, "message": "All notifications marked as read", "data": items}, err)
}
