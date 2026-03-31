package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountHandler_CreateAccount(t *testing.T) {
	type Mocker struct {
		service *mocks.AccountService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		body           interface{}
		expectedStatus int
		expectedCode   int
	}{
		{
			desc: "SUCCESS: Create Account",
			body: dto.CreateAccountRequest{
				AccountHolder: "Ayash",
			},
			mockSetup: func(m *Mocker) {
				m.service.On(
					"CreateAccount",
					mock.Anything,
					mock.AnythingOfType("dto.CreateAccountRequest"),
				).Return(dto.AccountResponse{
					ID:            "123",
					AccountHolder: "Ayash",
					Balance:       0,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   http.StatusCreated,
		},
		{
			desc: "ERROR: Missing Account Holder",
			body: dto.CreateAccountRequest{
				AccountHolder: "",
			},
			mockSetup:      func(m *Mocker) {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   http.StatusBadRequest,
		},
		{
			desc: "ERROR: Service Returns Error",
			body: dto.CreateAccountRequest{
				AccountHolder: "Ayash",
			},
			mockSetup: func(m *Mocker) {
				m.service.On(
					"CreateAccount",
					mock.Anything,
					mock.AnythingOfType("dto.CreateAccountRequest"),
				).Return(dto.AccountResponse{}, assert.AnError)
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
			svc := new(mocks.AccountService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewAccountHandler(mux, svc)

			var bodyBytes []byte
			switch v := tc.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handler.CreateAccount(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp dto.BaseResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tc.expectedCode, resp.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestAccountHandler_GetAllAccounts(t *testing.T) {
	type Mocker struct {
		service *mocks.AccountService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		expectedStatus int
		expectedCode   int
	}{
		{
			desc: "SUCCESS: Get All Accounts",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAllAccounts",
					mock.Anything,
				).Return([]dto.AccountResponse{
					{ID: "123", AccountHolder: "Ayash", Balance: 1000},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   http.StatusOK,
		},
		{
			desc: "ERROR: Service Returns Error",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAllAccounts",
					mock.Anything,
				).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			svc := new(mocks.AccountService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewAccountHandler(mux, svc)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			rec := httptest.NewRecorder()

			handler.GetAllAccounts(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp dto.BaseResponse
			json.NewDecoder(rec.Body).Decode(&resp)
			assert.Equal(t, tc.expectedCode, resp.Code)

			svc.AssertExpectations(t)
		})
	}
}

func TestAccountHandler_GetAccountByID(t *testing.T) {
	type Mocker struct {
		service *mocks.AccountService
	}

	testCases := []struct {
		desc           string
		mockSetup      func(m *Mocker)
		pathID         string
		expectedStatus int
		expectedCode   int
	}{
		{
			desc:   "SUCCESS: Get Account by ID",
			pathID: "123",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAccountByID",
					mock.Anything,
					"123",
				).Return(dto.AccountResponse{
					ID:            "123",
					AccountHolder: "Ayash",
					Balance:       1000,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   http.StatusOK,
		},
		{
			desc:   "ERROR: Account Not Found",
			pathID: "999",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAccountByID",
					mock.Anything,
					"999",
				).Return(dto.AccountResponse{}, errors.New("account not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   http.StatusNotFound,
		},
		{
			desc:   "ERROR: Server Error",
			pathID: "999",
			mockSetup: func(m *Mocker) {
				m.service.On(
					"GetAccountByID",
					mock.Anything,
					"999",
				).Return(dto.AccountResponse{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			svc := new(mocks.AccountService)
			mocker := &Mocker{service: svc}

			tc.mockSetup(mocker)

			mux := http.NewServeMux()
			handler := NewAccountHandler(mux, svc)

			mux.HandleFunc("GET /accounts/{id}", handler.GetAccountByID)

			req := httptest.NewRequest(http.MethodGet, "/accounts/"+tc.pathID, nil)
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
