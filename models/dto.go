package models

import "time"

// Request DTOs
type DisbursementRequest struct {
	PrincipalAmount     float64   `json:"principal_amount" validate:"gt=0"`
	InterestRate        float64   `json:"interest_rate" validate:"gt=0,lte=1"`
	InstallmentUnit     string    `json:"installment_unit" validate:"required,oneof=week month"`
	NumberOfInstallment int       `json:"number_of_installment" validate:"gt=0"`
	StartDate           time.Time `json:"start_date" validate:"required"`
	CustomerID          string    `json:"customer_id" validate:"required"`
}

type RepaymentRequest struct {
	LoanID        string  `json:"loan_id" validate:"required"`
	PaymentAmount float64 `json:"payment_amount" validate:"gt=0"`
}

// Response DTOs
type DisbursementResponse struct {
	LoanID              string    `json:"loan_id"`
	CustomerID          string    `json:"customer_id"`
	DisbursedAmount     float64   `json:"disbursed_amount"`
	InstallmentAmount   float64   `json:"installment_amount"`
	OutstandingAmount   float64   `json:"outstanding_amount"`
	InstallmentUnit     string    `json:"installment_unit"`
	NumberOfInstallment int       `json:"number_of_installment"`
	DisbursementDate    time.Time `json:"disbursement_date"`
	FirstDueDate        time.Time `json:"first_due_date"`
	FinalDueDate        time.Time `json:"final_due_date"`
}

type RepaymentResponse struct {
	LoanID                string    `json:"loan_id"`
	PaymentAmount         float64   `json:"payment_amount"`
	InstallmentsPaid      int       `json:"installments_paid"`
	InstallmentAmount     float64   `json:"installment_amount"`
	RemainingInstallments int       `json:"remaining_installments"`
	OutstandingAmount     float64   `json:"outstanding_amount"`
	NextDueDate           time.Time `json:"next_due_date"`
	PaymentDate           time.Time `json:"payment_date"`
}

type OutstandingBalanceResponse struct {
	LoanID                string              `json:"loan_id"`
	CustomerID            string              `json:"customer_id"`
	LoanDetails           LoanDetailsResponse `json:"loan_details"`
	OutstandingAmount     float64             `json:"outstanding_amount"`
	OverdueAmount         float64             `json:"overdue_amount"`
	OverdueInstallments   int                 `json:"overdue_installments"`
	PaidInstallments      int                 `json:"paid_installments"`
	RemainingInstallments int                 `json:"remaining_installments"`
}

type LoanDetailsResponse struct {
	InstallmentUnit   string  `json:"installment_unit"`
	InstallmentAmount float64 `json:"installment_amount"`
	TotalInstallments int     `json:"total_installments"`
}

type DelinquencyResponse struct {
	LoanID                string  `json:"loan_id"`
	CustomerID            string  `json:"customer_id"`
	IsDelinquent          bool    `json:"is_delinquent"`
	InstallmentUnit       string  `json:"installment_unit"`
	OverdueInstallments   int     `json:"overdue_installments"`
	OverdueAmount         float64 `json:"overdue_amount"`
	OutstandingAmount     float64 `json:"outstanding_amount"`
	RequiredPaymentAmount float64 `json:"required_payment_amount"`
}

type LoanScheduleResponse struct {
	LoanID      string                      `json:"loan_id"`
	LoanSummary LoanSummaryScheduleResponse `json:"loan_summary"`
	Schedule    []PaymentScheduleResponse   `json:"schedule"`
}

type LoanSummaryScheduleResponse struct {
	InstallmentUnit   string  `json:"installment_unit"`
	TotalInstallments int     `json:"total_installments"`
	InstallmentAmount float64 `json:"installment_amount"`
	DisbursedAmount   float64 `json:"disbursed_amount"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

type PaymentScheduleResponse struct {
	InstallmentNumber int        `json:"installment_number"`
	DueDate           time.Time  `json:"due_date"`
	InstallmentAmount float64    `json:"installment_amount"`
	InstallmentPaid   float64    `json:"installment_paid"`
	Status            string     `json:"status"`
	PaidDate          *time.Time `json:"paid_date"`
}
