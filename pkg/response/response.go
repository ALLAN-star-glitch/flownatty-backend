package response

// ErrorPayload represents error details
type ErrorPayload struct {
    Code    int         `json:"code"`
    Details interface{} `json:"details"`
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
    Total   int64 `json:"total"`
    Limit   int   `json:"limit"`
    Offset  int   `json:"offset"`
    HasMore bool  `json:"hasMore"`
}

// BaseResponse is the standard response wrapper - UNIFORM FOR ALL RESPONSES
type BaseResponse struct {
    Success    bool            `json:"success"`
    Message    string          `json:"message"`
    Code       int             `json:"code"`
    Data       interface{}     `json:"data,omitempty"`
    Pagination *PaginationMeta `json:"pagination,omitempty"`
    Error      *ErrorPayload   `json:"error,omitempty"`
}

// NewResponse creates a new BaseResponse
func NewResponse(success bool, message string, code int, data interface{}, pagination *PaginationMeta, err *ErrorPayload) *BaseResponse {
    return &BaseResponse{
        Success:    success,
        Message:    message,
        Code:       code,
        Data:       data,
        Pagination: pagination,
        Error:      err,
    }
}