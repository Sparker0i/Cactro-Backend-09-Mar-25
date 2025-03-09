package services

import (
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
)

// GitHubServiceInterface defines the interface for GitHub service operations
type GitHubServiceInterface interface {
	// GetUserProfile retrieves the user's GitHub profile
	GetUserProfile() (*models.GithubProfile, error)

	// GetRepository retrieves details about a specific repository
	GetRepository(repoName string) (*models.Repository, error)

	// CreateIssue creates a new issue in a repository
	CreateIssue(repoName string, issue *models.IssueRequest) (*models.IssueResponse, error)
}
