package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestError struct {
	Message string `json:"message"`
}

func CtxRequestErrorJson(ctx *gin.Context, msg string) {

	ctx.JSON(http.StatusBadRequest, RequestError{
		Message: msg,
	})
	ctx.Abort()
}
func CtxServerErrorJson(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, RequestError{
		Message: msg,
	})
	ctx.Abort()
}
