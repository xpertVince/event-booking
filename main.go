package main

import (
	"example.com/eveny-booking/db"
	"example.com/eveny-booking/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	server := gin.Default() // start a engine (http server), get a pointer of engine

	routes.RegisterRoutes(server) // outsource the logics

	server.Run(":8080") // listen for incoming requests, domain is localhost:8080
}
