package http

// ErrorResponse represent the response error struct
type ErrorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Message: err.Error(),
	}
}
