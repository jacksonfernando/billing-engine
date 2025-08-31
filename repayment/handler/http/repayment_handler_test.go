package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"billing-engine/global"
	"billing-engine/models"
	mocks "billing-engine/repayment/_mock"

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

func TestRepaymentHandler_ProcessRepayment_Success(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewRepaymentServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &RepaymentHandler{
		repaymentService: mockService,
		middleware:       mockMiddleware,
	}

	req := models.RepaymentRequest{
		LoanID:        "loan_123456789",
		PaymentAmount: 220000.00,
	}

	expectedResponse := &models.RepaymentResponse{
		LoanID:                "loan_123456789",
		PaymentAmount:         220000.00,
		InstallmentsPaid:      2,
		InstallmentAmount:     110000.00,
		RemainingInstallments: 48,
		OutstandingAmount:     5280000.00,
		NextDueDate:           time.Now().AddDate(0, 0, 7),
		PaymentDate:           time.Now(),
	}

	mockService.On("ProcessRepayment", mock.Anything, &req).Return(expectedResponse, nil)

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/repayment", bytes.NewBuffer(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.ProcessRepayment(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response global.RepaymentSuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestRepaymentHandler_ProcessRepayment_InvalidJSON(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewRepaymentServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &RepaymentHandler{
		repaymentService: mockService,
		middleware:       mockMiddleware,
	}

	// Create invalid JSON request
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/repayment", bytes.NewBuffer([]byte("invalid json")))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.ProcessRepayment(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response global.BadResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Invalid request body", response.Message)
}

func TestRepaymentHandler_ProcessRepayment_ValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		request       models.RepaymentRequest
		expectedError string
	}{
		{
			name: "Empty loan ID",
			request: models.RepaymentRequest{
				LoanID:        "",
				PaymentAmount: 220000.00,
			},
			expectedError: "loanid is required",
		},
		{
			name: "Zero payment amount",
			request: models.RepaymentRequest{
				LoanID:        "loan_123456789",
				PaymentAmount: 0,
			},
			expectedError: "paymentamount must be greater than 0",
		},
		{
			name: "Negative payment amount",
			request: models.RepaymentRequest{
				LoanID:        "loan_123456789",
				PaymentAmount: -100000.00,
			},
			expectedError: "paymentamount must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockService := mocks.NewRepaymentServiceInterface(t)
			mockMiddleware := new(MockMiddleware)

			handler := &RepaymentHandler{
				repaymentService: mockService,
				middleware:       mockMiddleware,
			}

			// Create request
			reqBody, _ := json.Marshal(tt.request)
			httpReq := httptest.NewRequest(http.MethodPost, "/v1/repayment", bytes.NewBuffer(reqBody))
			httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(httpReq, rec)

			// Execute
			err := handler.ProcessRepayment(c)

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

func TestRepaymentHandler_ProcessRepayment_ServiceError(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewRepaymentServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &RepaymentHandler{
		repaymentService: mockService,
		middleware:       mockMiddleware,
	}

	req := models.RepaymentRequest{
		LoanID:        "loan_123456789",
		PaymentAmount: 220000.00,
	}

	mockService.On("ProcessRepayment", mock.Anything, &req).Return(nil, errors.New("loan not found"))

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/repayment", bytes.NewBuffer(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.ProcessRepayment(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response global.BadResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "loan not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestRepaymentHandler_ProcessRepayment_IncorrectPaymentAmount(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := mocks.NewRepaymentServiceInterface(t)
	mockMiddleware := new(MockMiddleware)

	handler := &RepaymentHandler{
		repaymentService: mockService,
		middleware:       mockMiddleware,
	}

	req := models.RepaymentRequest{
		LoanID:        "loan_123456789",
		PaymentAmount: 100000.00,
	}

	mockService.On("ProcessRepayment", mock.Anything, &req).Return(nil, errors.New("payment amount 100000.00 does not match required amount 220000.00"))

	// Create request
	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/v1/repayment", bytes.NewBuffer(reqBody))
	httpReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(httpReq, rec)

	// Execute
	err := handler.ProcessRepayment(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response global.BadResponse
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Contains(t, response.Message, "payment amount 100000.00 does not match required amount 220000.00")

	mockService.AssertExpectations(t)
}
