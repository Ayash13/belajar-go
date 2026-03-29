package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"practice_04_unit_testing/dto"
	"practice_04_unit_testing/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersonHandler_CreatePerson(t *testing.T) {
	type Mocker struct {
		service *mocks.PersonService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		body           interface{}
		expectedStatus int
		expectedCode   int
	}{
		{
			desc: "SUCCESS: Create Person",
			body: dto.PersonCreateRequest{
				Name:  "ay",
				Email: "ay@mail.com",
			},
			mockSetup: func(m *Mocker) {
				m.service.On(
					"CreatePerson",
					mock.Anything,
					mock.AnythingOfType("dto.PersonCreateRequest"),
				).Return(dto.PersonResponse{
					ID:    1,
					Name:  "ay",
					Email: "ay@mail.com",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   http.StatusCreated,
		},
		{
			desc: "ERROR: Service Returns Error",
			body: dto.PersonCreateRequest{
				Name:  "ay",
				Email: "ay@mail.com",
			},
			mockSetup: func(m *Mocker) {
				m.service.On(
					"CreatePerson",
					mock.Anything,
					mock.AnythingOfType("dto.PersonCreateRequest"),
				).Return(dto.PersonResponse{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   http.StatusInternalServerError,
		},
		{
			desc:           "ERROR: Invalid JSON Body",
			body:           "invalid-json",
			mockSetup:      func(m *Mocker) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			svc := new(mocks.PersonService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewPersonHandler(mux, svc)

			var bodyBytes []byte
			switch v := tc.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/persons", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.CreatePerson(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp dto.BaseResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tc.expectedCode, resp.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestPersonHandler_GetAllPersons(t *testing.T) {
	type Mocker struct {
		service *mocks.PersonService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		expectedStatus int
		expectedCode   int
	}{
		{
			desc: "SUCCESS: Get All Persons",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAllPersons",
					mock.Anything,
				).Return([]dto.PersonResponse{
					{ID: 1, Name: "ay", Email: "ay@mail.com"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   http.StatusOK,
		},
		{
			desc: "ERROR: Service Returns Error",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAllPersons",
					mock.Anything,
				).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			svc := new(mocks.PersonService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewPersonHandler(mux, svc)

			req := httptest.NewRequest(http.MethodGet, "/persons", nil)
			rec := httptest.NewRecorder()

			handler.GetAllPersons(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp dto.BaseResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tc.expectedCode, resp.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestPersonHandler_GetPerson(t *testing.T) {
	type Mocker struct {
		service *mocks.PersonService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		pathID         string
		expectedStatus int
		expectedCode   int
	}{
		{
			desc:   "SUCCESS: Get Person by ID",
			pathID: "1",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetPerson",
					mock.Anything,
					mock.AnythingOfType("int"),
				).Return(dto.PersonResponse{
					ID:    1,
					Name:  "ay",
					Email: "ay@mail.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   http.StatusOK,
		},
		{
			desc:   "ERROR: Person Not Found",
			pathID: "999",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetPerson",
					mock.Anything,
					mock.AnythingOfType("int"),
				).Return(dto.PersonResponse{}, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   http.StatusNotFound,
		},
		{
			desc:           "ERROR: Invalid ID Format",
			pathID:         "abc",
			mockSetup:      func(m *Mocker) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			svc := new(mocks.PersonService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewPersonHandler(mux, svc)

			mux.HandleFunc("GET /persons/{id}", handler.GetPerson)

			req := httptest.NewRequest(http.MethodGet, "/persons/"+tc.pathID, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp dto.BaseResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tc.expectedCode, resp.Code)

			svc.AssertExpectations(t)
		})
	}
}
