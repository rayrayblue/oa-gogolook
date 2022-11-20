package domain

import "fmt"

var serviceCode = "TASK"

var (
	ErrInvalidParameters = NewErrorResponse(fmt.Sprintf("ERR_%s_0001", serviceCode), "invalid uri parameters")
	ErrInvalidPayload    = NewErrorResponse(fmt.Sprintf("ERR_%s_0002", serviceCode), "invalid request payload")
	ErrSystemError       = NewErrorResponse(fmt.Sprintf("ERR_%s_0003", serviceCode), "system error")
	ErrNetworkError      = NewErrorResponse(fmt.Sprintf("ERR_%s_0004", serviceCode), "network error")
	ErrUnauthorized      = NewErrorResponse(fmt.Sprintf("ERR_%s_0005", serviceCode), "unauthorized")
	ErrServiceTimeout    = NewErrorResponse(fmt.Sprintf("ERR_%s_0006", serviceCode), "service timeout")
	ErrDataNotFound      = NewErrorResponse(fmt.Sprintf("ERR_%s_0007", serviceCode), "data not found")
	ErrWrongID           = NewErrorResponse(fmt.Sprintf("ERR_%s_0008", serviceCode), "wrong task ID")
	ErrTaskNameNotMatch  = NewErrorResponse(fmt.Sprintf("ERR_%s_0009", serviceCode), "task name not match")
)

type ErrorResponse interface {
	Error() string
}

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func (e *errorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(errorCode string, message string) *errorResponse {
	return &errorResponse{
		ErrorCode: errorCode,
		Message:   message,
	}
}
