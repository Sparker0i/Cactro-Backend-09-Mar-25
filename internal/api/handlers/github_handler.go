package handlers

import (
	"net/http"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GitHubHandler handles GitHub-related API requests
type GitHubHandler struct {
	service *services.GitHubService
}

// NewGitHubHandler creates a new GitHubHandler
func NewGitHubHandler(service *services.GitHubService) *GitHubHandler {
	return &GitHubHandler{
		service: service,
	}
}

// GetUserProfile handles GET /github
func (h *GitHubHandler) GetUserProfile(c *gin.Context) {
	profile, err := h.service.GetUserProfile()
	if err != nil {
		logrus.WithError(err).Error("Failed to get user profile")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve GitHub profile",
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetRepository handles GET /github/:repo
func (h *GitHubHandler) GetRepository(c *gin.Context) {
	repoName := c.Param("repo")
	if repoName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Repository name is required",
		})
		return
	}

	repo, err := h.service.GetRepository(repoName)
	if err != nil {
		logrus.WithError(err).WithField("repo", repoName).Error("Failed to get repository")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve repository information",
		})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// CreateIssue handles POST /github/:repo/issues
func (h *GitHubHandler) CreateIssue(c *gin.Context) {
	repoName := c.Param("repo")
	if repoName == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Repository name is required",
		})
		return
	}

	var issueRequest models.IssueRequest
	if err := c.ShouldBindJSON(&issueRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request: title and body are required",
		})
		return
	}

	// Create the issue
	issue, err := h.service.CreateIssue(repoName, &issueRequest)
	if err != nil {
		logrus.WithError(err).WithField("repo", repoName).Error("Failed to create issue")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create issue",
		})
		return
	}

	c.JSON(http.StatusCreated, issue)
}
