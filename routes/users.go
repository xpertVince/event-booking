package routes

import (
	"net/http"

	"example.com/eveny-booking/models"
	"example.com/eveny-booking/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not save user!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "user created successfully."})
}

func login(context *gin.Context) {
	var user models.User

	// populate the user struct by front-end request
	err := context.ShouldBind(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// generate JWT Token
	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not generate Token!"})
		return
	}

	// valid user, also return token
	context.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
