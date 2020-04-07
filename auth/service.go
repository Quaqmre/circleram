package auth

import (
	"errors"
	"time"

	"github.com/Quaqmre/circleramkit/user"
	"github.com/dgrijalva/jwt-go"
)

// Service is a simple CRUD interface for users

type Service interface {
	IsUserValid(AuthModel) (bool, error)
	GetToken(AuthModel) (string, error)
	ParseToken(string) (bool, error)
	GetUser(tokenString string) (string, error)
}

type AuthModel struct {
	Name     string
	Password string
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// AuthModel cloning User Model

// ErrUnknown is used when a cargı could not be found
var ErrUnknown = errors.New("user not exist")
var ErrWorngContent = errors.New("user username or password is wrong")
var ErrInvalidPassword = errors.New("invalid Password struct")
var ErrInvalidUserName = errors.New("invalid username struct")

var ErrUnauthorize = errors.New("Unauthorize")

type service struct {
	secret      []byte
	userService user.Service
}

//NewUserService export for using inmem repo the other package
func NewAuthService(s string, srv user.Service) Service {
	return &service{
		secret:      []byte(s),
		userService: srv,
	}
}

func (s *service) IsUserValid(a AuthModel) (bool, error) {

	user, err := s.userService.FindWithUserName(a.Name)
	if err != nil {
		return false, err
	}

	if user.Password == a.Password {
		return true, nil
	}

	return false, ErrWorngContent
}
func (s *service) GetToken(a AuthModel) (string, error) {

	b, err := s.IsUserValid(a)

	if err != nil {
		return "", err
	}

	if !b {
		return "", ErrUnauthorize
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Username: a.Name,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.secret)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *service) ParseToken(tokenString string) (bool, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return false, err
	}
	if !tkn.Valid {
		return false, ErrUnauthorize
	}

	return tkn.Valid, nil
}
func (s *service) GetUser(tokenString string) (string, error) {

	b, err := s.ParseToken(tokenString)

	if b {
		claims := &Claims{}
		tkn, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		})
		if claims, ok := tkn.Claims.(*Claims); ok && tkn.Valid {
			return claims.Username, nil
		}
		return "", errors.New("not found user")
	}
	return "", err
}
