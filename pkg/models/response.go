package models

type ErrorResponse struct {
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
	Success      bool        `json:"success"`
}

func NewErrorResponse(code int, message error) *ErrorResponse {
	return &ErrorResponse{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: message.Error(),
	}
}
func NewErrorsResponse(code int, messages string) *ErrorResponse {
	return &ErrorResponse{
		Success:      false,
		ErrorCode:    code,
		ErrorMessage: messages,
	}
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total,omitempty"`
}

type Response[T any] struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Success      bool   `json:"success"`
	Data         T      `json:"data"`
	Total        int64  `json:"total"`
}

func NewSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
	}
}

func NewSuccessPagingResponse(data interface{}, total int64) *SuccessResponse {
	return &SuccessResponse{
		Success: true,
		Data:    data,
		Total:   total,
	}
}
