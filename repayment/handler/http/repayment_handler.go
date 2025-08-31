package http

import (
	"net/http"

	"billing-engine/global"
	"billing-engine/middlewares"
	"billing-engine/models"
	"billing-engine/repayment"
	"billing-engine/utils/validator"

	"github.com/labstack/echo/v4"
)

type RepaymentHandler struct {
	repaymentService repayment.RepaymentServiceInterface
	middleware       middlewares.GoMiddlewareInterface
}

// NewRepaymentHandler creates a new repayment handler instance
func NewRepaymentHandler(e *echo.Echo, repaymentService repayment.RepaymentServiceInterface, middleware middlewares.GoMiddlewareInterface) {
	handler := &RepaymentHandler{
		repaymentService: repaymentService,
		middleware:       middleware,
	}

	// Register routes
	v1 := e.Group("/v1")
	v1.POST("/repayment", handler.ProcessRepayment)
}

func (h *RepaymentHandler) ProcessRepayment(c echo.Context) error {
	var req models.RepaymentRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, global.BadResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate request using validator
	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, global.BadResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Process repayment
	response, err := h.repaymentService.ProcessRepayment(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, global.BadResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, global.RepaymentSuccessResponse{
		Status: "success",
		Data:   response,
	})
}
