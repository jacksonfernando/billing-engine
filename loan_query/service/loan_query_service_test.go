package service

import (
	mocks "billing-engine/loan_query/_mock"
	"billing-engine/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoanQueryService_GetOutstandingBalance_Success(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 5280000.00,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		InstallmentUnit:   "week",
	}

	overdueSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
		{
			ID:                2,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	paidSchedules := []*models.PaymentSchedule{
		{
			ID:     3,
			Status: models.StatusPaid,
		},
		{
			ID:     4,
			Status: models.StatusPaid,
		},
	}

	pendingSchedules := []*models.PaymentSchedule{
		{
			ID:     5,
			Status: models.StatusPending,
		},
		{
			ID:     6,
			Status: models.StatusPending,
		},
	}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)
	mockRepo.On("GetPaidPaymentSchedulesByLoanID", ctx, "loan_123").Return(paidSchedules, nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(pendingSchedules, nil)

	// Execute
	response, err := service.GetOutstandingBalance(ctx, "loan_123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, "customer_123", response.CustomerID)
	assert.Equal(t, 5280000.00, response.OutstandingAmount)
	assert.Equal(t, 220000.00, response.OverdueAmount) // 2 * 110000

	assert.Equal(t, 2, response.OverdueInstallments)
	assert.Equal(t, 2, response.PaidInstallments)
	assert.Equal(t, 2, response.RemainingInstallments)
	assert.Equal(t, "week", response.LoanDetails.InstallmentUnit)
	assert.Equal(t, 110000.00, response.LoanDetails.InstallmentAmount)
	assert.Equal(t, 50, response.LoanDetails.TotalInstallments)

	mockRepo.AssertExpectations(t)
}

func TestLoanQueryService_GetOutstandingBalance_LoanNotFound(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(nil, nil)

	// Execute
	response, err := service.GetOutstandingBalance(ctx, "loan_123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "loan not found")

	mockRepo.AssertExpectations(t)
}

func TestLoanQueryService_GetDelinquencyStatus_Success(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 5280000.00,
		InstallmentAmount: 110000.00,
		InstallmentUnit:   "week",
	}

	overdueSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
		{
			ID:                2,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)

	// Execute
	response, err := service.GetDelinquencyStatus(ctx, "loan_123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, "customer_123", response.CustomerID)
	assert.True(t, response.IsDelinquent)

	assert.Equal(t, "week", response.InstallmentUnit)
	assert.Equal(t, 2, response.OverdueInstallments)
	assert.Equal(t, 220000.00, response.OverdueAmount)
	assert.Equal(t, 5280000.00, response.OutstandingAmount)
	assert.Equal(t, 220000.00, response.RequiredPaymentAmount)

	mockRepo.AssertExpectations(t)
}

func TestLoanQueryService_GetDelinquencyStatus_NotDelinquent(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 5280000.00,
		InstallmentAmount: 110000.00,
		InstallmentUnit:   "week",
	}

	overdueSchedules := []*models.PaymentSchedule{} // No overdue schedules

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)

	// Execute
	response, err := service.GetDelinquencyStatus(ctx, "loan_123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, "customer_123", response.CustomerID)
	assert.False(t, response.IsDelinquent)

	assert.Equal(t, 0, response.OverdueInstallments)
	assert.Equal(t, 0.00, response.OverdueAmount)
	assert.Equal(t, 0.00, response.RequiredPaymentAmount)

	mockRepo.AssertExpectations(t)
}

func TestLoanQueryService_GetLoanSchedule_Success(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		PrincipalAmount:   5000000.00,
		OutstandingAmount: 5280000.00,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		InstallmentUnit:   "week",
	}

	paidDate := time.Now().AddDate(0, 0, -7)
	paymentSchedules := []*models.PaymentSchedule{
		{
			ID:                 1,
			InstallmentNumber:  1,
			InstallmentDueDate: time.Now().AddDate(0, 0, -14),
			InstallmentAmount:  110000.00,
			InstallmentPaid:    5390000.00,
			Status:             models.StatusPaid,
			UpdatedAt:          paidDate,
		},
		{
			ID:                 2,
			InstallmentNumber:  2,
			InstallmentDueDate: time.Now().AddDate(0, 0, -7),
			InstallmentAmount:  110000.00,
			InstallmentPaid:    5280000.00,
			Status:             models.StatusPending,
		},
	}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetPaymentSchedulesByLoanID", ctx, "loan_123").Return(paymentSchedules, nil)

	// Execute
	response, err := service.GetLoanSchedule(ctx, "loan_123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, "week", response.LoanSummary.InstallmentUnit)
	assert.Equal(t, 50, response.LoanSummary.TotalInstallments)
	assert.Equal(t, 110000.00, response.LoanSummary.InstallmentAmount)
	assert.Equal(t, 5000000.00, response.LoanSummary.DisbursedAmount)
	assert.Equal(t, 5280000.00, response.LoanSummary.OutstandingAmount)

	assert.Len(t, response.Schedule, 2)

	// Check first schedule (paid)
	assert.Equal(t, 1, response.Schedule[0].InstallmentNumber)
	assert.Equal(t, 110000.00, response.Schedule[0].InstallmentAmount)
	assert.Equal(t, 5390000.00, response.Schedule[0].InstallmentPaid)
	assert.Equal(t, models.StatusPaid, response.Schedule[0].Status)
	assert.NotNil(t, response.Schedule[0].PaidDate)

	// Check second schedule (pending)
	assert.Equal(t, 2, response.Schedule[1].InstallmentNumber)
	assert.Equal(t, 110000.00, response.Schedule[1].InstallmentAmount)
	assert.Equal(t, 5280000.00, response.Schedule[1].InstallmentPaid)
	assert.Equal(t, models.StatusPending, response.Schedule[1].Status)
	assert.Nil(t, response.Schedule[1].PaidDate)

	mockRepo.AssertExpectations(t)
}

func TestLoanQueryService_GetLoanSchedule_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewLoanQueryMySQLRepositoryInterface(t)
	service := NewLoanQueryService(mockRepo)
	ctx := context.Background()

	// Mock repository error
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(nil, errors.New("database error"))

	// Execute
	response, err := service.GetLoanSchedule(ctx, "loan_123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get loan summary")

	mockRepo.AssertExpectations(t)
}
