package utils

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId, userName, email, role string) (string, error) {

	var signedTokenString string
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresIn := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":      userId,
		"username": userName,
		"email":    email,
		"role":     role,
	}

	if jwtExpiresIn != "" {
		duration, err := time.ParseDuration(jwtExpiresIn)
		if err != nil {
			return "", ErrorHandler(err, "Internal error")
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(time.Minute * 15))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrorHandler(err, "Unable to generate token")
	}
	return signedTokenString, nil
}

var JwtStore = JWTStore{
	tokens: make(map[string]time.Time),
}

type JWTStore struct {
	tokens map[string]time.Time
	mu     sync.Mutex
}

// func NewJWTStore() *JWTStore {
// 	return &JWTStore{
// 		tokens: make(map[string]time.Time),
// 	}
// }

func (store *JWTStore) AddToken(token string, exiryTime time.Time) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.tokens[token] = exiryTime
}

func (store *JWTStore) CleanUpExpireTokens() {
	for {
		time.Sleep(time.Second * 2)
		store.mu.Lock()
		for token, timeStamp := range store.tokens {
			if time.Now().After(timeStamp) {
				delete(store.tokens, token)
				log.Println("Cleaning up expired tokens")
			}
		}
		store.mu.Unlock()
	}
}

func (store *JWTStore) IsLoggedOut(token string) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	_, ok := store.tokens[token]
	return ok
}
