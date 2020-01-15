package message

import (
	"crypto/rand"
	"fmt"
)

// GetRandomHexString returns a random hexadecimal string. It returns an error if it has any problem with the local PRNG.
func GetRandomHexString(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

