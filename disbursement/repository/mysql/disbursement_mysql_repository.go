package mysql

import (
	"context"
	"errors"

	"billing-engine/disbursement"
	"billing-engine/models"

	"gorm.io/gorm"
)

type disbursementMySQLRepository struct {
	db *gorm.DB
}

// NewDisbursementMySQLRepository creates a new disbursement repository instance
func NewDisbursementMySQLRepository(db *gorm.DB) disbursement.DisbursementMySQLRepositoryInterface {
	return &disbursementMySQLRepository{db: db}
}

func (r *disbursementMySQLRepository) CreateDisbursement(ctx context.Context, disbursementDetail *models.DisbursementDetail) error {
	return r.db.WithContext(ctx).Create(disbursementDetail).Error
}

func (r *disbursementMySQLRepository) CreateLoanSummary(ctx context.Context, loanSummary *models.LoanSummary) error {
	return r.db.WithContext(ctx).Create(loanSummary).Error
}

func (r *disbursementMySQLRepository) CreatePaymentSchedules(ctx context.Context, paymentSchedules []*models.PaymentSchedule) error {
	return r.db.WithContext(ctx).Create(&paymentSchedules).Error
}

func (r *disbursementMySQLRepository) GetDisbursementByLoanID(ctx context.Context, loanID string) (*models.DisbursementDetail, error) {
	var disbursementDetail models.DisbursementDetail
	err := r.db.WithContext(ctx).Where("loan_id = ? AND deleted_at IS NULL", loanID).First(&disbursementDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &disbursementDetail, nil
}

func (r *disbursementMySQLRepository) GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error) {
	var loanSummary models.LoanSummary
	err := r.db.WithContext(ctx).Where("loan_id = ? AND deleted_at IS NULL", loanID).First(&loanSummary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &loanSummary, nil
}
