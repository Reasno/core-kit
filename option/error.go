package option

import (
	"context"
	"net/http"

	"github.com/DoNewsCode/core/srvhttp"
)

// ErrorEncoder is a go kit style http error encoder. Internally it uses
// srvhttp.ResponseEncoder to encode the error.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	encoder := srvhttp.NewResponseEncoder(w)
	encoder.EncodeError(err)
}
