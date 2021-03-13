package kits3

import (
	"github.com/DoNewsCode/core/ots3"
	"net/http"
	"testing"

	"github.com/DoNewsCode/core"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestModule(t *testing.T) {
	c := core.New(core.WithInline("s3.default.accessKey", ""))
	c.ProvideEssentials()
	c.Provide(ots3.Providers())
	c.AddModuleFunc(New)
	router := mux.NewRouter()
	c.ApplyRouter(router)
	request, _ := http.NewRequest("POST", "/upload", nil)
	assert.True(t, router.Match(request, &mux.RouteMatch{}))
}
