package auth

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Service is a simple CRUD interface for users
type loggingService struct {
	logger log.Logger
	next   Service
}

//NewLoggingService returns a new instance of logging Service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) IsUserValid(a AuthModel) (b bool, err error) {
	defer func(since time.Time) {
		s.logger.Log(
			"method", "IsUserValid",
			"username", a.Name,
			"took", time.Since(since),
			"err", err,
		)
	}(time.Now())
	return s.next.IsUserValid(a)
}
func (s *loggingService) GetToken(a AuthModel) (str string, err error) {
	defer func(since time.Time) {
		s.logger.Log(
			"method", "GetToken",
			"username", a.Name,
			"took", time.Since(since),
			"err", err,
		)
	}(time.Now())
	return s.next.GetToken(a)
}
func (s *loggingService) ParseToken(tokenString string) (b bool, err error) {
	defer func(since time.Time) {
		s.logger.Log(
			"method", "ParseToken",
			"token", tokenString,
			"took", time.Since(since),
			"err", err,
		)
	}(time.Now())
	return s.next.ParseToken(tokenString)
}

func (s *loggingService) GetUser(tokenString string) (string, error) {
	return s.next.GetUser(tokenString)
}
