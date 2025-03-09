package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/config"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
	"github.com/sirupsen/logrus"
)

// GitHubService provides methods for interacting with the GitHub API
type GitHubService struct {
	config *config.Config
	client *http.Client
}

// NewGitHubService creates a new GitHubService
func NewGitHubService(config *config.Config) GitHubServiceInterface {
	return &GitHubService{
		config: config,
		client: &http.Client{},
	}
}

// GetUserProfile retrieves the user's GitHub profile
func (s *GitHubService) GetUserProfile() (*models.GithubProfile, error) {
	// Get user data
	user, err := s.getUser()
	if err != nil {
		return nil, err
	}

	// Get user repositories
	repos, err := s.getUserRepositories()
	if err != nil {
		return nil, err
	}

	return &models.GithubProfile{
		User:         *user,
		Repositories: repos,
	}, nil
}

// GetRepository retrieves details about a specific repository
func (s *GitHubService) GetRepository(repoName string) (*models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", s.config.GitHub.Username, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", s.config.GitHub.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("GitHub API error")
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var repo models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repo, nil
}

// CreateIssue creates a new issue in a repository
func (s *GitHubService) CreateIssue(repoName string, issue *models.IssueRequest) (*models.IssueResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", s.config.GitHub.Username, repoName)

	// Create request body
	body, err := json.Marshal(issue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", s.config.GitHub.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("GitHub API error")
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var issueResponse models.IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &issueResponse, nil
}

// getUser retrieves the user's GitHub profile
func (s *GitHubService) getUser() (*models.UserResponse, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", s.config.GitHub.Username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", s.config.GitHub.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("GitHub API error")
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var user models.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// getUserRepositories retrieves the user's repositories
func (s *GitHubService) getUserRepositories() ([]models.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=owner&sort=updated&per_page=100", s.config.GitHub.Username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("token %s", s.config.GitHub.Token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"response":    string(body),
		}).Error("GitHub API error")
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	var repos []models.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return repos, nil
}
