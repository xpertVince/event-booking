package routes

import (
	"example.com/eveny-booking/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// dont need to return anything, original server
	server.GET("/events", getEvents)    // incoming get request
	server.GET("/events/:id", getEvent) // get event by ID

	// protected by Authorization
	// group the routes that needs to be authenticated
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate) // use which middleware, run middleware before the handlers
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	// registration also need authentication
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)
	//
	// server.POST("/events", middlewares.Authenticate, createEvent) // execute from left to right
	// server.PUT("/events/:id", updateEvent)                        // update an event
	// server.DELETE("/events/:id", deleteEvent)

	// user related
	server.POST("/signup", signup)
	server.POST("/login", login)
}
