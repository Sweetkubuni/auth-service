package controllers

import (
	"auth-service/database"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterResponse struct {
	UserId   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// RegisterUser godoc
//
//	@Summary		Regiser a new user
//	@Description	Register a new user in the sytem
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.User			true	"User details for registration"
//	@Success		201		{object}	RegisterResponse	"User created successfully"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		500		{object}	map[string]string	"Internal Server Error"
//	@Router			/user/register [post]
func (app *App) RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "email": user.Email, "username": user.Username})
}
