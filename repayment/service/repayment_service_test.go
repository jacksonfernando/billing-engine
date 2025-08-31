package service

import (
	"billing-engine/models"
	mocks "billing-engine/repayment/_mock"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRepaymentService_ProcessRepayment_Success(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 220000.00,
	}

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 5500000.00,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		Status:            models.StatusPending,
	}

	overdueSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			LoanID:            "loan_123",
			InstallmentNumber: 1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
		{
			ID:                2,
			LoanID:            "loan_123",
			InstallmentNumber: 2,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	pendingSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			LoanID:            "loan_123",
			InstallmentNumber: 1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
		{
			ID:                2,
			LoanID:            "loan_123",
			InstallmentNumber: 2,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	remainingSchedules := []*models.PaymentSchedule{
		{
			ID:                3,
			LoanID:            "loan_123",
			InstallmentNumber: 3,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	nextDueDate := time.Now().AddDate(0, 0, 7)

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(pendingSchedules, nil).Once()
	mockRepo.On("UpdatePaymentSchedules", ctx, mock.AnythingOfType("[]*models.PaymentSchedule")).Return(nil)
	mockRepo.On("CreatePaymentHistory", ctx, mock.AnythingOfType("[]*models.PaymentScheduleHistory")).Return(nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(remainingSchedules, nil).Once()
	mockRepo.On("UpdateLoanSummary", ctx, mock.AnythingOfType("*models.LoanSummary")).Return(nil)
	mockRepo.On("GetNextDueDate", ctx, "loan_123").Return(&nextDueDate, nil)

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, 220000.00, response.PaymentAmount)
	assert.Equal(t, 2, response.InstallmentsPaid)
	assert.Equal(t, 110000.00, response.InstallmentAmount)
	assert.Equal(t, 1, response.RemainingInstallments)
	assert.Equal(t, nextDueDate, response.NextDueDate)

	mockRepo.AssertExpectations(t)
}

func TestRepaymentService_ProcessRepayment_LoanNotFound(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 220000.00,
	}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(nil, nil)

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "loan not found")

	mockRepo.AssertExpectations(t)
}

func TestRepaymentService_ProcessRepayment_IncorrectPaymentAmount(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 100000.00, // Incorrect amount
	}

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 5500000.00,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		Status:            models.StatusPending,
	}

	overdueSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			LoanID:            "loan_123",
			InstallmentNumber: 1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	pendingSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			LoanID:            "loan_123",
			InstallmentNumber: 1,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(pendingSchedules, nil)

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "payment amount 100000.00 does not match required amount 110000.00")

	mockRepo.AssertExpectations(t)
}

func TestRepaymentService_ProcessRepayment_NoPendingInstallments(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 220000.00,
	}

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 0,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		Status:            models.StatusPaid,
	}

	overdueSchedules := []*models.PaymentSchedule{}
	pendingSchedules := []*models.PaymentSchedule{}

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(pendingSchedules, nil)

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "no pending installments found")

	mockRepo.AssertExpectations(t)
}

func TestRepaymentService_ProcessRepayment_AllInstallmentsPaid(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 110000.00,
	}

	loanSummary := &models.LoanSummary{
		LoanID:            "loan_123",
		CustomerID:        "customer_123",
		OutstandingAmount: 110000.00,
		InstallmentAmount: 110000.00,
		NoOfInstallment:   50,
		Status:            models.StatusPending,
	}

	overdueSchedules := []*models.PaymentSchedule{}
	pendingSchedules := []*models.PaymentSchedule{
		{
			ID:                1,
			LoanID:            "loan_123",
			InstallmentNumber: 50,
			InstallmentAmount: 110000.00,
			Status:            models.StatusPending,
		},
	}

	remainingSchedules := []*models.PaymentSchedule{} // No remaining schedules

	// Mock repository calls
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(loanSummary, nil)
	mockRepo.On("GetOverduePaymentSchedulesByLoanID", ctx, "loan_123").Return(overdueSchedules, nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(pendingSchedules, nil).Once()
	mockRepo.On("UpdatePaymentSchedules", ctx, mock.AnythingOfType("[]*models.PaymentSchedule")).Return(nil)
	mockRepo.On("CreatePaymentHistory", ctx, mock.AnythingOfType("[]*models.PaymentScheduleHistory")).Return(nil)
	mockRepo.On("GetPendingPaymentSchedulesByLoanID", ctx, "loan_123").Return(remainingSchedules, nil).Once()
	mockRepo.On("UpdateLoanSummary", ctx, mock.MatchedBy(func(ls *models.LoanSummary) bool {
		return ls.Status == models.StatusPaid // Should be marked as PAID
	})).Return(nil)
	mockRepo.On("GetNextDueDate", ctx, "loan_123").Return(nil, nil)

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "loan_123", response.LoanID)
	assert.Equal(t, 110000.00, response.PaymentAmount)
	assert.Equal(t, 1, response.InstallmentsPaid)
	assert.Equal(t, 0, response.RemainingInstallments)

	mockRepo.AssertExpectations(t)
}

func TestRepaymentService_ProcessRepayment_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewRepaymentMySQLRepositoryInterface(t)
	service := NewRepaymentService(mockRepo)
	ctx := context.Background()

	req := &models.RepaymentRequest{
		LoanID:        "loan_123",
		PaymentAmount: 220000.00,
	}

	// Mock repository error
	mockRepo.On("GetLoanSummaryByLoanID", ctx, "loan_123").Return(nil, errors.New("database error"))

	// Execute
	response, err := service.ProcessRepayment(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get loan summary")

	mockRepo.AssertExpectations(t)
}
