package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"golang.org/x/time/rate"
)

// Endpoints collects all of the endpoints that compose a user service
type Endpoints struct {
	StoreEndpoint endpoint.Endpoint
	FindEndpoint  endpoint.Endpoint
	ListEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes.
func MakeServerEndpoints(s Service) Endpoints {
	var storeEndpoint endpoint.Endpoint
	{
		storeEndpoint = makeStoreEndpoints(s)
		storeEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Minute), 15))(storeEndpoint)
	}

	var findEndpoint endpoint.Endpoint
	{
		findEndpoint = makeFindEndpoints(s)
	}

	var listdEndpoint endpoint.Endpoint
	{
		listdEndpoint = makeListEndpoints(s)
	}

	return Endpoints{
		StoreEndpoint: storeEndpoint,
		FindEndpoint:  findEndpoint,
		ListEndpoint:  listdEndpoint,
	}
}

// MakeStoreEndpoints returns an endpoint via the passed service.
// Primarily useful in a server.
func makeStoreEndpoints(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(storeRequest)
		e := s.Store(&req.User)
		return storeResponse{Err: e}, nil
	}
}

// MakeFindEndpoints returns an endpoint via the passed service.
// Primarily useful in a server.
func makeFindEndpoints(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findRequest)
		c, e := s.Find(req.UserID)
		if e != nil {
			return findResponse{Err: e}, nil
		}
		return findResponse{User: *c, Err: e}, nil
	}
}

func makeListEndpoints(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(listRequest)
		c, e := s.List()
		if e != nil {
			return listResponse{Err: e}, nil
		}
		return listResponse{Users: c, Err: e}, nil
	}
}

type storeRequest struct {
	User User
}

type storeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r storeResponse) error() error { return r.Err }

type findRequest struct {
	UserID UserID
}
type findResponse struct {
	User User
	Err  error `json:"err,omitempty"`
}

func (r findResponse) error() error { return r.Err }

type listRequest struct{}
type listResponse struct {
	Users []User
	Err   error `json:"err,omitempty"`
}
