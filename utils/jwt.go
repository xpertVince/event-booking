package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// in reality, should be very complex
const secretKey = "supersecret"

func GenerateToken(email string, userId int64) (string, error) {
	// generate with new token (Claims: with data)
	// first argument: sign-in method, Signature of the method
	// second argument: JWP claims: expiration time is 2 hours, Unix format -> "exp" used internally
	// return a token pointer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // DO NOT Include Password !!!!
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	// convert token to string to client, need a Key to verify incoming tokens!!!!!!
	return token.SignedString([]byte(secretKey)) // return a token or error.
	// Key must be []byte!!!!
}

// return user_id and error
func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) { // secret key, error
		// check if the value stored in Method is the HMAC type
		_, ok := token.Method.(*jwt.SigningMethodHMAC) // type checking, end version HMAC. return bool if it is the correct type
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, errors.New("Could not parse token")
	}

	// extract data from token
	tokenIsValid := parsedToken.Valid
	if !tokenIsValid {
		return 0, errors.New("Invalid Token!") // 0 as user_id for no user
	}

	// if valid token: get holded data, Map type claims
	// check if it is map type claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Invalid token claism")
	}

	// // access data in the claims
	// email := claims["email"].(string) // tell email is string
	userId := int64(claims["userId"].(float64)) // actual value for userId (float), should convert to int64
	return userId, nil                          // no error
}
