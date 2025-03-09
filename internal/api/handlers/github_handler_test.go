package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGitHubService is a mock implementation of the GitHub service interface
type MockGitHubService struct {
	mock.Mock
}

// Ensure MockGitHubService implements GitHubServiceInterface
var _ services.GitHubServiceInterface = (*MockGitHubService)(nil)

// GetUserProfile mocks the GetUserProfile method
func (m *MockGitHubService) GetUserProfile() (*models.GithubProfile, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GithubProfile), args.Error(1)
}

// GetRepository mocks the GetRepository method
func (m *MockGitHubService) GetRepository(repoName string) (*models.Repository, error) {
	args := m.Called(repoName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Repository), args.Error(1)
}

// CreateIssue mocks the CreateIssue method
func (m *MockGitHubService) CreateIssue(repoName string, issue *models.IssueRequest) (*models.IssueResponse, error) {
	args := m.Called(repoName, issue)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IssueResponse), args.Error(1)
}

// SetupTestRouter creates a router for testing
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// TestGetUserProfile tests the GetUserProfile handler
func TestGetUserProfile(t *testing.T) {
	// Create test data
	mockProfile := &models.GithubProfile{
		User: models.UserResponse{
			Login:     "test-user",
			Name:      "Test User",
			Followers: 10,
			Following: 20,
		},
		Repositories: []models.Repository{
			{
				Name:     "repo1",
				FullName: "test-user/repo1",
			},
			{
				Name:     "repo2",
				FullName: "test-user/repo2",
			},
		},
	}

	// Test cases
	tests := []struct {
		name               string
		setupMock          func(mockService *MockGitHubService)
		expectedStatusCode int
	}{
		{
			name: "Success",
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("GetUserProfile").Return(mockProfile, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Service Error",
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("GetUserProfile").Return(nil, errors.New("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			router := SetupTestRouter()
			mockService := new(MockGitHubService)
			tc.setupMock(mockService)

			// Create a handler with our mock service - using the constructor now
			handler := NewGitHubHandler(mockService)

			router.GET("/github", handler.GetUserProfile)

			// Create a test request
			req, _ := http.NewRequest("GET", "/github", nil)
			resp := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(resp, req)

			// Check the response
			assert.Equal(t, tc.expectedStatusCode, resp.Code)

			if tc.expectedStatusCode == http.StatusOK {
				var response models.GithubProfile
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, mockProfile.User.Login, response.User.Login)
				assert.Len(t, response.Repositories, 2)
			}

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

// TestGetRepository tests the GetRepository handler
func TestGetRepository(t *testing.T) {
	// Create test data
	mockRepo := &models.Repository{
		Name:     "test-repo",
		FullName: "test-user/test-repo",
	}

	// Test cases
	tests := []struct {
		name               string
		repoName           string
		setupMock          func(mockService *MockGitHubService)
		expectedStatusCode int
	}{
		{
			name:     "Success",
			repoName: "test-repo",
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("GetRepository", "test-repo").Return(mockRepo, nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:     "Repository Not Found",
			repoName: "non-existent-repo",
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("GetRepository", "non-existent-repo").Return(nil, errors.New("not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:     "Empty Repository Name",
			repoName: "",
			setupMock: func(mockService *MockGitHubService) {
				// No mock setup needed as the request will be rejected by the handler
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			router := SetupTestRouter()
			mockService := new(MockGitHubService)
			tc.setupMock(mockService)

			// Create a handler with our mock service
			handler := NewGitHubHandler(mockService)

			router.GET("/github/:repo", handler.GetRepository)

			// Create a test request
			req, _ := http.NewRequest("GET", "/github/"+tc.repoName, nil)
			resp := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(resp, req)

			// Check the response
			assert.Equal(t, tc.expectedStatusCode, resp.Code)

			if tc.expectedStatusCode == http.StatusOK {
				var response models.Repository
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, mockRepo.Name, response.Name)
				assert.Equal(t, mockRepo.FullName, response.FullName)
			}

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

// TestCreateIssue tests the CreateIssue handler
func TestCreateIssue(t *testing.T) {
	// Create test data
	mockIssueRequest := &models.IssueRequest{
		Title: "Test Issue",
		Body:  "This is a test issue",
	}

	mockIssueResponse := &models.IssueResponse{
		Number:  1,
		Title:   "Test Issue",
		HTMLURL: "https://github.com/test-user/test-repo/issues/1",
	}

	// Test cases
	tests := []struct {
		name               string
		repoName           string
		requestBody        interface{}
		setupMock          func(mockService *MockGitHubService)
		expectedStatusCode int
	}{
		{
			name:        "Success",
			repoName:    "test-repo",
			requestBody: mockIssueRequest,
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("CreateIssue", "test-repo", mock.MatchedBy(func(req *models.IssueRequest) bool {
					return req.Title == mockIssueRequest.Title && req.Body == mockIssueRequest.Body
				})).Return(mockIssueResponse, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:        "Repository Not Found",
			repoName:    "non-existent-repo",
			requestBody: mockIssueRequest,
			setupMock: func(mockService *MockGitHubService) {
				mockService.On("CreateIssue", "non-existent-repo", mock.MatchedBy(func(req *models.IssueRequest) bool {
					return req.Title == mockIssueRequest.Title && req.Body == mockIssueRequest.Body
				})).Return(nil, errors.New("not found"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:     "Invalid Request",
			repoName: "test-repo",
			requestBody: map[string]string{
				"title": "", // Empty title
			},
			setupMock: func(mockService *MockGitHubService) {
				// No mock setup needed as the request will be rejected by the handler
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Empty Repository Name",
			repoName:    "",
			requestBody: mockIssueRequest,
			setupMock: func(mockService *MockGitHubService) {
				// No mock setup needed as the request will be rejected by the handler
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			router := SetupTestRouter()
			mockService := new(MockGitHubService)
			tc.setupMock(mockService)

			// Create a handler with our mock service
			handler := NewGitHubHandler(mockService)

			router.POST("/github/:repo/issues", handler.CreateIssue)

			// Create request body
			requestBody, _ := json.Marshal(tc.requestBody)

			// Create a test request
			req, _ := http.NewRequest("POST", "/github/"+tc.repoName+"/issues", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(resp, req)

			// Check the response
			assert.Equal(t, tc.expectedStatusCode, resp.Code)

			if tc.expectedStatusCode == http.StatusCreated {
				var response models.IssueResponse
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, mockIssueResponse.Number, response.Number)
				assert.Equal(t, mockIssueResponse.Title, response.Title)
				assert.Equal(t, mockIssueResponse.HTMLURL, response.HTMLURL)
			}

			// Verify that all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
