package httpapi

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) listComments(c *gin.Context) {
	comments, err := h.service.ListComments(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, comments, err)
}

func (h *Handler) createComment(c *gin.Context) {
	var input service.CommentInput
	if !bindJSON(c, &input) {
		return
	}
	input.TaskID = pathInt(c, "id")
	comment, err := h.service.CreateComment(c.Request.Context(), currentUser(c), input)
	respond(c, http.StatusCreated, comment, err)
}

func (h *Handler) updateComment(c *gin.Context) {
	var input struct {
		Content string `json:"content"`
	}
	if !bindJSON(c, &input) {
		return
	}
	comment, err := h.service.UpdateComment(c.Request.Context(), currentUser(c), pathInt(c, "id"), input.Content)
	respond(c, http.StatusOK, comment, err)
}

func (h *Handler) deleteComment(c *gin.Context) {
	_, err := h.service.DeleteComment(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Comment deleted successfully"}, err)
}

func (h *Handler) listAttachments(c *gin.Context) {
	attachments, err := h.service.ListAttachments(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, attachments, err)
}

func (h *Handler) createAttachment(c *gin.Context) {
	var input service.AttachmentInput
	if !bindJSON(c, &input) {
		return
	}
	attachment, err := h.service.CreateAttachment(c.Request.Context(), currentUser(c), pathInt(c, "id"), input)
	respond(c, http.StatusCreated, attachment, err)
}

func (h *Handler) getAttachment(c *gin.Context) {
	attachment, err := h.service.GetAttachmentByID(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, attachment, err)
}

func (h *Handler) deleteAttachment(c *gin.Context) {
	_, err := h.service.DeleteAttachment(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Attachment deleted successfully"}, err)
}

func (h *Handler) uploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		respondError(c, domain.NewError(domain.ErrValidation, "No file uploaded"))
		return
	}
	if file.Size > 5*1024*1024 {
		respondError(c, domain.NewError(domain.ErrValidation, "File too large"))
		return
	}
	if err := os.MkdirAll(h.cfg.UploadDir, 0o755); err != nil {
		respondError(c, err)
		return
	}
	name := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	target := filepath.Join(h.cfg.UploadDir, name)
	if err := c.SaveUploadedFile(file, target); err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, gin.H{"url": "/uploads/" + name}, nil)
}
