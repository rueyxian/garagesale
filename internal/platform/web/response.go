package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// ================================================================================
// Response
func Response(ctx context.Context, w http.ResponseWriter, v interface{}, statusCode int) error {

	ctxVal, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return errors.New("web values missing from context")
	}

	ctxVal.StatusCode = statusCode

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, err := json.Marshal(&v)
	if err != nil {
		return errors.Wrap(err, "marshalling value to json")
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "writing to client")
	}
	return nil
}

// ================================================================================
// ResponseError knows how to handle errors going out to the client
func ResponseError(ctx context.Context, w http.ResponseWriter, err error) error {

	if webErr, ok := errors.Cause(err).(*Error); ok {
		resp := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		return Response(ctx, w, resp, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Response(ctx, w, resp, http.StatusInternalServerError)

}
