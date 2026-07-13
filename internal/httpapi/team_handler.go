package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/service"
)

func (h *Handler) listTeams(c *gin.Context) {
	result, err := h.service.ListTeams(c.Request.Context(), pageFilter(c))
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) createTeam(c *gin.Context) {
	var input service.TeamInput
	if !bindJSON(c, &input) {
		return
	}
	team, err := h.service.CreateTeam(c.Request.Context(), currentUser(c), input)
	respond(c, http.StatusCreated, team, err)
}

func (h *Handler) getTeam(c *gin.Context) {
	team, err := h.service.GetTeamByID(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, team, err)
}

func (h *Handler) updateTeam(c *gin.Context) {
	var input service.TeamInput
	if !bindJSON(c, &input) {
		return
	}
	team, err := h.service.UpdateTeam(c.Request.Context(), currentUser(c), pathInt(c, "id"), input)
	respond(c, http.StatusOK, team, err)
}

func (h *Handler) deleteTeam(c *gin.Context) {
	_, err := h.service.DeleteTeam(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Team deleted successfully"}, err)
}

func (h *Handler) listTeamMembers(c *gin.Context) {
	members, err := h.service.ListTeamMembers(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, members, err)
}

func (h *Handler) addTeamMember(c *gin.Context) {
	var input service.TeamMemberInput
	if !bindJSON(c, &input) {
		return
	}
	member, err := h.service.AddTeamMember(c.Request.Context(), currentUser(c), pathInt(c, "id"), input)
	respond(c, http.StatusCreated, member, err)
}

func (h *Handler) removeTeamMember(c *gin.Context) {
	_, err := h.service.RemoveTeamMember(
		c.Request.Context(),
		currentUser(c),
		pathInt(c, "id"),
		pathInt(c, "userId"),
	)
	respond(c, http.StatusOK, gin.H{"message": "Member removed from team"}, err)
}
