package service

import (
	mocks "billing-engine/disbursement/_mock"
	"billing-engine/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDisbursementService_CreateDisbursement_Success(t *testing.T) {
	mockRepo := mocks.NewDisbursementMySQLRepositoryInterface(t)
	service := NewDisbursementService(mockRepo)
	ctx := context.Background()

	req := &models.DisbursementRequest{
		PrincipalAmount:     5000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
		CustomerID:          "12312312",
	}

	// Mock repository calls
	mockRepo.On("CreateDisbursement", ctx, mock.AnythingOfType("*models.DisbursementDetail")).Return(nil)
	mockRepo.On("CreateLoanSummary", ctx, mock.AnythingOfType("*models.LoanSummary")).Return(nil)
	mockRepo.On("CreatePaymentSchedules", ctx, mock.AnythingOfType("[]*models.PaymentSchedule")).Return(nil)

	// Execute
	response, err := service.CreateDisbursement(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.CustomerID, response.CustomerID)
	assert.Equal(t, req.PrincipalAmount, response.DisbursedAmount) // only principal is disbursed
	assert.Equal(t, 110000.00, response.InstallmentAmount)         // (5000000 + 500000) / 50
	assert.Equal(t, 5500000.00, response.OutstandingAmount)        // 5000000 + 500000
	assert.Equal(t, req.InstallmentUnit, response.InstallmentUnit)
	assert.Equal(t, req.NumberOfInstallment, response.NumberOfInstallment)
	assert.Contains(t, response.LoanID, "loan_")

	mockRepo.AssertExpectations(t)
}

func TestDisbursementService_CreateDisbursement_ZeroStartDate(t *testing.T) {
	mockRepo := mocks.NewDisbursementMySQLRepositoryInterface(t)
	service := NewDisbursementService(mockRepo)
	ctx := context.Background()

	req := &models.DisbursementRequest{
		PrincipalAmount:     5000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		StartDate:           time.Time{}, // Zero time value
		CustomerID:          "12312312",
	}

	// Execute
	response, err := service.CreateDisbursement(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "start date cannot be zero")
}

func TestDisbursementService_CreateDisbursement_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewDisbursementMySQLRepositoryInterface(t)
	service := NewDisbursementService(mockRepo)
	ctx := context.Background()

	req := &models.DisbursementRequest{
		PrincipalAmount:     5000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "week",
		NumberOfInstallment: 50,
		StartDate:           time.Date(2025, 8, 31, 11, 43, 0, 0, time.UTC),
		CustomerID:          "12312312",
	}

	// Mock repository error
	mockRepo.On("CreateDisbursement", ctx, mock.AnythingOfType("*models.DisbursementDetail")).Return(errors.New("database error"))

	// Execute
	response, err := service.CreateDisbursement(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to create disbursement")

	mockRepo.AssertExpectations(t)
}

func TestDisbursementService_CreateDisbursement_MonthlyInstallments(t *testing.T) {
	mockRepo := mocks.NewDisbursementMySQLRepositoryInterface(t)
	service := NewDisbursementService(mockRepo)
	ctx := context.Background()

	req := &models.DisbursementRequest{
		PrincipalAmount:     1000000.00,
		InterestRate:        0.10, // 10% interest rate
		InstallmentUnit:     "month",
		NumberOfInstallment: 12,
		StartDate:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		CustomerID:          "customer123",
	}

	// Mock repository calls
	mockRepo.On("CreateDisbursement", ctx, mock.AnythingOfType("*models.DisbursementDetail")).Return(nil)
	mockRepo.On("CreateLoanSummary", ctx, mock.AnythingOfType("*models.LoanSummary")).Return(nil)
	mockRepo.On("CreatePaymentSchedules", ctx, mock.AnythingOfType("[]*models.PaymentSchedule")).Return(nil)

	// Execute
	response, err := service.CreateDisbursement(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.InDelta(t, 91666.67, response.InstallmentAmount, 0.01) // (1000000 + 100000) / 12 rounded
	assert.Equal(t, "month", response.InstallmentUnit)

	mockRepo.AssertExpectations(t)
}

func TestGeneratePaymentSchedules_WeeklyInstallments(t *testing.T) {
	service := &disbursementService{}

	req := &models.DisbursementRequest{
		InstallmentUnit:     "week",
		NumberOfInstallment: 3,
	}

	startDate, _ := time.Parse("2006-01-02T15:04:05", "2025-01-01T00:00:00")
	installmentAmount := decimal.NewFromFloat(100000.00)
	totalAmount := decimal.NewFromFloat(300000.00)

	schedules := service.generatePaymentSchedules("loan123", req, installmentAmount, totalAmount, startDate)

	assert.Len(t, schedules, 3)

	// Check first installment
	assert.Equal(t, 1, schedules[0].InstallmentNumber)
	assert.Equal(t, startDate.AddDate(0, 0, 7), schedules[0].InstallmentDueDate)
	assert.Equal(t, 0.00, schedules[0].InstallmentPaid) // No payment made yet

	// Check second installment
	assert.Equal(t, 2, schedules[1].InstallmentNumber)
	assert.Equal(t, startDate.AddDate(0, 0, 14), schedules[1].InstallmentDueDate)
	assert.Equal(t, 0.00, schedules[1].InstallmentPaid) // No payment made yet

	// Check third installment
	assert.Equal(t, 3, schedules[2].InstallmentNumber)
	assert.Equal(t, startDate.AddDate(0, 0, 21), schedules[2].InstallmentDueDate)
	assert.Equal(t, 0.00, schedules[2].InstallmentPaid) // No payment made yet
}

func TestGeneratePaymentSchedules_MonthlyInstallments(t *testing.T) {
	service := &disbursementService{}

	req := &models.DisbursementRequest{
		InstallmentUnit:     "month",
		NumberOfInstallment: 2,
	}

	startDate, _ := time.Parse("2006-01-02T15:04:05", "2025-01-01T00:00:00")
	installmentAmount := decimal.NewFromFloat(150000.00)
	totalAmount := decimal.NewFromFloat(300000.00)

	schedules := service.generatePaymentSchedules("loan123", req, installmentAmount, totalAmount, startDate)

	assert.Len(t, schedules, 2)

	// Check first installment
	assert.Equal(t, 1, schedules[0].InstallmentNumber)
	assert.Equal(t, startDate.AddDate(0, 1, 0), schedules[0].InstallmentDueDate)
	assert.Equal(t, 0.00, schedules[0].InstallmentPaid)

	// Check second installment
	assert.Equal(t, 2, schedules[1].InstallmentNumber)
	assert.Equal(t, startDate.AddDate(0, 2, 0), schedules[1].InstallmentDueDate)
	assert.Equal(t, 0.00, schedules[1].InstallmentPaid)
}
