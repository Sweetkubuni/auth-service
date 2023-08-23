package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
//
//	@Summary		Ping the secured endpoint
//	@Description	Returns a pong message
//	@Tags			utils
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	map[string]string	"Ping sucessful"
//	@Router			/secured/ping [get]
func (app *App) Ping(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "pong"})
}
