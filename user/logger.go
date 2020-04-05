package user

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

func (s *loggingService) Store(c *User) (err error) {
	defer func(since time.Time) {
		s.logger.Log(
			"method", "Store",
			"took", time.Since(since),
			"err", err,
		)
	}(time.Now())
	return s.next.Store(c)
}

func (s *loggingService) Find(id UserID) (c *User, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Find",
			"userId", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Find(id)
}
func (s *loggingService) List() (c []User, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "List",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.List()
}
func (s *loggingService) FindWithUserName(name string) (c *User, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "FindWithUserName",
			"userId", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.FindWithUserName(name)
}
