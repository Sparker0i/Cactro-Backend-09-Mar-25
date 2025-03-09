package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/config"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/api/routes"
	"github.com/Sparker0i/Cactro-Backend-09-Mar-25/internal/models"
	"github.com/stretchr/testify/assert"
)

// setupIntegrationTestEnv sets up the environment for integration tests
func setupIntegrationTestEnv() {
	// Set environment variables for testing
	os.Setenv("GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_USERNAME", "test-user")
	os.Setenv("GIN_MODE", "test")
}

// TestAPIIntegration tests the API endpoints together
func TestAPIIntegration(t *testing.T) {
	// Skip this test during automated testing
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("Skipping integration test")
	}

	// Set up the test environment
	setupIntegrationTestEnv()

	// Load the configuration
	cfg := config.LoadConfig()

	// Set up the router
	router := routes.SetupRouter(cfg)

	// Test health check endpoint
	t.Run("Health Check", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response map[string]string
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})

	// Test GitHub profile endpoint
	if os.Getenv("GITHUB_TOKEN") != "test-token" {
		t.Run("Get GitHub Profile", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/github", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			// Should return a success status code if the API is working
			assert.True(t, resp.Code == http.StatusOK || resp.Code == http.StatusUnauthorized)

			// If successful, check the response structure
			if resp.Code == http.StatusOK {
				var profile models.GithubProfile
				err := json.Unmarshal(resp.Body.Bytes(), &profile)
				assert.NoError(t, err)
				assert.NotEmpty(t, profile.User.Login)
			}
		})
	}

	// Test Create Issue endpoint with invalid input
	t.Run("Create Issue - Invalid Input", func(t *testing.T) {
		// Create an invalid issue request (missing required fields)
		invalidIssue := map[string]string{
			"title": "", // Empty title
		}

		requestBody, _ := json.Marshal(invalidIssue)
		req := httptest.NewRequest("POST", "/github/test-repo/issues", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)

		var errorResponse models.ErrorResponse
		err := json.Unmarshal(resp.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.NotEmpty(t, errorResponse.Error)
	})

	// Test non-existent endpoint
	t.Run("Non-existent Endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/non-existent", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
