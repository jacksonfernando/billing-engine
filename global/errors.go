package global

import "errors"

var (
	INTERNAL_SERVER_ERROR = errors.New("Internal server error")
	ERROR_NOT_FOUND       = errors.New("Your request item not found")
	ERROR_CONFLICT        = errors.New("Your item already exist")
	ERROR_BAD_PARAM_INPUT = errors.New("Given param is not valid")
)

type BadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
