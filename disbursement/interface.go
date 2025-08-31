package disbursement

import (
	"billing-engine/models"
	"context"
)

// DisbursementMySQLRepositoryInterface defines the interface for disbursement repository
type DisbursementMySQLRepositoryInterface interface {
	CreateDisbursement(ctx context.Context, disbursement *models.DisbursementDetail) error
	CreateLoanSummary(ctx context.Context, loanSummary *models.LoanSummary) error
	CreatePaymentSchedules(ctx context.Context, paymentSchedules []*models.PaymentSchedule) error
	GetDisbursementByLoanID(ctx context.Context, loanID string) (*models.DisbursementDetail, error)
	GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error)
}

// DisbursementServiceInterface defines the interface for disbursement service
type DisbursementServiceInterface interface {
	CreateDisbursement(ctx context.Context, req *models.DisbursementRequest) (*models.DisbursementResponse, error)
}
