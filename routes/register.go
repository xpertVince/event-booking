package routes

import (
	"net/http"
	"strconv"

	"example.com/eveny-booking/models"
	"github.com/gin-gonic/gin"
)

func registerForEvent(context *gin.Context) {
	userId := context.GetInt64("userId")                          // get user id from token in authentication process
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64) // extract the id from incoming request
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID!"})
		return
	}

	// look for event
	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event!"})
		return
	}

	// add registration into DB table
	err = event.Register(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not register event for the user!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"Message": "Registered!"})
}

func cancelRegistration(context *gin.Context) {
	userId := context.GetInt64("userId")                          // get user id from token in authentication process
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64) // extract the id from incoming request
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID!"})
		return
	}

	// event struct, because CancelRegistration is a Method !!! But no need to populate all, event ID is enough
	var event models.Event
	event.ID = eventId // assign the ID

	err = event.CancelRegistration(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not cancel the registration!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Message": "Cancelled!"})
}
