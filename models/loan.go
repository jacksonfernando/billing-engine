package models

import (
	"time"
)

// User represents the users table (for reference only)
type User struct {
	ID         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID string     `json:"customer_id" gorm:"uniqueIndex;not null;type:varchar(36)"`
	Name       string     `json:"name" gorm:"type:varchar(255)"`
	Email      string     `json:"email" gorm:"type:varchar(255)"`
	Phone      string     `json:"phone" gorm:"type:varchar(50)"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}

// DisbursementDetail represents the disbursement_details table
type DisbursementDetail struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	LoanID            string     `json:"loan_id" gorm:"uniqueIndex;not null;type:varchar(50)"`
	CustomerID        string     `json:"customer_id" gorm:"not null;type:varchar(36);index"`
	DisbursementDate  time.Time  `json:"disbursement_date" gorm:"not null"`
	DisbursedAmount   float64    `json:"disbursed_amount" gorm:"not null;type:decimal(15,2)"`
	DisbursedCurrency string     `json:"disbursed_currency" gorm:"default:'IDR';type:char(3)"`
	Status            string     `json:"status" gorm:"not null;type:varchar(100)"`
	CreatedAt         time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy         string     `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy         string     `json:"updated_by" gorm:"type:varchar(255)"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"index"`
}

// LoanSummary represents the loan_summary table
type LoanSummary struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	LoanID                string     `json:"loan_id" gorm:"uniqueIndex;not null;type:varchar(50)"`
	CustomerID            string     `json:"customer_id" gorm:"not null;type:varchar(36);index"`
	PrincipalAmount       float64    `json:"principal_amount" gorm:"not null;type:decimal(15,2)"`
	InterestAmount        float64    `json:"interest_amount" gorm:"not null;type:decimal(15,2)"`
	OutstandingAmount     float64    `json:"outstanding_amount" gorm:"not null;type:decimal(15,2)"`
	NoOfInstallment       int        `json:"no_of_installment" gorm:"not null"`
	InstallmentUnit       string     `json:"installment_unit" gorm:"not null;type:varchar(100)"`
	InstallmentAmount     float64    `json:"installment_amount" gorm:"not null;type:decimal(15,2)"`
	EffectiveInterestRate float64    `json:"effective_interest_rate" gorm:"not null;type:decimal(5,4)"`
	Status                string     `json:"status" gorm:"not null;type:varchar(100);index"`
	LoanStartDate         time.Time  `json:"loan_start_date" gorm:"not null;type:date"`
	CreatedAt             time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy             string     `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedAt             time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy             string     `json:"updated_by" gorm:"type:varchar(255)"`
	DeletedAt             *time.Time `json:"deleted_at" gorm:"index"`
}

// PaymentSchedule represents the payment_schedule table
type PaymentSchedule struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	LoanID             string     `json:"loan_id" gorm:"not null;type:varchar(50);index"`
	InstallmentNumber  int        `json:"installment_number" gorm:"not null"`
	InstallmentAmount  float64    `json:"installment_amount" gorm:"not null;type:decimal(15,2)"`
	InstallmentDueDate time.Time  `json:"installment_due_date" gorm:"not null;type:date;index"`
	InstallmentPaid    float64    `json:"installment_paid" gorm:"not null;type:decimal(15,2);default:0"`
	Status             string     `json:"status" gorm:"default:'PENDING';type:varchar(100);index"`
	Currency           string     `json:"currency" gorm:"default:'IDR';type:char(3)"`
	CreatedAt          time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy          string     `json:"created_by" gorm:"type:varchar(255)"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy          string     `json:"updated_by" gorm:"type:varchar(255)"`
	DeletedAt          *time.Time `json:"deleted_at" gorm:"index"`
	DeletedBy          string     `json:"deleted_by" gorm:"type:varchar(255)"`
}

// PaymentScheduleHistory represents the payment_schedule_history table
type PaymentScheduleHistory struct {
	ID                 uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID         uint      `json:"schedule_id" gorm:"not null;index"`
	LoanID             string    `json:"loan_id" gorm:"not null;type:varchar(50);index"`
	Action             string    `json:"action" gorm:"not null;type:varchar(100)"`
	InstallmentNumber  int       `json:"installment_number" gorm:"not null"`
	InstallmentAmount  float64   `json:"installment_amount" gorm:"not null;type:decimal(15,2)"`
	InstallmentDueDate time.Time `json:"installment_due_date" gorm:"not null;type:date"`
	Status             string    `json:"status" gorm:"type:varchar(100)"`
	Currency           string    `json:"currency" gorm:"default:'IDR';type:char(3)"`
	CreatedAt          time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;index"`
	CreatedBy          string    `json:"created_by" gorm:"type:varchar(255)"`
}

// Constants
const (
	StatusPending    = "PENDING"
	StatusPaid       = "PAID"
	StatusDelinquent = "DELINQUENT"

	InstallmentUnitWeek  = "week"
	InstallmentUnitMonth = "month"

	CurrencyIDR = "IDR"

	ActionPayment = "PAYMENT"
)
