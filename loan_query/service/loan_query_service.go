package service

import (
	"context"
	"fmt"
	"time"

	"billing-engine/loan_query"
	"billing-engine/models"
)

type loanQueryService struct {
	loanQueryRepo loan_query.LoanQueryMySQLRepositoryInterface
}

// NewLoanQueryService creates a new loan query service instance
func NewLoanQueryService(loanQueryRepo loan_query.LoanQueryMySQLRepositoryInterface) loan_query.LoanQueryServiceInterface {
	return &loanQueryService{
		loanQueryRepo: loanQueryRepo,
	}
}

func (s *loanQueryService) GetOutstandingBalance(ctx context.Context, loanID string) (*models.OutstandingBalanceResponse, error) {
	// Get loan summary
	loanSummary, err := s.loanQueryRepo.GetLoanSummaryByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan summary: %v", err)
	}
	if loanSummary == nil {
		return nil, fmt.Errorf("loan not found")
	}

	// Get overdue schedules
	overdueSchedules, err := s.loanQueryRepo.GetOverduePaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue schedules: %v", err)
	}

	// Get paid schedules
	paidSchedules, err := s.loanQueryRepo.GetPaidPaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paid schedules: %v", err)
	}

	// Get pending schedules
	pendingSchedules, err := s.loanQueryRepo.GetPendingPaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending schedules: %v", err)
	}

	// Calculate overdue amount
	var overdueAmount float64
	for _, schedule := range overdueSchedules {
		overdueAmount += schedule.InstallmentAmount
	}

	return &models.OutstandingBalanceResponse{
		LoanID:     loanID,
		CustomerID: loanSummary.CustomerID,
		LoanDetails: models.LoanDetailsResponse{
			InstallmentUnit:   loanSummary.InstallmentUnit,
			InstallmentAmount: loanSummary.InstallmentAmount,
			TotalInstallments: loanSummary.NoOfInstallment,
		},
		OutstandingAmount: loanSummary.OutstandingAmount,
		OverdueAmount:     overdueAmount,

		OverdueInstallments:   len(overdueSchedules),
		PaidInstallments:      len(paidSchedules),
		RemainingInstallments: len(pendingSchedules),
	}, nil
}

func (s *loanQueryService) GetDelinquencyStatus(ctx context.Context, loanID string) (*models.DelinquencyResponse, error) {
	// Get loan summary
	loanSummary, err := s.loanQueryRepo.GetLoanSummaryByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan summary: %v", err)
	}
	if loanSummary == nil {
		return nil, fmt.Errorf("loan not found")
	}

	// Get overdue schedules
	overdueSchedules, err := s.loanQueryRepo.GetOverduePaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue schedules: %v", err)
	}

	// Calculate overdue amount and required payment
	var overdueAmount float64
	for _, schedule := range overdueSchedules {
		overdueAmount += schedule.InstallmentAmount
	}

	// Determine if delinquent (has overdue installments)
	isDelinquent := len(overdueSchedules) > 1

	return &models.DelinquencyResponse{
		LoanID:       loanID,
		CustomerID:   loanSummary.CustomerID,
		IsDelinquent: isDelinquent,

		InstallmentUnit:       loanSummary.InstallmentUnit,
		OverdueInstallments:   len(overdueSchedules),
		OverdueAmount:         overdueAmount,
		OutstandingAmount:     loanSummary.OutstandingAmount,
		RequiredPaymentAmount: overdueAmount, // Must pay all overdue amounts
	}, nil
}

func (s *loanQueryService) GetLoanSchedule(ctx context.Context, loanID string) (*models.LoanScheduleResponse, error) {
	// Get loan summary
	loanSummary, err := s.loanQueryRepo.GetLoanSummaryByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan summary: %v", err)
	}
	if loanSummary == nil {
		return nil, fmt.Errorf("loan not found")
	}

	// Get all payment schedules
	paymentSchedules, err := s.loanQueryRepo.GetPaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment schedules: %v", err)
	}

	// Convert to response format
	scheduleResponses := make([]models.PaymentScheduleResponse, 0, len(paymentSchedules))
	for _, schedule := range paymentSchedules {
		var paidDate *time.Time
		if schedule.Status == models.StatusPaid {
			paidDate = &schedule.UpdatedAt
		}

		scheduleResponse := models.PaymentScheduleResponse{
			InstallmentNumber: schedule.InstallmentNumber,
			DueDate:           schedule.InstallmentDueDate,
			InstallmentAmount: schedule.InstallmentAmount,
			InstallmentPaid:   schedule.InstallmentPaid,
			Status:            schedule.Status,
			PaidDate:          paidDate,
		}
		scheduleResponses = append(scheduleResponses, scheduleResponse)
	}

	return &models.LoanScheduleResponse{
		LoanID: loanID,
		LoanSummary: models.LoanSummaryScheduleResponse{
			InstallmentUnit:   loanSummary.InstallmentUnit,
			TotalInstallments: loanSummary.NoOfInstallment,
			InstallmentAmount: loanSummary.InstallmentAmount,
			DisbursedAmount:   loanSummary.PrincipalAmount,
			OutstandingAmount: loanSummary.OutstandingAmount,
		},
		Schedule: scheduleResponses,
	}, nil
}
