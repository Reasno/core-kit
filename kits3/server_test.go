// +build integration

package kits3

import (
	"context"
	"github.com/DoNewsCode/core/ots3"
	"github.com/opentracing/opentracing-go"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/DoNewsCode/core/config"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Parallel()
	manager := setupManager()
	service := UploadService{
		Logger: log.NewNopLogger(),
		S3:     manager,
	}
	endpoint := MakeUploadEndpoint(&service)
	endpoint = Middleware(log.NewNopLogger(), config.EnvTesting)(endpoint)
	handler := httptransport.NewServer(endpoint, decodeRequest, httptransport.EncodeJSONResponse)
	ln, _ := net.Listen("tcp", ":8888")
	server := &http.Server{
		Handler: handler,
	}
	go server.Serve(ln)
	defer server.Shutdown(context.Background())

	uri, _ := url.Parse("http://localhost:8888/")
	uploader := NewClientUploaderFromUrl(uri)
	urlStr, err := uploader.Upload(context.Background(), "foo", strings.NewReader("bar"))
	assert.NoError(t, err)
	assert.NotEmpty(t, urlStr)
}

func setupManager() *ots3.Manager {
	return setupManagerWithTracer(nil)
}

func setupManagerWithTracer(tracer opentracing.Tracer) *ots3.Manager {
	m := ots3.NewManager(
		"Q3AM3UQ867SPQQA43P2F",
		"zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		"https://play.minio.io:9000",
		"asia",
		"mybucket",
		ots3.WithTracer(tracer),
	)
	_ = m.CreateBucket(context.Background(), "mybucket")
	return m
}
