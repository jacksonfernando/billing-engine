package global

import (
	"billing-engine/models"
	"time"
)

type PostgresDefault struct {
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UpdatedBy string    `json:"updated_by,omitempty"`
	DeletedBy string    `json:"deleted_by,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// Standardized API Response Structures

// DisbursementSuccessResponse represents a successful disbursement response
type DisbursementSuccessResponse struct {
	Status string                       `json:"status"`
	Data   *models.DisbursementResponse `json:"data"`
}

// RepaymentSuccessResponse represents a successful repayment response
type RepaymentSuccessResponse struct {
	Status string                    `json:"status"`
	Data   *models.RepaymentResponse `json:"data"`
}

// OutstandingBalanceSuccessResponse represents a successful outstanding balance query response
type OutstandingBalanceSuccessResponse struct {
	Status string                             `json:"status"`
	Data   *models.OutstandingBalanceResponse `json:"data"`
}

// DelinquencySuccessResponse represents a successful delinquency query response
type DelinquencySuccessResponse struct {
	Status string                      `json:"status"`
	Data   *models.DelinquencyResponse `json:"data"`
}

// LoanScheduleSuccessResponse represents a successful loan schedule query response
type LoanScheduleSuccessResponse struct {
	Status string                       `json:"status"`
	Data   *models.LoanScheduleResponse `json:"data"`
}
