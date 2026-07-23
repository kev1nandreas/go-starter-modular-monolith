package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"-"`
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Error      string      `json:"error,omitempty"`
}

func NewSuccessResponse(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    message,
		Data:       data,
	}
}

func NewPaginatedResponse(message string, data interface{}, meta interface{}) Response {
	return Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    message,
		Data:       data,
		Meta:       meta,
	}
}

func NewErrorResponse(statusCode int, message string, err string) Response {
	return Response{
		StatusCode: statusCode,
		Success:    false,
		Message:    message,
		Error:      err,
	}
}

func (r Response) Send(ctx *gin.Context) {
	ctx.JSON(r.StatusCode, r)
}
