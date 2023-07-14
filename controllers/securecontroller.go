package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
