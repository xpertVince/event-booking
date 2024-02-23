package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// first argument is byte slice, so convert first
	// second argument: const, complexity of the hash
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err // convert to string
}

// compare the password
func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil // true if no error
}
