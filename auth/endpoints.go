package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose a user service
type Endpoints struct {
	GetTokenEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes.
func MakeAuthEndpoints(s Service) Endpoints {

	var getTokenEndpoint endpoint.Endpoint
	{
		getTokenEndpoint = makeGetTokenEndpoint(s)
	}

	return Endpoints{
		GetTokenEndpoint: getTokenEndpoint,
	}
}

// MakeStoreEndpoints returns an endpoint via the passed service.
// Primarily useful in a server.
func makeGetTokenEndpoint(s Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTokenRequest)
		tokenStr, e := s.GetToken(req.authModel)

		if e != nil {
			return getTokenResponse{Err: e}, nil
		}

		return getTokenResponse{Token: tokenStr, Err: e}, nil
	}
}

type getTokenRequest struct {
	authModel AuthModel
}

type getTokenResponse struct {
	Token string
	Err   error `json:"err,omitempty"`
}

func (r getTokenResponse) error() error { return r.Err }
