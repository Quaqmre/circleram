package customer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a customer service
type Endpoints struct {
	StoreEndpoint endpoint.Endpoint
	FindEndpoint  endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes.
func MakeServerEndpoints(s Service) Endpoints {
	var storeEndpoint endpoint.Endpoint
	{
		storeEndpoint = makeStoreEndpoints(s)
		// storeEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Minute), 1))(storeEndpoint)
	}

	var findEndpoint endpoint.Endpoint
	{
		findEndpoint = makeFindEndpoints(s)
	}

	return Endpoints{
		StoreEndpoint: storeEndpoint,
		FindEndpoint:  findEndpoint,
	}
}

// MakeStoreEndpoints returns an endpoint via the passed service.
// Primarily useful in a server.
func makeStoreEndpoints(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(storeRequest)
		e := s.Store(&req.Customer)
		return storeResponse{Err: e}, nil
	}
}

// MakeFindEndpoints returns an endpoint via the passed service.
// Primarily useful in a server.
func makeFindEndpoints(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(findRequest)
		c, e := s.Find(req.CustomerID)
		if e != nil {
			return findResponse{Err: e}, nil
		}
		return findResponse{Customer: *c, Err: e}, nil
	}
}

type storeRequest struct {
	Customer Customer
}

type storeResponse struct {
	Err error `json:"err,omitempty"`
}

func (r storeResponse) error() error { return r.Err }

type findRequest struct {
	CustomerID CustomerID
}
type findResponse struct {
	Customer Customer
	Err      error `json:"err,omitempty"`
}

func (r findResponse) error() error { return r.Err }
