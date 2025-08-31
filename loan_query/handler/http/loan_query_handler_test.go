package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"billing-engine/global"
	mocks "billing-engine/loan_query/_mock"
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

func TestLoanQueryHandler_GetOutstandingBalance_Success(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	expectedResponse := &models.OutstandingBalanceResponse{
		LoanID:     "loan_123456789",
		CustomerID: "12312312",
		LoanDetails: models.LoanDetailsResponse{
			InstallmentUnit:   "week",
			InstallmentAmount: 110000.00,
			TotalInstallments: 50,
		},
		OutstandingAmount: 3300000.00,
		OverdueAmount:     220000.00,

		OverdueInstallments:   2,
		PaidInstallments:      20,
		RemainingInstallments: 30,
	}

	mockService.On("GetOutstandingBalance", mock.Anything, "loan_123456789").Return(expectedResponse, nil)

	// Create request
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans/loan_123456789/outstanding", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("loan_123456789")

	// Execute
	err := handler.GetOutstandingBalance(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response global.OutstandingBalanceSuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestLoanQueryHandler_GetOutstandingBalance_EmptyLoanID(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	// Create request with empty loan_id
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans//outstanding", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("")

	// Execute
	err := handler.GetOutstandingBalance(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Equal(t, "Loan ID is required", response["message"])
}

func TestLoanQueryHandler_GetOutstandingBalance_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	mockService.On("GetOutstandingBalance", mock.Anything, "loan_123456789").Return(nil, errors.New("loan not found"))

	// Create request
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans/loan_123456789/outstanding", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("loan_123456789")

	// Execute
	err := handler.GetOutstandingBalance(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusInternalServerError), response["code"])
	assert.Equal(t, "loan not found", response["message"])

	mockService.AssertExpectations(t)
}

func TestLoanQueryHandler_GetDelinquencyStatus_Success(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	expectedResponse := &models.DelinquencyResponse{
		LoanID:       "loan_123456789",
		CustomerID:   "12312312",
		IsDelinquent: true,

		InstallmentUnit:       "week",
		OverdueInstallments:   2,
		OverdueAmount:         220000.00,
		OutstandingAmount:     3300000.00,
		RequiredPaymentAmount: 220000.00,
	}

	mockService.On("GetDelinquencyStatus", mock.Anything, "loan_123456789").Return(expectedResponse, nil)

	// Create request
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans/loan_123456789/delinquency", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("loan_123456789")

	// Execute
	err := handler.GetDelinquencyStatus(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response global.DelinquencySuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestLoanQueryHandler_GetLoanSchedule_Success(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	paidDate := time.Now().AddDate(0, 0, -7)
	expectedResponse := &models.LoanScheduleResponse{
		LoanID: "loan_123456789",
		LoanSummary: models.LoanSummaryScheduleResponse{
			InstallmentUnit:   "week",
			TotalInstallments: 50,
			InstallmentAmount: 110000.00,
			DisbursedAmount:   5000000.00,
			OutstandingAmount: 3300000.00,
		},
		Schedule: []models.PaymentScheduleResponse{
			{
				InstallmentNumber: 1,
				DueDate:           time.Now().AddDate(0, 0, -14),
				InstallmentAmount: 110000.00,
				InstallmentPaid:   5390000.00,
				Status:            "PAID",
				PaidDate:          &paidDate,
			},
			{
				InstallmentNumber: 2,
				DueDate:           time.Now().AddDate(0, 0, -7),
				InstallmentAmount: 110000.00,
				InstallmentPaid:   5280000.00,
				Status:            "PENDING",
				PaidDate:          nil,
			},
		},
	}

	mockService.On("GetLoanSchedule", mock.Anything, "loan_123456789").Return(expectedResponse, nil)

	// Create request
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans/loan_123456789/schedule", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("loan_123456789")

	// Execute
	err := handler.GetLoanSchedule(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response global.LoanScheduleSuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestLoanQueryHandler_GetDelinquencyStatus_EmptyLoanID(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	// Create request with empty loan_id
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans//delinquency", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("")

	// Execute
	err := handler.GetDelinquencyStatus(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Equal(t, "Loan ID is required", response["message"])
}

func TestLoanQueryHandler_GetLoanSchedule_EmptyLoanID(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewLoanQueryServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &LoanQueryHandler{
		loanQueryService: mockService,
		middleware:       mockMiddleware,
	}

	// Create request with empty loan_id
	httpReq := httptest.NewRequest(http.MethodGet, "/v1/loans//schedule", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)
	c.SetParamNames("loan_id")
	c.SetParamValues("")

	// Execute
	err := handler.GetLoanSchedule(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Equal(t, "Loan ID is required", response["message"])
}
