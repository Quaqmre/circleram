package customer

import (
	"time"

	"github.com/go-kit/kit/log"
)

// Service is a simple CRUD interface for customers
type loggingService struct {
	logger log.Logger
	next   Service
}

//NewLoggingService returns a new instance of logging Service
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Store(c *Customer) (err error) {
	defer func(since time.Time) {
		s.logger.Log(
			"method", "Store",
			"took", time.Since(since),
			"err", err,
		)
	}(time.Now())
	return s.next.Store(c)
}

func (s *loggingService) Find(id CustomerID) (c *Customer, err error) {
	defer func(begin time.Time) {
		s.logger.Log("method", "Find",
			"customerId", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.next.Find(id)
}
