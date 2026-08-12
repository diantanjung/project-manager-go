package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) listProjects(c *gin.Context) {
	result, err := h.service.ListProjects(c.Request.Context(), currentUser(c), service.ProjectFilter{
		PageFilter: pageFilter(c),
		Search:     c.Query("search"),
		TeamID:     queryIntPtr(c, "teamId"),
		SortBy:     c.DefaultQuery("sortBy", "createdAt"),
		Order:      c.DefaultQuery("order", "desc"),
	})
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) createProject(c *gin.Context) {
	var request createProjectRequest
	if !bindJSON(c, &request) {
		return
	}
	project, err := h.service.CreateProject(c.Request.Context(), currentUser(c), request.input())
	respond(c, http.StatusCreated, project, err)
}

func (h *Handler) getProject(c *gin.Context) {
	project, err := h.service.GetProjectByID(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, project, err)
}

func (h *Handler) updateProject(c *gin.Context) {
	var input service.ProjectInput
	if !bindJSON(c, &input) {
		return
	}
	project, err := h.service.UpdateProject(c.Request.Context(), currentUser(c), pathInt(c, "id"), input)
	respond(c, http.StatusOK, project, err)
}

func (h *Handler) deleteProject(c *gin.Context) {
	_, err := h.service.DeleteProject(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Project deleted successfully"}, err)
}

func (h *Handler) listProjectTasks(c *gin.Context) {
	result, err := h.service.ListProjectTasks(c.Request.Context(), currentUser(c), pathInt(c, "id"), pageFilter(c))
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) listProjectTeams(c *gin.Context) {
	teams, err := h.service.ListProjectTeams(c.Request.Context(), pathInt(c, "id"))
	respond(c, http.StatusOK, teams, err)
}

func (h *Handler) listProjectTeamsByProjectID(c *gin.Context) {
	teams, err := h.service.ListProjectTeams(c.Request.Context(), pathInt(c, "projectId"))
	respond(c, http.StatusOK, teams, err)
}

func (h *Handler) assignTeamToProject(c *gin.Context) {
	var input service.ProjectTeamInput
	if !bindJSON(c, &input) {
		return
	}
	assignment, exists, err := h.service.AssignTeamToProject(c.Request.Context(), currentUser(c), input)
	if exists && err == nil {
		err = domain.NewError(domain.ErrDuplicate, "Team is already assigned to this project")
	}
	respond(c, http.StatusCreated, assignment, err)
}

func (h *Handler) removeTeamFromProject(c *gin.Context) {
	_, err := h.service.RemoveTeamFromProject(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Team removed from project successfully"}, err)
}
