package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.`
var hmacSampleSecret []byte

func main() {
	hmacSampleSecret = []byte("u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU")
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	``
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)
}
