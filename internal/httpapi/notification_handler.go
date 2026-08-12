package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
)

func (h *Handler) listNotifications(c *gin.Context) {
	items, err := h.service.ListNotifications(c.Request.Context(), currentUser(c).ID)
	respondWithMeta(c, http.StatusOK, items, gin.H{"unreadCount": unreadNotificationCount(items)}, err)
}

func (h *Handler) markNotificationRead(c *gin.Context) {
	item, err := h.service.MarkNotificationRead(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, item, err)
}

func (h *Handler) markAllNotificationsRead(c *gin.Context) {
	items, err := h.service.MarkAllNotificationsRead(c.Request.Context(), currentUser(c))
	respond(c, http.StatusOK, items, err)
}

func unreadNotificationCount(items []domain.Notification) int {
	count := 0
	for _, item := range items {
		if !item.IsRead {
			count++
		}
	}
	return count
}
