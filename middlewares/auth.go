package middlewares

import (
	"net/http"

	"example.com/eveny-booking/utils"
	"github.com/gin-gonic/gin"
)

// run in the middle of request, so use abort with status
func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization") // token is in Authorization header
	if token == "" {                                     // if client returns no Token in incoming request
		// abort the current response, send this response
		// no other code on the server runs
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized !"})
		return
	}

	// Do have a token, but NOT valid
	userId, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token Not Authorized."})
		return
	}

	// add some data to context value: Key:Value
	context.Set("userId", userId)

	// valid token: next request will be executed
	context.Next()
}
