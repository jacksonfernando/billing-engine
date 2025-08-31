package http

import (
	"net/http"

	"billing-engine/disbursement"
	"billing-engine/global"
	"billing-engine/middlewares"
	"billing-engine/models"
	"billing-engine/utils/validator"

	"github.com/labstack/echo/v4"
)

type DisbursementHandler struct {
	disbursementService disbursement.DisbursementServiceInterface
	middleware          middlewares.GoMiddlewareInterface
}

// NewDisbursementHandler creates a new disbursement handler instance
func NewDisbursementHandler(e *echo.Echo, disbursementService disbursement.DisbursementServiceInterface, middleware middlewares.GoMiddlewareInterface) {
	handler := &DisbursementHandler{
		disbursementService: disbursementService,
		middleware:          middleware,
	}

	// Register routes
	v1 := e.Group("/v1")
	v1.POST("/disbursement", handler.CreateDisbursement)
}

func (h *DisbursementHandler) CreateDisbursement(c echo.Context) error {
	var req models.DisbursementRequest

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

	// Create disbursement
	response, err := h.disbursementService.CreateDisbursement(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, global.BadResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, global.DisbursementSuccessResponse{
		Status: "success",
		Data:   response,
	})
}
