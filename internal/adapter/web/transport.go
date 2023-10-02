package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// Endpoint defines endpoint logic
type Endpoint interface {
	Do(ctx context.Context, request interface{}) (response interface{}, err error)
}

// Decoder defines behavior to extract a user-domain request object from an HTTP
// request object. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward Decoder could be something that
// JSON decodes from the request body to the concrete request type.
type Decoder interface {
	Decode(context.Context, *http.Request) (request interface{}, err error)
}

// Encoder encodes the passed response object to the HTTP response
// writer. It's designed to be used in HTTP servers, for server-side
// endpoints. One straightforward Encoder could be something that
// JSON encodes the object directly to the response body.
type Encoder interface {
	Encode(context.Context, http.ResponseWriter, interface{}) error
}

type Handler struct {
	endpoint Endpoint
	decoder  Decoder
	encoder  Encoder
	logger   *slog.Logger
}

// ErrorResponse define response.
type ErrorResponse struct {
	Message string `json:"message"`
}

const (
	jsonContentType = "application/json; charset=utf-8"
)

var (
	defaultErrorResponse = []byte(`{"message": "unable to process request"}`)
)

func NewRouter() *mux.Router {
	return mux.NewRouter()
}

func NewHandler() *Handler {
	newHandler := Handler{}

	return &newHandler
}

func (h *Handler) WithEndpoint(endpoint Endpoint) *Handler {
	h.endpoint = endpoint

	return h
}

func (h *Handler) WithDecoder(decoder Decoder) *Handler {
	h.decoder = decoder

	return h
}

func (h *Handler) WithEncoder(encoder Encoder) *Handler {
	h.encoder = encoder

	return h
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var err error

	ctx := req.Context()

	request, err := h.decoder.Decode(ctx, req)
	if err != nil {
		h.encodeError(err, rw)
		return
	}

	response, err := h.endpoint.Do(ctx, request)
	if err != nil {
		h.encodeError(err, rw)
		return
	}

	err = h.encoder.Encode(ctx, rw, response)
	if err != nil {
		h.encodeError(err, rw)
		return
	}
}

func (h *Handler) encodeError(err error, w http.ResponseWriter) {
	newErrorMessage := ErrorResponse{
		Message: err.Error(),
	}

	content, errMarshal := json.Marshal(newErrorMessage)
	if errMarshal != nil {
		h.logger.Error("unable to marshal error response into json", "error", errMarshal)

		content = defaultErrorResponse
	}

	w.Header().Set("Content-Type", jsonContentType)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(content)
}
