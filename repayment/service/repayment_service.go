package service

import (
	"context"
	"fmt"
	"time"

	"billing-engine/models"
	"billing-engine/repayment"
)

type repaymentService struct {
	repaymentRepo repayment.RepaymentMySQLRepositoryInterface
}

// NewRepaymentService creates a new repayment service instance
func NewRepaymentService(repaymentRepo repayment.RepaymentMySQLRepositoryInterface) repayment.RepaymentServiceInterface {
	return &repaymentService{
		repaymentRepo: repaymentRepo,
	}
}

func (s *repaymentService) ProcessRepayment(ctx context.Context, req *models.RepaymentRequest) (*models.RepaymentResponse, error) {
	// 1. Validate loan exists
	loanSummary, err := s.validateLoanExists(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	// 2. Get payment schedules
	overdueSchedules, pendingSchedules, err := s.getPaymentSchedules(ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	// 3. Calculate payment plan
	schedulesToPay, requiredAmount, err := s.calculatePaymentPlan(overdueSchedules, pendingSchedules, req.PaymentAmount)
	if err != nil {
		return nil, err
	}

	// 4. Validate payment amount
	if err := s.validatePaymentAmount(req.PaymentAmount, requiredAmount); err != nil {
		return nil, err
	}

	// 5. Process payment
	paymentDate := time.Now()
	if err := s.processPaymentSchedules(ctx, schedulesToPay, paymentDate); err != nil {
		return nil, err
	}

	// 6. Update loan summary
	remainingSchedules, err := s.updateLoanSummary(ctx, loanSummary, req.PaymentAmount, paymentDate)
	if err != nil {
		return nil, err
	}

	// 7. Build response
	return s.buildRepaymentResponse(ctx, req, loanSummary, schedulesToPay, remainingSchedules, paymentDate)
}

// validateLoanExists checks if the loan exists and returns the loan summary
func (s *repaymentService) validateLoanExists(ctx context.Context, loanID string) (*models.LoanSummary, error) {
	loanSummary, err := s.repaymentRepo.GetLoanSummaryByLoanID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan summary: %v", err)
	}
	if loanSummary == nil {
		return nil, fmt.Errorf("loan not found")
	}
	return loanSummary, nil
}

// getPaymentSchedules retrieves overdue and pending payment schedules
func (s *repaymentService) getPaymentSchedules(ctx context.Context, loanID string) ([]*models.PaymentSchedule, []*models.PaymentSchedule, error) {
	// Get overdue payment schedules first (must be paid first)
	overdueSchedules, err := s.repaymentRepo.GetOverduePaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get overdue schedules: %v", err)
	}

	// Get all pending payment schedules
	pendingSchedules, err := s.repaymentRepo.GetPendingPaymentSchedulesByLoanID(ctx, loanID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get pending schedules: %v", err)
	}

	if len(pendingSchedules) == 0 {
		return nil, nil, fmt.Errorf("no pending installments found")
	}

	return overdueSchedules, pendingSchedules, nil
}

// calculatePaymentPlan determines which schedules to pay and the required amount
func (s *repaymentService) calculatePaymentPlan(overdueSchedules, pendingSchedules []*models.PaymentSchedule, paymentAmount float64) ([]*models.PaymentSchedule, float64, error) {
	var requiredAmount float64
	var schedulesToPay []*models.PaymentSchedule

	// If there are overdue installments, customer MUST pay ALL overdue installments at once
	if len(overdueSchedules) > 0 {
		for _, schedule := range overdueSchedules {
			requiredAmount += schedule.InstallmentAmount
			schedulesToPay = append(schedulesToPay, schedule)
		}
		// Customer must pay the exact total amount of all overdue installments
		return schedulesToPay, requiredAmount, nil
	}

	// If no overdue installments, customer can pay exactly one pending installment
	if len(pendingSchedules) > 0 {
		nextInstallment := pendingSchedules[0]
		requiredAmount = nextInstallment.InstallmentAmount
		schedulesToPay = append(schedulesToPay, nextInstallment)
		return schedulesToPay, requiredAmount, nil
	}

	// No installments to pay
	return nil, 0, fmt.Errorf("no pending installments found")
}

// validatePaymentAmount ensures the payment amount matches the required amount exactly
func (s *repaymentService) validatePaymentAmount(paymentAmount, requiredAmount float64) error {
	if paymentAmount != requiredAmount {
		return fmt.Errorf("payment amount %.2f does not match required amount %.2f. You must pay the exact amount for all overdue installments or the next pending installment", paymentAmount, requiredAmount)
	}
	return nil
}

// processPaymentSchedules updates payment schedules and creates history records
func (s *repaymentService) processPaymentSchedules(ctx context.Context, schedulesToPay []*models.PaymentSchedule, paymentDate time.Time) error {
	var histories []*models.PaymentScheduleHistory

	// Prepare schedules and history records
	for _, schedule := range schedulesToPay {
		// Create history record before updating
		history := s.createPaymentHistory(schedule)
		histories = append(histories, history)

		// Update schedule
		s.updateScheduleForPayment(schedule, paymentDate)
	}

	// Update payment schedules in database
	if err := s.repaymentRepo.UpdatePaymentSchedules(ctx, schedulesToPay); err != nil {
		return fmt.Errorf("failed to update payment schedules: %v", err)
	}

	// Create payment history
	if err := s.repaymentRepo.CreatePaymentHistory(ctx, histories); err != nil {
		return fmt.Errorf("failed to create payment history: %v", err)
	}

	return nil
}

// createPaymentHistory creates a payment history record for a schedule
func (s *repaymentService) createPaymentHistory(schedule *models.PaymentSchedule) *models.PaymentScheduleHistory {
	return &models.PaymentScheduleHistory{
		ScheduleID:         schedule.ID,
		LoanID:             schedule.LoanID,
		Action:             models.ActionPayment,
		InstallmentNumber:  schedule.InstallmentNumber,
		InstallmentAmount:  schedule.InstallmentAmount,
		InstallmentDueDate: schedule.InstallmentDueDate,
		Status:             models.StatusPaid,
		Currency:           schedule.Currency,
		CreatedBy:          "system",
	}
}

// updateScheduleForPayment updates a schedule to mark it as paid
func (s *repaymentService) updateScheduleForPayment(schedule *models.PaymentSchedule, paymentDate time.Time) {
	schedule.Status = models.StatusPaid
	schedule.InstallmentPaid = schedule.InstallmentAmount
	schedule.UpdatedBy = "system"
	schedule.UpdatedAt = paymentDate
}

// updateLoanSummary updates the loan summary and returns remaining schedules
func (s *repaymentService) updateLoanSummary(ctx context.Context, loanSummary *models.LoanSummary, paymentAmount float64, paymentDate time.Time) ([]*models.PaymentSchedule, error) {
	// Update loan summary
	loanSummary.OutstandingAmount -= paymentAmount
	loanSummary.UpdatedBy = "system"
	loanSummary.UpdatedAt = paymentDate

	// Check if all installments are paid
	remainingPendingSchedules, err := s.repaymentRepo.GetPendingPaymentSchedulesByLoanID(ctx, loanSummary.LoanID)
	if err != nil {
		return nil, fmt.Errorf("failed to check remaining schedules: %v", err)
	}

	if len(remainingPendingSchedules) == 0 {
		loanSummary.Status = models.StatusPaid
	}

	if err := s.repaymentRepo.UpdateLoanSummary(ctx, loanSummary); err != nil {
		return nil, fmt.Errorf("failed to update loan summary: %v", err)
	}

	return remainingPendingSchedules, nil
}

// buildRepaymentResponse constructs the final response
func (s *repaymentService) buildRepaymentResponse(ctx context.Context, req *models.RepaymentRequest, loanSummary *models.LoanSummary, schedulesToPay []*models.PaymentSchedule, remainingSchedules []*models.PaymentSchedule, paymentDate time.Time) (*models.RepaymentResponse, error) {
	// Get next due date
	nextDueDate, err := s.repaymentRepo.GetNextDueDate(ctx, req.LoanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get next due date: %v", err)
	}

	var nextDue time.Time
	if nextDueDate != nil {
		nextDue = *nextDueDate
	}

	return &models.RepaymentResponse{
		LoanID:                req.LoanID,
		PaymentAmount:         req.PaymentAmount,
		InstallmentsPaid:      len(schedulesToPay),
		InstallmentAmount:     loanSummary.InstallmentAmount,
		RemainingInstallments: len(remainingSchedules),
		OutstandingAmount:     loanSummary.OutstandingAmount,
		NextDueDate:           nextDue,
		PaymentDate:           paymentDate,
	}, nil
}
