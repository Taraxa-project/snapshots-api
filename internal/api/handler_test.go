package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taraxa/snapshots-api/internal/auth"
	"github.com/taraxa/snapshots-api/internal/config"
	"github.com/taraxa/snapshots-api/internal/models"
)

func createTestHandler(apiKeys []string) (*Handler, *MockSnapshotService) {
	mockService := &MockSnapshotService{}
	cfg := &config.Config{APIKeys: apiKeys}
	authMiddleware := auth.NewMiddleware(cfg)
	handler := NewHandler(mockService, authMiddleware)
	return handler, mockService
}

func TestHandler_GetSnapshots(t *testing.T) {
	handler, mockService := createTestHandler([]string{"valid-api-key"})

	tests := []struct {
		name           string
		queryParams    string
		authHeader     string
		expectedStatus int
		checkResponse  bool
		mockError      error
		checkFullData  bool // Check if full snapshots are included
	}{
		{
			name:           "authenticated request - should get full and light snapshots",
			queryParams:    "?network=mainnet",
			authHeader:     "Bearer valid-api-key",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
			checkFullData:  true,
			mockError:      nil,
		},
		{
			name:           "unauthenticated request - should get only light snapshots",
			queryParams:    "?network=mainnet",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
			checkFullData:  false,
			mockError:      nil,
		},
		{
			name:           "invalid API key - should get only light snapshots",
			queryParams:    "?network=mainnet",
			authHeader:     "Bearer invalid-key",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
			checkFullData:  false,
			mockError:      nil,
		},
		{
			name:           "malformed auth header - should get only light snapshots",
			queryParams:    "?network=mainnet",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
			checkFullData:  false,
			mockError:      nil,
		},
		{
			name:           "missing network parameter",
			queryParams:    "",
			authHeader:     "",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
			checkFullData:  false,
			mockError:      nil,
		},
		{
			name:           "invalid network",
			queryParams:    "?network=invalid",
			authHeader:     "",
			expectedStatus: http.StatusBadRequest,
			checkResponse:  false,
			checkFullData:  false,
			mockError:      nil,
		},
		{
			name:           "service error with auth",
			queryParams:    "?network=mainnet",
			authHeader:     "Bearer valid-api-key",
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  false,
			checkFullData:  false,
			mockError:      errors.New("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			if tt.mockError != nil {
				mockService.GetSnapshotsWithAuthFunc = func(network models.Network, authenticated bool) (*models.NetworkSnapshots, error) {
					return nil, tt.mockError
				}
			} else {
				mockService.GetSnapshotsWithAuthFunc = nil // Use default
			}

			req, err := http.NewRequest("GET", "/"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Add auth header if provided
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()
			handler.getSnapshots(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.checkResponse && tt.expectedStatus == http.StatusOK {
				// Check if response is valid JSON
				var result models.NetworkSnapshots
				if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				// Check authentication filtering
				if tt.checkFullData {
					// Should have both full and light snapshots
					if result.Full == nil {
						t.Errorf("Expected full snapshot data for authenticated request, got nil")
					}
					if result.Light == nil {
						t.Errorf("Expected light snapshot data, got nil")
					}
					// Should have previous-full for authenticated requests
					if result.PreviousFull == nil {
						t.Errorf("Expected previous-full data for authenticated request, got nil")
					}
					// Should have previous-light for all requests
					if result.PreviousLight == nil {
						t.Errorf("Expected previous-light data, got nil")
					}
				} else {
					// Should have only light snapshots
					if result.Full != nil {
						t.Errorf("Expected no full snapshot data for unauthenticated request, got %+v", result.Full)
					}
					if result.Light == nil {
						t.Errorf("Expected light snapshot data, got nil")
					}
					// Should NOT have previous-full for unauthenticated requests
					if result.PreviousFull != nil {
						t.Errorf("Expected no previous-full data for unauthenticated request, got %+v", result.PreviousFull)
					}
					// Should have previous-light for all requests
					if result.PreviousLight == nil {
						t.Errorf("Expected previous-light data, got nil")
					}
				}

				// Check content type
				expectedContentType := "application/json"
				if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
					t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
				}

				// Check cache control header
				if cacheControl := rr.Header().Get("Cache-Control"); cacheControl == "" {
					t.Errorf("handler should set Cache-Control header")
				}
			}
		})
	}
}

func TestHandler_GetSnapshots_InvalidMethods(t *testing.T) {
	handler, _ := createTestHandler([]string{})

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/?network=mainnet", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.getSnapshots(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for %s: got %v want %v", method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestHandler_Health(t *testing.T) {
	handler, _ := createTestHandler([]string{})

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.health(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if response["service"] != "snapshots-api" {
		t.Errorf("Expected service 'snapshots-api', got %v", response["service"])
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}
}

func TestHandler_Ready(t *testing.T) {
	handler, mockService := createTestHandler([]string{})

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "service ready",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "ready",
		},
		{
			name:           "service not ready",
			mockError:      errors.New("connection failed"),
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "not ready",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock behavior
			if tt.mockError != nil {
				mockService.GetSnapshotsFunc = func(network models.Network) (*models.NetworkSnapshots, error) {
					return nil, tt.mockError
				}
			} else {
				mockService.GetSnapshotsFunc = nil // Use default
			}

			req, err := http.NewRequest("GET", "/ready", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ready(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check response body
			var response map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			if response["status"] != tt.expectedBody {
				t.Errorf("Expected status '%s', got %v", tt.expectedBody, response["status"])
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
			}
		})
	}
}

func TestHandler_Health_InvalidMethods(t *testing.T) {
	handler, _ := createTestHandler([]string{})

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.health(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for %s: got %v want %v", method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestHandler_Ready_InvalidMethods(t *testing.T) {
	handler, _ := createTestHandler([]string{})

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, "/ready", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ready(rr, req)

			if status := rr.Code; status != http.StatusMethodNotAllowed {
				t.Errorf("handler returned wrong status code for %s: got %v want %v", method, status, http.StatusMethodNotAllowed)
			}
		})
	}
}

func TestHandler_Routes(t *testing.T) {
	handler, _ := createTestHandler([]string{})

	routes := handler.Routes()
	if routes == nil {
		t.Error("Routes() returned nil")
	}

	// Test that routes are properly configured by making requests
	endpoints := []string{"/", "/health", "/ready"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			var req *http.Request
			var err error

			if endpoint == "/" {
				req, err = http.NewRequest("GET", "/?network=mainnet", nil)
			} else {
				req, err = http.NewRequest("GET", endpoint, nil)
			}

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			routes.ServeHTTP(rr, req)

			// Should not return 404 (endpoint exists)
			if rr.Code == http.StatusNotFound {
				t.Errorf("Endpoint %s not found", endpoint)
			}
		})
	}
}
