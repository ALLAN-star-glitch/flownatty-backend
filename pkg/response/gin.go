package response

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Send sends a BaseResponse with the appropriate HTTP status
func Send(c *gin.Context, status int, resp *BaseResponse) {
    c.JSON(status, resp)
}

// Success sends a 200 OK response
func Success(c *gin.Context, message string, data interface{}) {
    Send(c, http.StatusOK, &BaseResponse{
        Success: true,
        Message: message,
        Code:    http.StatusOK,
        Data:    data,
    })
}

// SuccessWithPagination sends a 200 OK response with pagination
func SuccessWithPagination(c *gin.Context, message string, data interface{}, total int64, limit, offset int) {
    hasMore := total > int64(offset+limit)
    Send(c, http.StatusOK, &BaseResponse{
        Success: true,
        Message: message,
        Code:    http.StatusOK,
        Data:    data,
        Pagination: &PaginationMeta{
            Total:   total,
            Limit:   limit,
            Offset:  offset,
            HasMore: hasMore,
        },
    })
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
    Send(c, http.StatusCreated, &BaseResponse{
        Success: true,
        Message: message,
        Code:    http.StatusCreated,
        Data:    data,
    })
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusBadRequest, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusBadRequest,
        Error: &ErrorPayload{
            Code:    http.StatusBadRequest,
            Details: details,
        },
    })
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusUnauthorized, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusUnauthorized,
        Error: &ErrorPayload{
            Code:    http.StatusUnauthorized,
            Details: details,
        },
    })
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusForbidden, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusForbidden,
        Error: &ErrorPayload{
            Code:    http.StatusForbidden,
            Details: details,
        },
    })
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusNotFound, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusNotFound,
        Error: &ErrorPayload{
            Code:    http.StatusNotFound,
            Details: details,
        },
    })
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusConflict, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusConflict,
        Error: &ErrorPayload{
            Code:    http.StatusConflict,
            Details: details,
        },
    })
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string, details interface{}) {
    Send(c, http.StatusInternalServerError, &BaseResponse{
        Success: false,
        Message: message,
        Code:    http.StatusInternalServerError,
        Error: &ErrorPayload{
            Code:    http.StatusInternalServerError,
            Details: details,
        },
    })
}