package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

func FormatResponse(c *gin.Context, code int, status string, message string, data interface{}) {
	c.JSON(code, Response{
		Meta: Meta{
			Code:    code,
			Status:  status,
			Message: message,
		},
		Data: data,
	})
}

func ErrorResponse(c *gin.Context, code int, status string, message string) {
	c.JSON(code, Response{
		Meta: Meta{
			Code:    code,
			Status:  status,
			Message: message,
		},
		Data: nil,
	})
}

func BadRequestResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Status:  "bad request",
			Message: message,
		},
		Data: data,
	})
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	FormatResponse(c, code, "success", message, data)
}
