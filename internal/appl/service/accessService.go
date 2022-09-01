package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/johnHPX/blog-hard-backend/internal/infra/repository"
	"github.com/johnHPX/blog-hard-backend/internal/infra/utils/messages"
)

type userToken struct {
	UserID string
	Kind   string
}

type accessServiceInterface interface {
	CreateAToken(userID, kind string) (string, error)
	CreateRToken() (string, error)
	ValidateAToken(r *http.Request) error
	ValidateRToken(rtoken string) error
	ExtractTokenInfo(r *http.Request) (*userToken, error)
	ExtractInvalideToken(r *http.Request) (string, error)
	GenerateNewToken(userID string) (string, error)
	GenerateTokenRecovery(userID string) (string, error)
	ValidateAndExtractTokenRecovery(r *http.Request) (string, error)
}

type accessServiceImpl struct {
	SecretKey string
}

func (s *accessServiceImpl) GenerateNewToken(userID string) (string, error) {
	// repository access
	repAcces := repository.NewAccessRepository()
	// find rtoken by userID of a user
	access, err := repAcces.FindToken(userID)
	if err != nil {
		return "", err
	}

	// verific if rtoken was blocked
	if access.IsBlocked {
		return "", errors.New(messages.TokenBlocked)
	}

	// validate rtoken
	err = s.ValidateRToken(access.Token)
	if err != nil {
		err := repAcces.RemoveToken(userID)
		if err != nil {
			return "", err
		}

		return "", err
	}

	// create a new rtoken
	newRToken, err := s.CreateRToken()
	if err != nil {
		return "", err
	}

	// update rtoken
	err = repAcces.UpdateToken(newRToken, userID)
	if err != nil {
		return "", err
	}

	// find user by userID
	repUser := repository.NewUserRepository()
	user, err := repUser.Find(access.UserID)
	if err != nil {
		return "", err
	}

	// create a new atoken
	newAToken, err := s.CreateAToken(user.UserID, user.Kind)
	if err != nil {
		return "", err
	}

	// return atoken

	return newAToken, nil
}

func (s *accessServiceImpl) CreateAToken(userID, kind string) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 4).Unix()
	permissions["userID"] = userID
	permissions["kind"] = kind
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *accessServiceImpl) CreateRToken() (string, error) {
	permissions := jwt.MapClaims{}
	permissions["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *accessServiceImpl) GenerateTokenRecovery(userID string) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Minute * 10).Unix()
	permissions["userID-recovery"] = userID
	permissions["recovery"] = true
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *accessServiceImpl) ValidateAToken(r *http.Request) error {
	tokenString := s.getToken(r)
	token, err := jwt.Parse(tokenString, s.returnCheckKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New(messages.InvalideToken)
}

func (s *accessServiceImpl) ValidateRToken(rtoken string) error {
	token, err := jwt.Parse(rtoken, s.returnCheckKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New(messages.InvalideToken)
}

func (s *accessServiceImpl) ValidateAndExtractTokenRecovery(r *http.Request) (string, error) {
	tokenString := s.getToken(r)
	token, err := jwt.Parse(tokenString, s.returnCheckKey)
	if err != nil {
		return "", err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := permissions["userID-recovery"]
		recovery := permissions["recovery"]

		if !recovery.(bool) {
			return "", errors.New(messages.InvalideToken)
		}

		return userID.(string), nil
	}

	return "", errors.New(messages.InvalideToken)
}

func (s *accessServiceImpl) ExtractTokenInfo(r *http.Request) (*userToken, error) {
	tokenString := s.getToken(r)
	token, err := jwt.Parse(tokenString, s.returnCheckKey)
	if err != nil {
		return nil, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := permissions["userID"]
		kind := permissions["kind"]

		return &userToken{
			UserID: userID.(string),
			Kind:   kind.(string),
		}, nil
	}

	return nil, errors.New(messages.InvalideToken)
}

func (s *accessServiceImpl) ExtractInvalideToken(r *http.Request) (string, error) {
	tokenString := s.getToken(r)
	token, err := jwt.Parse(tokenString, s.returnCheckKey)
	if err != nil {
		if err.Error() == "Token is expired" {
			permissions, ok := token.Claims.(jwt.MapClaims)
			if ok {
				userID := permissions["userID"]
				return userID.(string), nil
			}
		}
	}

	return "", err
}

func (s *accessServiceImpl) getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	// Bearer asdlkdjsakl -> asdlkdjsakl

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return token
}
func (s *accessServiceImpl) returnCheckKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("método de assinatura inesperado! %v", token.Header["alg"])
	}

	return []byte(s.SecretKey), nil
}

func NewAccessService() accessServiceInterface {
	return &accessServiceImpl{
		SecretKey: "X46LJy3eytWRGQlVFXGqVXC/QUnI/OcVtIPpzCtpPdySMOs9PGfTqanJ5Ri3RVPugjA3BwW0hW4H8LveAoRhLw==",
	}
}
