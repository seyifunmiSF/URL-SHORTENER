package utils

import (
	"math/rand"
	"time"
)

// GenerateUniqueID generates a unique ID using cryptographic random numbers
func GenerateUniqueID() (string, error) {
	// Define the character set containing alphabetic characters
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Create a random generator based on the character set
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	// Generate a unique ID by randomly selecting 5 letters from the character set
	uniqueID := make([]byte, 5)
	for i := range uniqueID {
		uniqueID[i] = charSet[rnd.Intn(len(charSet))]
	}

	return string(uniqueID), nil
}
