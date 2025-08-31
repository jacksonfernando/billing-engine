package http

import (
	"net/http"

	"billing-engine/global"
	"billing-engine/loan_query"
	"billing-engine/middlewares"

	"github.com/labstack/echo/v4"
)

type LoanQueryHandler struct {
	loanQueryService loan_query.LoanQueryServiceInterface
	middleware       middlewares.GoMiddlewareInterface
}

// NewLoanQueryHandler creates a new loan query handler instance
func NewLoanQueryHandler(e *echo.Echo, loanQueryService loan_query.LoanQueryServiceInterface, middleware middlewares.GoMiddlewareInterface) {
	handler := &LoanQueryHandler{
		loanQueryService: loanQueryService,
		middleware:       middleware,
	}

	// Register routes
	v1 := e.Group("/v1")
	v1.GET("/loans/:loan_id/outstanding", handler.GetOutstandingBalance)
	v1.GET("/loans/:loan_id/delinquency", handler.GetDelinquencyStatus)
	v1.GET("/loans/:loan_id/schedule", handler.GetLoanSchedule)
}

func (h *LoanQueryHandler) GetOutstandingBalance(c echo.Context) error {
	loanID := c.Param("loan_id")
	if loanID == "" {
		return c.JSON(http.StatusBadRequest, global.BadResponse{
			Code:    http.StatusBadRequest,
			Message: "Loan ID is required",
		})
	}

	response, err := h.loanQueryService.GetOutstandingBalance(c.Request().Context(), loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, global.BadResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, global.OutstandingBalanceSuccessResponse{
		Status: "success",
		Data:   response,
	})
}

func (h *LoanQueryHandler) GetDelinquencyStatus(c echo.Context) error {
	loanID := c.Param("loan_id")
	if loanID == "" {
		return c.JSON(http.StatusBadRequest, global.BadResponse{
			Code:    http.StatusBadRequest,
			Message: "Loan ID is required",
		})
	}

	response, err := h.loanQueryService.GetDelinquencyStatus(c.Request().Context(), loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, global.BadResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, global.DelinquencySuccessResponse{
		Status: "success",
		Data:   response,
	})
}

func (h *LoanQueryHandler) GetLoanSchedule(c echo.Context) error {
	loanID := c.Param("loan_id")
	if loanID == "" {
		return c.JSON(http.StatusBadRequest, global.BadResponse{
			Code:    http.StatusBadRequest,
			Message: "Loan ID is required",
		})
	}

	response, err := h.loanQueryService.GetLoanSchedule(c.Request().Context(), loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, global.BadResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, global.LoanScheduleSuccessResponse{
		Status: "success",
		Data:   response,
	})
}
