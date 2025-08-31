package repayment

import (
	"billing-engine/models"
	"context"
	"time"
)

// RepaymentMySQLRepositoryInterface defines the interface for repayment repository
type RepaymentMySQLRepositoryInterface interface {
	GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error)
	GetPendingPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
	GetOverduePaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error)
	UpdatePaymentSchedules(ctx context.Context, schedules []*models.PaymentSchedule) error
	UpdateLoanSummary(ctx context.Context, loanSummary *models.LoanSummary) error
	CreatePaymentHistory(ctx context.Context, histories []*models.PaymentScheduleHistory) error
	GetNextDueDate(ctx context.Context, loanID string) (*time.Time, error)
}

// RepaymentServiceInterface defines the interface for repayment service
type RepaymentServiceInterface interface {
	ProcessRepayment(ctx context.Context, req *models.RepaymentRequest) (*models.RepaymentResponse, error)
}
