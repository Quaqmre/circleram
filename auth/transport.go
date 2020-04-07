package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Quaqmre/circleramkit/user"
	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s Service) http.Handler {

	authEndpoints := MakeAuthEndpoints(s)

	getTokenHandler := httptransport.NewServer(
		authEndpoints.GetTokenEndpoint,
		decodeGetTokenRequest,
		encodeResponse,
	)

	r := chi.NewRouter()
	r.Handle("/gettoken", getTokenHandler)
	return r
}

type errorer interface {
	error() error
}

func decodeGetTokenRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getTokenRequest
	if e := json.NewDecoder(r.Body).Decode(&req.authModel); e != nil {
		return nil, e
	}
	return req, nil
}

//encodeerror first access to request
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a kit transport error, only business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
func codeFrom(err error) int {
	switch err {
	case user.ErrUnknown:
		return http.StatusNotFound
	case ErrWorngContent:
		return http.StatusNotFound
	case ErrUnauthorize:
		return http.StatusUnauthorized
	case ErrInvalidPassword, ErrInvalidUserName:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
