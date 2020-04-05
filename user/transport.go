package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s Service) (http.Handler, http.Handler) {

	userEndpoints := MakeServerEndpoints(s)

	storeHandler := httptransport.NewServer(
		userEndpoints.StoreEndpoint,
		decodeStoreRequest,
		httptransport.EncodeJSONResponse,
	)

	findHandler := httptransport.NewServer(
		userEndpoints.FindEndpoint,
		decodeFindRequest,
		encodeerror,
	)
	listHandler := httptransport.NewServer(
		userEndpoints.ListEndpoint,
		decodeListRequest,
		httptransport.EncodeJSONResponse,
	)
	rpwless := chi.NewRouter()
	r := chi.NewRouter()
	rpwless.Post("/store", func(w http.ResponseWriter, r *http.Request) {
		storeHandler.ServeHTTP(w, r)
	})

	r.Handle("/find", findHandler)
	r.Handle("/list", listHandler)
	return r, rpwless
}

type errorer interface {
	error() error
}

func decodeStoreRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req storeRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeFindRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req findRequest
	vals := r.URL.Query()
	id := vals.Get("userid")
	userID, _ := strconv.Atoi(id)

	req.UserID = UserID(userID)

	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req listRequest

	return req, nil
}

//encodeerror first access to request
func encodeerror(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case ErrUnknown:
		w.WriteHeader(http.StatusAccepted)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
