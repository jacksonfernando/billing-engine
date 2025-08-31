package loan_query

import (
	"billing-engine/models"
	"context"
)

// LoanQueryMySQLRepositoryInterface defines the interface for loan query repository
type LoanQueryMySQLRepositoryInterface interface {
	GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error)
	GetPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
	GetOverduePaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
	GetPaidPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
	GetPendingPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
}

// LoanQueryServiceInterface defines the interface for loan query service
type LoanQueryServiceInterface interface {
	GetOutstandingBalance(ctx context.Context, loanID string) (*models.OutstandingBalanceResponse, error)
	GetDelinquencyStatus(ctx context.Context, loanID string) (*models.DelinquencyResponse, error)
	GetLoanSchedule(ctx context.Context, loanID string) (*models.LoanScheduleResponse, error)
}
