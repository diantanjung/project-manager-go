package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/service"
)

func (h *Handler) listTeams(c *gin.Context) {
	filter, ok := pageFilter(c)
	if !ok {
		return
	}
	result, err := h.service.ListTeams(c.Request.Context(), filter)
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) createTeam(c *gin.Context) {
	var request createTeamRequest
	if !bindJSON(c, &request) {
		return
	}
	team, err := h.service.CreateTeam(c.Request.Context(), currentUser(c), request.input())
	respond(c, http.StatusCreated, team, err)
}

func (h *Handler) getTeam(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	team, err := h.service.GetTeamByID(c.Request.Context(), id)
	respond(c, http.StatusOK, team, err)
}

func (h *Handler) updateTeam(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var input service.TeamInput
	if !bindJSON(c, &input) {
		return
	}
	team, err := h.service.UpdateTeam(c.Request.Context(), currentUser(c), id, input)
	respond(c, http.StatusOK, team, err)
}

func (h *Handler) deleteTeam(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	_, err := h.service.DeleteTeam(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, gin.H{"message": "Team deleted successfully"}, err)
}

func (h *Handler) listTeamMembers(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	members, err := h.service.ListTeamMembers(c.Request.Context(), id)
	respond(c, http.StatusOK, members, err)
}

func (h *Handler) addTeamMember(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var input service.TeamMemberInput
	if !bindJSON(c, &input) {
		return
	}
	member, err := h.service.AddTeamMember(c.Request.Context(), currentUser(c), id, input)
	respond(c, http.StatusCreated, member, err)
}

func (h *Handler) removeTeamMember(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	userID, ok := pathInt(c, "userId")
	if !ok {
		return
	}
	_, err := h.service.RemoveTeamMember(
		c.Request.Context(),
		currentUser(c),
		id,
		userID,
	)
	respond(c, http.StatusOK, gin.H{"message": "Member removed from team"}, err)
}
