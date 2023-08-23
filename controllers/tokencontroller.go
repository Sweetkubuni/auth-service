package controllers

import (
	"auth-service/auth"
	"auth-service/database"
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Token Request represents the credentials required to fetch a token.
//
// TokenRequest godoc
//	@Schema		TokenRequest
//	@Property	email (string) "User's email address" true
//	@Property	password(string) "User's password" true
type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// GenerateToken godoc
//
//	@Summary		Generate a token
//	@Description	Generate a JWT token for a user given their email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		TokenRequest		true	"Credentials for token generation"
//	@Success		200		{object}	TokenResponse		"Token generated successfully"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized"
//	@Failure		500		{object map[string]string "Internal Server Error"
//	@Router			/token [post]
func (app *App) GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Id, user.Email, user.Username, app.Config.GetAppConfig().JwtSecret)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
