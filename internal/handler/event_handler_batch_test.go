package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mohamedelhefni/siraaj/geolocation"
	"github.com/mohamedelhefni/siraaj/internal/domain"
	"github.com/mohamedelhefni/siraaj/internal/mocks"
)

func TestTrackBatchEvents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEventService(ctrl)
	geoService, _ := geolocation.NewService()
	handler := NewEventHandler(mockService, geoService)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func()
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful batch processing",
			requestBody: map[string]interface{}{
				"events": []domain.Event{
					{
						EventName: "page_view",
						UserID:    "user1",
						SessionID: "session1",
						URL:       "https://example.com",
						Timestamp: time.Now(),
					},
					{
						EventName: "button_click",
						UserID:    "user2",
						SessionID: "session2",
						URL:       "https://example.com/page",
						Timestamp: time.Now(),
					},
				},
			},
			setupMock: func() {
				mockService.EXPECT().
					TrackEventBatch(gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["status"] != "ok" {
					t.Errorf("Expected status ok, got %v", body["status"])
				}
				if body["successful"].(float64) != 2 {
					t.Errorf("Expected 2 successful events, got %v", body["successful"])
				}
				if body["failed"].(float64) != 0 {
					t.Errorf("Expected 0 failed events, got %v", body["failed"])
				}
			},
		},
		{
			name:           "empty events array",
			requestBody:    map[string]interface{}{"events": []domain.Event{}},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "no events field",
			requestBody:    map[string]interface{}{},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name:           "invalid json",
			requestBody:    "invalid json",
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "batch processing with error",
			requestBody: map[string]interface{}{
				"events": []domain.Event{
					{
						EventName: "page_view",
						UserID:    "user1",
						SessionID: "session1",
						URL:       "https://example.com",
						Timestamp: time.Now(),
					},
					{
						EventName: "button_click",
						UserID:    "user2",
						SessionID: "session2",
						URL:       "https://example.com/page",
						Timestamp: time.Now(),
					},
				},
			},
			setupMock: func() {
				// Batch insert fails
				mockService.EXPECT().
					TrackEventBatch(gomock.Any()).
					Return(gomock.Errorf("database error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/track/batch", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.TrackBatchEvents(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil && w.Code == http.StatusOK || w.Code == http.StatusPartialContent {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestTrackBatchEvents_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEventService(ctrl)
	handler := NewEventHandler(mockService, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/track/batch", nil)
	w := httptest.NewRecorder()

	handler.TrackBatchEvents(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestTrackBatchEvents_MaxBatchSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEventService(ctrl)
	handler := NewEventHandler(mockService, nil)

	// Create 101 events (exceeds max of 100)
	events := make([]domain.Event, 101)
	for i := 0; i < 101; i++ {
		events[i] = domain.Event{
			EventName: "test_event",
			UserID:    "user1",
			SessionID: "session1",
			Timestamp: time.Now(),
		}
	}

	requestBody := map[string]interface{}{
		"events": events,
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/track/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.TrackBatchEvents(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for batch size exceeded, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTrackBatchEvents_WithGeolocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockEventService(ctrl)
	geoService, _ := geolocation.NewService()
	handler := NewEventHandler(mockService, geoService)

	requestBody := map[string]interface{}{
		"events": []domain.Event{
			{
				EventName: "page_view",
				UserID:    "user1",
				SessionID: "session1",
				IP:        "8.8.8.8", // Google DNS - should resolve to US
				Timestamp: time.Now(),
			},
		},
	}

	mockService.EXPECT().
		TrackEventBatch(gomock.Any()).
		DoAndReturn(func(events []domain.Event) error {
			// Verify geolocation was enriched for the first event
			if len(events) > 0 && events[0].Country == "" {
				t.Error("Expected country to be enriched from IP")
			}
			return nil
		}).
		Times(1)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/track/batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.TrackBatchEvents(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
