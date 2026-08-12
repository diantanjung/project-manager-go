package httpapi

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"project-manager-go/internal/config"
	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

type Handler struct {
	service *service.Service
	cfg     config.Config
}

func NewRouter(cfg config.Config, svc *service.Service) *gin.Engine {
	if cfg.NodeEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := &Handler{service: svc, cfg: cfg}
	router.Static("/uploads", cfg.UploadDir)
	router.GET("/", h.root)
	router.GET("/health", h.health)

	h.mountAPI(router.Group("/api"))
	h.mountAPI(router.Group("/api/v1"))
	return router
}

func (h *Handler) mountAPI(api *gin.RouterGroup) {
	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", h.register)
	authRoutes.POST("/login", h.login)
	authRoutes.POST("/refresh", h.refresh)
	authRoutes.GET("/me", h.authRequired(), h.me)
	authRoutes.POST("/logout", h.authRequired(), h.logout)

	protected := api.Group("")
	protected.Use(h.authRequired())

	h.mountUserRoutes(protected)
	h.mountTeamRoutes(protected)
	h.mountProjectRoutes(protected)
	h.mountTaskRoutes(protected)
	h.mountCollaborationRoutes(protected)
	h.mountNotificationRoutes(protected)
}

func (h *Handler) root(c *gin.Context) {
	respond(c, http.StatusOK, gin.H{"message": "Welcome to Project Manager API"}, nil)
}

func (h *Handler) health(c *gin.Context) {
	respond(c, http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"env":       h.cfg.NodeEnv,
	}, nil)
}

func (h *Handler) mountUserRoutes(protected *gin.RouterGroup) {
	users := protected.Group("/users")
	users.GET("", h.listUsers)
	users.POST("", h.requireRole(domain.RoleAdmin), h.createUser)
	users.GET("/:id", h.getUser)
	users.PATCH("/:id", h.updateUser)
	users.DELETE("/:id", h.requireRole(domain.RoleAdmin), h.deleteUser)
	users.GET("/:id/tasks", h.listUserTasks)
}

func (h *Handler) mountTeamRoutes(protected *gin.RouterGroup) {
	teams := protected.Group("/teams")
	teams.GET("", h.listTeams)
	teams.POST("", h.requireRole(domain.RoleProductOwner), h.createTeam)
	teams.GET("/:id", h.getTeam)
	teams.PATCH("/:id", h.requireRole(domain.RoleProductOwner), h.updateTeam)
	teams.DELETE("/:id", h.requireRole(domain.RoleAdmin), h.deleteTeam)
	teams.GET("/:id/members", h.listTeamMembers)
	teams.POST("/:id/members", h.requireRole(domain.RoleProjectManager), h.addTeamMember)
	teams.DELETE("/:id/members/:userId", h.requireRole(domain.RoleProjectManager), h.removeTeamMember)
}

func (h *Handler) mountProjectRoutes(protected *gin.RouterGroup) {
	projects := protected.Group("/projects")
	projects.GET("", h.listProjects)
	projects.POST("", h.requireRole(domain.RoleProjectManager), h.createProject)
	projects.GET("/:id", h.getProject)
	projects.PATCH("/:id", h.requireRole(domain.RoleProjectManager), h.updateProject)
	projects.DELETE("/:id", h.requireRole(domain.RoleProductOwner), h.deleteProject)
	projects.GET("/:id/tasks", h.listProjectTasks)
	projects.GET("/:id/teams", h.listProjectTeams)

	projectTeams := protected.Group("/project-teams")
	projectTeams.GET("/projects/:projectId/teams", h.listProjectTeamsByProjectID)
	projectTeams.POST("", h.requireRole(domain.RoleProjectManager), h.assignTeamToProject)
	projectTeams.DELETE("/:id", h.requireRole(domain.RoleProductOwner), h.removeTeamFromProject)
}

func (h *Handler) mountTaskRoutes(protected *gin.RouterGroup) {
	tasks := protected.Group("/tasks")
	tasks.GET("", h.listTasks)
	tasks.POST("", h.requireRole(domain.RoleProjectManager), h.createTask)
	tasks.GET("/:id", h.getTask)
	tasks.PATCH("/:id", h.updateTask)
	tasks.PATCH("/:id/status", h.updateTaskStatus)
	tasks.DELETE("/:id", h.requireRole(domain.RoleProjectManager), h.deleteTask)
	tasks.GET("/:id/comments", h.listComments)
	tasks.POST("/:id/comments", h.createComment)
	tasks.GET("/:id/attachments", h.listAttachments)
	tasks.POST("/:id/attachments", h.createAttachment)

	taskAssignments := protected.Group("/task-assignments")
	taskAssignments.GET("/tasks/:taskId/assignments", h.listTaskAssignments)
	taskAssignments.POST("", h.requireRole(domain.RoleProjectManager), h.assignUserToTask)
	taskAssignments.DELETE("/:id", h.requireRole(domain.RoleProjectManager), h.removeTaskAssignment)
}

func (h *Handler) mountCollaborationRoutes(protected *gin.RouterGroup) {
	protected.PATCH("/comments/:id", h.updateComment)
	protected.DELETE("/comments/:id", h.deleteComment)
	protected.GET("/attachments/:id", h.getAttachment)
	protected.DELETE("/attachments/:id", h.deleteAttachment)
	protected.POST("/upload", h.uploadAvatar)
}

func (h *Handler) mountNotificationRoutes(protected *gin.RouterGroup) {
	notifications := protected.Group("/notifications")
	notifications.GET("", h.listNotifications)
	notifications.PATCH("/:id/read", h.markNotificationRead)
	notifications.PATCH("/read-all", h.markAllNotificationsRead)
}
