package mysql

import (
	"context"
	"errors"
	"time"

	"billing-engine/models"
	"billing-engine/repayment"

	"gorm.io/gorm"
)

type repaymentMySQLRepository struct {
	db *gorm.DB
}

// NewRepaymentMySQLRepository creates a new repayment repository instance
func NewRepaymentMySQLRepository(db *gorm.DB) repayment.RepaymentMySQLRepositoryInterface {
	return &repaymentMySQLRepository{db: db}
}

func (r *repaymentMySQLRepository) GetLoanSummaryByLoanID(ctx context.Context, loanID string) (*models.LoanSummary, error) {
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

func (r *repaymentMySQLRepository) GetPendingPaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
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

func (r *repaymentMySQLRepository) GetOverduePaymentSchedulesByLoanID(ctx context.Context, loanID string) ([]*models.PaymentSchedule, error) {
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

func (r *repaymentMySQLRepository) UpdatePaymentSchedules(ctx context.Context, schedules []*models.PaymentSchedule) error {
	for _, schedule := range schedules {
		if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *repaymentMySQLRepository) UpdateLoanSummary(ctx context.Context, loanSummary *models.LoanSummary) error {
	return r.db.WithContext(ctx).Save(loanSummary).Error
}

func (r *repaymentMySQLRepository) CreatePaymentHistory(ctx context.Context, histories []*models.PaymentScheduleHistory) error {
	return r.db.WithContext(ctx).Create(&histories).Error
}

func (r *repaymentMySQLRepository) GetNextDueDate(ctx context.Context, loanID string) (*time.Time, error) {
	var schedule models.PaymentSchedule
	err := r.db.WithContext(ctx).
		Where("loan_id = ? AND status = ? AND deleted_at IS NULL", loanID, models.StatusPending).
		Order("installment_number ASC").
		First(&schedule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &schedule.InstallmentDueDate, nil
}
