package services

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/config"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
	"github.com/stretchr/testify/assert"
)

// mockTransport is a custom RoundTripper that returns predefined responses
type mockTransport struct {
	mockResponse func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse(req)
}

// TestGetUserProfile tests the GetUserProfile function
func TestGetUserProfile(t *testing.T) {
	// Create config for testing
	cfg := &config.Config{
		GitHub: config.GitHubConfig{
			Token:    "test-token",
			Username: "test-user",
		},
	}

	// Create test cases
	tests := []struct {
		name           string
		userResponse   string
		reposResponse  string
		statusCode     int
		expectedError  bool
		expectedResult bool
	}{
		{
			name:           "Success",
			userResponse:   `{"login": "test-user", "name": "Test User", "followers": 10, "following": 20}`,
			reposResponse:  `[{"name": "repo1", "full_name": "test-user/repo1"}, {"name": "repo2", "full_name": "test-user/repo2"}]`,
			statusCode:     http.StatusOK,
			expectedError:  false,
			expectedResult: true,
		},
		{
			name:           "User API Error",
			userResponse:   `{"message": "Not Found"}`,
			reposResponse:  `[]`,
			statusCode:     http.StatusNotFound,
			expectedError:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a service instance
			// Need to type assert since we're now using an interface return type
			service := NewGitHubService(cfg).(*GitHubService)

			// Create a counter to track API calls
			callCount := 0

			// Replace the http client's transport with our mock
			service.client.Transport = &mockTransport{
				mockResponse: func(req *http.Request) (*http.Response, error) {
					callCount++

					var responseBody string
					if callCount == 1 {
						// First call should be to get user info
						responseBody = tc.userResponse
					} else {
						// Second call should be to get repositories
						responseBody = tc.reposResponse
					}

					// Create a mock response
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(strings.NewReader(responseBody)),
						Header:     make(http.Header),
					}, nil
				},
			}

			// Call the function
			result, err := service.GetUserProfile()

			// Check the results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedResult {
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, "test-user", result.User.Login)
					if callCount > 1 {
						assert.Len(t, result.Repositories, 2)
					}
				}
			} else {
				if !tc.expectedError {
					assert.Nil(t, result)
				}
			}
		})
	}
}

// TestGetRepository tests the GetRepository function
func TestGetRepository(t *testing.T) {
	// Create config for testing
	cfg := &config.Config{
		GitHub: config.GitHubConfig{
			Token:    "test-token",
			Username: "test-user",
		},
	}

	// Create test cases
	tests := []struct {
		name           string
		repoName       string
		response       string
		statusCode     int
		expectedError  bool
		expectedResult bool
	}{
		{
			name:           "Success",
			repoName:       "test-repo",
			response:       `{"name": "test-repo", "full_name": "test-user/test-repo", "description": "Test repository"}`,
			statusCode:     http.StatusOK,
			expectedError:  false,
			expectedResult: true,
		},
		{
			name:           "Repository Not Found",
			repoName:       "non-existent-repo",
			response:       `{"message": "Not Found"}`,
			statusCode:     http.StatusNotFound,
			expectedError:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a service instance
			// Need to type assert since we're now using an interface return type
			service := NewGitHubService(cfg).(*GitHubService)

			// Replace the http client's transport with our mock
			service.client.Transport = &mockTransport{
				mockResponse: func(req *http.Request) (*http.Response, error) {
					// Create a mock response
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(strings.NewReader(tc.response)),
						Header:     make(http.Header),
					}, nil
				},
			}

			// Call the function
			result, err := service.GetRepository(tc.repoName)

			// Check the results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedResult {
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, tc.repoName, result.Name)
					assert.Equal(t, "test-user/"+tc.repoName, result.FullName)
				}
			} else {
				if !tc.expectedError {
					assert.Nil(t, result)
				}
			}
		})
	}
}

// TestCreateIssue tests the CreateIssue function
func TestCreateIssue(t *testing.T) {
	// Create config for testing
	cfg := &config.Config{
		GitHub: config.GitHubConfig{
			Token:    "test-token",
			Username: "test-user",
		},
	}

	// Create test cases
	tests := []struct {
		name           string
		repoName       string
		issueRequest   *models.IssueRequest
		response       string
		statusCode     int
		expectedError  bool
		expectedResult bool
	}{
		{
			name:     "Success",
			repoName: "test-repo",
			issueRequest: &models.IssueRequest{
				Title: "Test Issue",
				Body:  "This is a test issue",
			},
			response:       `{"html_url": "https://github.com/test-user/test-repo/issues/1", "number": 1, "title": "Test Issue"}`,
			statusCode:     http.StatusCreated,
			expectedError:  false,
			expectedResult: true,
		},
		{
			name:     "Repository Not Found",
			repoName: "non-existent-repo",
			issueRequest: &models.IssueRequest{
				Title: "Test Issue",
				Body:  "This is a test issue",
			},
			response:       `{"message": "Not Found"}`,
			statusCode:     http.StatusNotFound,
			expectedError:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a service instance
			// Need to type assert since we're now using an interface return type
			service := NewGitHubService(cfg).(*GitHubService)

			// Replace the http client's transport with our mock
			service.client.Transport = &mockTransport{
				mockResponse: func(req *http.Request) (*http.Response, error) {
					// Create a mock response
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(strings.NewReader(tc.response)),
						Header:     make(http.Header),
					}, nil
				},
			}

			// Call the function
			result, err := service.CreateIssue(tc.repoName, tc.issueRequest)

			// Check the results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedResult {
				assert.NotNil(t, result)
				if result != nil {
					assert.Equal(t, 1, result.Number)
					assert.Equal(t, "Test Issue", result.Title)
					assert.Contains(t, result.HTMLURL, "issues/1")
				}
			} else {
				if !tc.expectedError {
					assert.Nil(t, result)
				}
			}
		})
	}
}
