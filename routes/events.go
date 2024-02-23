package routes

import (
	"net/http"
	"strconv"

	"example.com/eveny-booking/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch events. Try again later!"})
		return
	}

	// do not return, transform to json by gin
	context.JSON(http.StatusOK, events) // gin.H any type, string:any
}

func getEvent(context *gin.Context) {
	// decimal system: 10, base
	// 64L bit size, int64
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64) // extract the id from incoming request
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID!"})
		return
	}

	event, nil := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch event!"})
		return
	}

	context.JSON(http.StatusOK, event)
}

// always get gin context when use a func as endpoint handler func
func createEvent(context *gin.Context) {
	// BEFORE DO ANYTHING ELSE, SHOULD EXTRACT TOKEN FROM INCOMING REQUEST

	var event models.Event
	err := context.ShouldBindJSON(&event) // store the data from the body, return type is error

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	// event created succressfully
	// user ID should be the creator user_id
	userId := context.GetInt64("userId") // get userid from middleware
	event.UserID = userId

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not create events. Try again later!"})
		return
	}

	// send the event back to user
	context.JSON(http.StatusCreated, gin.H{"Message": "Event created", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64) // extract the id from incoming request
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID!"})
		return
	}

	userId := context.GetInt64("userId") // get userId from token
	event, err := models.GetEventByID(eventId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch the event!"})
		return
	}

	// check if user is updating his own event (creator)
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized to update this event!"})
		return
	}

	// create new event and populate updated info
	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent) // populate the updated event

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	updatedEvent.ID = eventId   // use the ID generated (existing ID), to the Struct
	err = updatedEvent.Update() // update in DB

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not update the event!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Message": "Event updated successfully!"})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64) // extract the id from incoming request
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Message": "Could not parse event ID!"})
		return
	}

	// get userId from token authenticated
	userId := context.GetInt64("userId")

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not fetch the event!"})
		return
	}

	// check if user is updating his own event (creator)
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "User not authorized to delete this event!"})
		return
	}

	// when everything is done, delete from DB
	err = event.Delete()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Message": "Could not delete the event!"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"Message": "Event deleted successfully!"})
}
