package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIPath(t *testing.T) {
	testCases := []struct {
		desc     string
		method   string
		path     string
		expected string
	}{
		{
			desc:     "SUCCESS: GET Path",
			method:   http.MethodGet,
			path:     "/persons",
			expected: "GET /persons",
		},
		{
			desc:     "SUCCESS: POST Path",
			method:   http.MethodPost,
			path:     "/persons",
			expected: "POST /persons",
		},
		{
			desc:     "SUCCESS: GET Path with ID",
			method:   http.MethodGet,
			path:     "/persons/{id}",
			expected: "GET /persons/{id}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := NewAPIPath(tc.method, tc.path)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestApplicationMiddlewareResponse(t *testing.T) {
	testCases := []struct {
		desc            string
		expectedHeader  string
	}{
		{
			desc:           "SUCCESS: Sets Content-Type Header",
			expectedHeader: "application/json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			innerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := ApplicationMiddlewareResponse(innerHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedHeader, rec.Header().Get("Content-Type"))
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestHandleRoutesNotFound(t *testing.T) {
	testCases := []struct {
		desc           string
		path           string
		registerRoute  bool
		expectedStatus int
		expectedBody   string
	}{
		{
			desc:           "SUCCESS: Route Exists",
			path:           "/exists",
			registerRoute:  true,
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			desc:           "ERROR: Route Not Found",
			path:           "/not-found",
			registerRoute:  false,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":"error","message":"Route not found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			mux := http.NewServeMux()

			if tc.registerRoute {
				mux.HandleFunc(tc.path, func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("OK"))
				})
			}

			handler := HandleRoutesNotFound(mux)

			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedBody, rec.Body.String())
		})
	}
}
