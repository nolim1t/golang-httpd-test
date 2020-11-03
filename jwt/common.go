package jwt

// Define imports
import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"strings"
	"time"
)

// Local methods
func stringtofile(filename string) string {
	byte_output, err := ioutil.ReadFile(filename)
	if err == nil {
		return strings.Trim(string(byte_output), "\n")
	} else {
		return "-1"
	}
}

// Methods
func SignKey(keyfile string, Username string) string {
	// Get byte output for filename
	signed_key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		fmt.Printf("Error has occured reading file: %s", err)
		return "-1"
	}
	// If no error
	claims := jwt.MapClaims{}
	claims["user"] = Username
	// 2 hours
	claims["exp"] = time.Now().Add(time.Minute * 60 * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_string, errSign := token.SignedString(signed_key)
	if errSign != nil {
		fmt.Printf("Error signing string: %s\n", errSign)
		return "-2"
	}
	return signed_string
}
func ValidateKey(keyfile string, Token string) (string, error) {
	jwtfile, err := ioutil.ReadFile(keyfile)
	if err != nil {
		fmt.Printf("Error has occured reading file: %s", err)
		return "-1", errors.New("Error reading file")
	}
	token, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
		return jwtfile, nil
	})
	if token.Valid {
		// TODO: Return a struct with valid and error
		claims := token.Claims.(jwt.MapClaims)
		fmt.Printf("user: %s\n", claims["user"])
		return "valid", nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "not a token", errors.New("Not a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return "not active", errors.New("Token expired")
		} else {
			return "unhandled token", errors.New("Unhandled token exception")
		}
	} else {
		return "-2", errors.New("Token not valid")
	}
}
