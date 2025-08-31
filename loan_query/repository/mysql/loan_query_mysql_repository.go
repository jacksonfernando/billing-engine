package mysql

import (
	"context"
	"errors"
	"time"

	"billing-engine/loan_query"
	"billing-engine/models"

	"gorm.io/gorm"
)

type loanQueryMySQLRepository struct {
	db *gorm.DB
}

// NewLoanQueryMySQLRepository creates a new loan query repository instance
func NewLoanQueryMySQLRepository(db *gorm.DB) loan_query.LoanQueryMySQLRepositoryInterface {
	return &loanQueryMySQLRepository{db: db}
}

func (r *loanQueryMySQLRepository) GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error) {
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

func (r *loanQueryMySQLRepository) GetPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
	var schedules []*models.PaymentSchedule
	err := r.db.WithContext(ctx).
		Where("loan_id = ? AND deleted_at IS NULL", loanID).
		Order("installment_number ASC").
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *loanQueryMySQLRepository) GetOverduePaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
	var schedules []*models.PaymentSchedule
	now := time.Now()
	err := r.db.WithContext(ctx).
		Where("loan_id = ? AND status = ? AND installment_due_date < ? AND deleted_at IS NULL", loanID, models.StatusPending, now).
		Order("installment_number ASC").
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *loanQueryMySQLRepository) GetPaidPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
	var schedules []*models.PaymentSchedule
	err := r.db.WithContext(ctx).
		Where("loan_id = ? AND status = ? AND deleted_at IS NULL", loanID, models.StatusPaid).
		Order("installment_number ASC").
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *loanQueryMySQLRepository) GetPendingPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
	var schedules []*models.PaymentSchedule
	err := r.db.WithContext(ctx).
		Where("loan_id = ? AND status = ? AND deleted_at IS NULL", loanID, models.StatusPending).
		Order("installment_number ASC").
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}
