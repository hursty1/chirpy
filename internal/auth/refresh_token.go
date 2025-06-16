package auth

import (
	"crypto/rand"
	"encoding/hex"
)


func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	// fmt.Println(key)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	
	// fmt.Println(key)
	hex_token := hex.EncodeToString(key)
	return hex_token, nil
}