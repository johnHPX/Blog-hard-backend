package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SecretKey = "X46LJy3eytWRGQlVFXGqVXC/QUnI/OcVtIPpzCtpPdySMOs9PGfTqanJ5Ri3RVPugjA3BwW0hW4H8LveAoRhLw=="

func CreateToken(userID, kind string) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 4).Unix()
	permissions["userID"] = userID
	permissions["kind"] = kind
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(SecretKey))
}

func validateToken(r *http.Request) error {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, returnCheckKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("token invalido")
}

type userToken struct {
	UserID string
	Kind   string
}

func ExtractUserID(r *http.Request) (userToken, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, returnCheckKey)
	if err != nil {
		return userToken{}, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := permissions["userID"]
		kind := permissions["kind"]

		return userToken{
			UserID: userID.(string),
			Kind:   kind.(string),
		}, nil
	}

	return userToken{}, errors.New("token inválido")
}

func ExtractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	// Bearer asdlkdjsakl -> asdlkdjsakl
	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return token
}
func returnCheckKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}

	return []byte(SecretKey), nil
}
