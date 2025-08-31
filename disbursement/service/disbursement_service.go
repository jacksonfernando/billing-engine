package service

import (
	"context"
	"fmt"
	"time"

	"billing-engine/disbursement"
	"billing-engine/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type disbursementService struct {
	disbursementRepo disbursement.DisbursementMySQLRepositoryInterface
}

// NewDisbursementService creates a new disbursement service instance
func NewDisbursementService(disbursementRepo disbursement.DisbursementMySQLRepositoryInterface) disbursement.DisbursementServiceInterface {
	return &disbursementService{
		disbursementRepo: disbursementRepo,
	}
}

func (s *disbursementService) CreateDisbursement(ctx context.Context, req *models.DisbursementRequest) (*models.DisbursementResponse, error) {
	loanID := fmt.Sprintf("loan_%s", uuid.New().String())
	startDate := req.StartDate
	if startDate.IsZero() {
		return nil, fmt.Errorf("start date cannot be zero")
	}
	principal := decimal.NewFromFloat(req.PrincipalAmount)
	interestRate := decimal.NewFromFloat(req.InterestRate)
	numberOfInstallments := decimal.NewFromInt(int64(req.NumberOfInstallment))
	// Calculate interest amount: principal * interest_rate (total interest for the loan)
	interestAmount := principal.Mul(interestRate)
	// Calculate total amount: principal + interest_amount
	totalAmount := principal.Add(interestAmount)
	// Calculate installment amount: total_amount / number_of_installments
	installmentAmount := totalAmount.Div(numberOfInstallments)
	// The effective interest rate is the same as the input interest rate
	effectiveInterestRate := interestRate
	// Create disbursement detail
	disbursementDetail := &models.DisbursementDetail{
		LoanID:            loanID,
		CustomerID:        req.CustomerID,
		DisbursementDate:  startDate,
		DisbursedAmount:   principal.InexactFloat64(),
		DisbursedCurrency: models.CurrencyIDR,
		Status:            models.StatusPending,
		CreatedBy:         "system",
		UpdatedBy:         "system",
	}

	// Create loan summary
	loanSummary := &models.LoanSummary{
		LoanID:                loanID,
		CustomerID:            req.CustomerID,
		PrincipalAmount:       principal.InexactFloat64(),
		InterestAmount:        interestAmount.InexactFloat64(),
		OutstandingAmount:     totalAmount.InexactFloat64(),
		NoOfInstallment:       req.NumberOfInstallment,
		InstallmentUnit:       req.InstallmentUnit,
		InstallmentAmount:     installmentAmount.InexactFloat64(),
		EffectiveInterestRate: effectiveInterestRate.InexactFloat64(),

		Status:        models.StatusPending,
		LoanStartDate: startDate,
		CreatedBy:     "system",
		UpdatedBy:     "system",
	}

	// Generate payment schedules
	paymentSchedules := s.generatePaymentSchedules(loanID, req, installmentAmount, totalAmount, startDate)
	if err := s.disbursementRepo.CreateDisbursement(ctx, disbursementDetail); err != nil {
		return nil, fmt.Errorf("failed to create disbursement: %v", err)
	}
	if err := s.disbursementRepo.CreateLoanSummary(ctx, loanSummary); err != nil {
		return nil, fmt.Errorf("failed to create loan summary: %v", err)
	}
	if err := s.disbursementRepo.CreatePaymentSchedules(ctx, paymentSchedules); err != nil {
		return nil, fmt.Errorf("failed to create payment schedules: %v", err)
	}

	firstDueDate := paymentSchedules[0].InstallmentDueDate
	finalDueDate := paymentSchedules[len(paymentSchedules)-1].InstallmentDueDate

	return &models.DisbursementResponse{
		LoanID:              loanID,
		CustomerID:          req.CustomerID,
		DisbursedAmount:     principal.InexactFloat64(),
		InstallmentAmount:   installmentAmount.InexactFloat64(),
		OutstandingAmount:   totalAmount.InexactFloat64(),
		InstallmentUnit:     req.InstallmentUnit,
		NumberOfInstallment: req.NumberOfInstallment,
		DisbursementDate:    startDate,
		FirstDueDate:        firstDueDate,
		FinalDueDate:        finalDueDate,
	}, nil
}

func (s *disbursementService) generatePaymentSchedules(loanID string, req *models.DisbursementRequest, installmentAmount, totalAmount decimal.Decimal, startDate time.Time) []*models.PaymentSchedule {
	schedules := make([]*models.PaymentSchedule, 0, req.NumberOfInstallment)
	for i := 1; i <= req.NumberOfInstallment; i++ {
		var dueDate time.Time
		if req.InstallmentUnit == models.InstallmentUnitWeek {
			dueDate = startDate.AddDate(0, 0, 7*i)
		} else { // month
			dueDate = startDate.AddDate(0, i, 0)
		}

		schedule := &models.PaymentSchedule{
			LoanID:             loanID,
			InstallmentNumber:  i,
			InstallmentAmount:  installmentAmount.InexactFloat64(),
			InstallmentDueDate: dueDate,
			InstallmentPaid:    0,
			Status:             models.StatusPending,
			Currency:           models.CurrencyIDR,
			CreatedBy:          "system",
			UpdatedBy:          "system",
		}

		schedules = append(schedules, schedule)
	}

	return schedules
}
