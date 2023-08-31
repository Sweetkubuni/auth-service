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

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get the profile of a specific user
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User profile"
// @Failure 404 {object} map[string]string "User not found"
// @Router /user/{id} [get]
func (app *App) GetUserProfile(context *gin.Context) {
	userID := context.Param("id")
	var user models.User
	result := database.Instance.First(&user, userID)
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	context.JSON(http.StatusOK, user)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update the profile of a specific user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body models.User true "User profile details for update"
// @Success 200 {object} models.User "Updated user profile"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Router /user/{id} [put]
func (app *App) UpdateUserProfile(context *gin.Context) {
	userID := context.Param("id")
	var user models.User
	result := database.Instance.First(&user, userID)
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update user attributes
	database.Instance.Save(&user)
	context.JSON(http.StatusOK, user)
}

// DeleteUserProfile godoc
// @Summary Delete user profile
// @Description Delete the profile of a specific user
// @Tags user
// @Param id path int true "User ID"
// @Success 204 "No content"
// @Failure 404 {object} map[string]string "User not found"
// @Router /user/{id} [delete]
func (app *App) DeleteUserProfile(context *gin.Context) {
	userID := context.Param("id")
	var user models.User
	result := database.Instance.First(&user, userID)
	if result.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Delete user profile
	database.Instance.Delete(&user)
	context.Status(http.StatusNoContent)
}
