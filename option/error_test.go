package option

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/DoNewsCode/core/unierr"
	"github.com/stretchr/testify/assert"
)

func TestErrorEncoder(t *testing.T) {
	var err = unierr.InternalErr(errors.New("server bug"), "whoops")
	recorder := httptest.NewRecorder()
	ErrorEncoder(context.Background(), err, recorder)
	resp := recorder.Result()
	assert.Equal(t, "500 Internal Server Error", resp.Status)
	content, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"code\":13,\"message\":\"whoops\"}\n", string(content))
}
