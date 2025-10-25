package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLogging(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})

	middleware := Logging(handler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rec.Code)
	}

	if rec.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", rec.Body.String())
	}
}

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		corsEnv        string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "GET request with default CORS",
			method:         "GET",
			corsEnv:        "",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request with custom CORS",
			method:         "POST",
			corsEnv:        "https://example.com",
			expectedOrigin: "https://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS preflight request",
			method:         "OPTIONS",
			corsEnv:        "*",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.corsEnv != "" {
				os.Setenv("CORS", tt.corsEnv)
			} else {
				os.Unsetenv("CORS")
			}
			defer os.Unsetenv("CORS")

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte("OK")); err != nil {
					t.Errorf("Failed to write response: %v", err)
				}
			})

			middleware := CORS(handler)

			req := httptest.NewRequest(tt.method, "/api/test", nil)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			origin := rec.Header().Get("Access-Control-Allow-Origin")
			if origin != tt.expectedOrigin {
				t.Errorf("Expected Access-Control-Allow-Origin to be '%s', got '%s'", tt.expectedOrigin, origin)
			}

			methods := rec.Header().Get("Access-Control-Allow-Methods")
			if methods == "" {
				t.Error("Expected Access-Control-Allow-Methods header to be set")
			}

			headers := rec.Header().Get("Access-Control-Allow-Headers")
			if headers == "" {
				t.Error("Expected Access-Control-Allow-Headers header to be set")
			}
		})
	}
}

func TestCORSChaining(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Chained")); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	})

	// Chain CORS and Logging
	chained := CORS(Logging(handler))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()

	chained.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rec.Code)
	}

	if rec.Body.String() != "Chained" {
		t.Errorf("Expected body 'Chained', got '%s'", rec.Body.String())
	}

	// Verify CORS headers are set
	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin == "" {
		t.Error("Expected CORS headers to be set in chained middleware")
	}
}
