package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type error_message struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, status_code int, message string) {
	logrus.Errorf(message)
	ctx.AbortWithStatusJSON(status_code, error_message{message})
}
