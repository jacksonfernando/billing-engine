package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "billing-engine/disbursement/_mock"
	"billing-engine/global"
	"billing-engine/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMiddleware is a mock implementation of GoMiddlewareInterface
type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) ValidateCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return next
}

func (m *MockMiddleware) ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return next
}

func TestDisbursementHandler_CreateDisbursement_Success(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewDisbursementServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &DisbursementHandler{
		disbursementService: mockService,
		middleware:          mockMiddleware,
	}

	req := models.DisbursementRequest{
		PrincipalAmount:     5000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
		CustomerID:          "12312312",
	}

	expectedResponse := &models.DisbursementResponse{
		LoanID:              "loan_123456789",
		CustomerID:          "12312312",
		DisbursedAmount:     5000000.00,
		InstallmentAmount:   110000.00,
		OutstandingAmount:   5500000.00,
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		DisbursementDate:    time.Now(),
		FirstDueDate:        time.Now().AddDate(0, 0, 7),
		FinalDueDate:        time.Now().AddDate(0, 0, 350),
	}

	mockService.On("CreateDisbursement", mock.Anything, &req).Return(expectedResponse, nil)

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/disbursement", bytes.NewBuffer(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.CreateDisbursement(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response global.DisbursementSuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestDisbursementHandler_CreateDisbursement_InvalidJSON(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewDisbursementServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &DisbursementHandler{
		disbursementService: mockService,
		middleware:          mockMiddleware,
	}

	// Create invalid JSON request
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/disbursement", bytes.NewBuffer([]byte("invalid json")))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.CreateDisbursement(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response global.BadResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request body", response.Message)
}

func TestDisbursementHandler_CreateDisbursement_ValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		request       models.DisbursementRequest
		expectedError string
	}{
		{
			name: "Zero principal amount",
			request: models.DisbursementRequest{
				PrincipalAmount:     0,
				InterestRate:        0.10, // 10% interest rate
				InstallmentUnit:     "week",
				NumberOfInstallment: 50,
				StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
				CustomerID:          "12312312",
			},
			expectedError: "principalamount must be greater than 0",
		},
		{
			name: "Zero interest rate",
			request: models.DisbursementRequest{
				PrincipalAmount:     5000000.00,
				InterestRate:        0,
				InstallmentUnit:     "week",
				NumberOfInstallment: 50,
				StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
				CustomerID:          "12312312",
			},
			expectedError: "interestrate must be greater than 0",
		},
		{
			name: "Zero installments",
			request: models.DisbursementRequest{
				PrincipalAmount:     5000000.00,
				InterestRate:        0.10, // 10% interest rate
				InstallmentUnit:     "week",
				NumberOfInstallment: 0,
				StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
				CustomerID:          "12312312",
			},
			expectedError: "numberofinstallment must be greater than 0",
		},
		{
			name: "Invalid installment unit",
			request: models.DisbursementRequest{
				PrincipalAmount:     5000000.00,
				InterestRate:        0.10, // 10% interest rate
				InstallmentUnit:     "daily",
				NumberOfInstallment: 50,
				StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
				CustomerID:          "12312312",
			},
			expectedError: "installmentunit must be one of: week, month",
		},
		{
			name: "Empty customer ID",
			request: models.DisbursementRequest{
				PrincipalAmount:     5000000.00,
				InterestRate:        0.10, // 10% interest rate
				InstallmentUnit:     "week",
				NumberOfInstallment: 50,
				StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
				CustomerID:          "",
			},
			expectedError: "customerid is required",
		},
		{
			name: "Empty start date",
			request: models.DisbursementRequest{
				PrincipalAmount:     5000000.00,
				InterestRate:        0.10, // 10% interest rate
				InstallmentUnit:     "week",
				NumberOfInstallment: 50,
				StartDate:           time.Time{},
				CustomerID:          "12312312",
			},
			expectedError: "startdate is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockService := mocks.NewDisbursementServiceInterface(t)
			mockMiddleware := new(MockMiddleware)

			handler := &DisbursementHandler{
				disbursementService: mockService,
				middleware:          mockMiddleware,
			}

			// Create request
			reqBody, _ := json.Marshal(tt.request)
			httpReq := httptest.NewRequest(http.MethodPost, "/v1/disbursement", bytes.NewBuffer(reqBody))
			httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(httpReq, rec)

			// Execute
			err := handler.CreateDisbursement(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			var response global.BadResponse
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, tt.expectedError, response.Message)
		})
	}
}

func TestDisbursementHandler_CreateDisbursement_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewDisbursementServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &DisbursementHandler{
		disbursementService: mockService,
		middleware:          mockMiddleware,
	}

	req := models.DisbursementRequest{
		PrincipalAmount:     5000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
		CustomerID:          "12312312",
	}

	mockService.On("CreateDisbursement", mock.Anything, &req).Return(nil, errors.New("service error"))

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/disbursement", bytes.NewBuffer(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.CreateDisbursement(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response global.BadResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "service error", response.Message)

	mockService.AssertExpectations(t)
}
