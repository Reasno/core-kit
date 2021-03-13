package mw_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DoNewsCode/core-kit/mw"
	"github.com/DoNewsCode/core/key"
	"github.com/DoNewsCode/core/logging"
	"github.com/go-kit/kit/endpoint"
)

func ExampleError() {
	var (
		err      error
		original endpoint.Endpoint
		wrapped  endpoint.Endpoint
	)
	original = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, errors.New("error")
	}
	_, err = original(context.Background(), nil)
	fmt.Printf("%T\n", err)

	wrapped = mw.Error(mw.ErrorOption{})(original)

	_, err = wrapped(context.Background(), nil)
	fmt.Printf("%T\n", err)

	// Output:
	// *errors.errorString
	// *unierr.Error
}

func ExampleLog() {
	var (
		original endpoint.Endpoint
		wrapped  endpoint.Endpoint
	)
	original = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "respData", nil
	}

	wrapped = mw.Log(
		logging.NewLogger("json"),
		key.New(),
		false,
	)(original)

	wrapped(context.Background(), "reqData")

	// Output:
	// {"request":"reqData","response":"respData"}
}

func ExampleRetry() {
	var (
		original endpoint.Endpoint
		wrapped  endpoint.Endpoint
	)
	original = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		fmt.Println("attempt")
		return nil, errors.New("")
	}

	wrapped = mw.Retry(mw.RetryOption{
		Max:     5,
		Timeout: time.Second,
	})(original)

	wrapped(context.Background(), nil)

	// Output:
	// attempt
	// attempt
	// attempt
	// attempt
	// attempt
}

func ExampleTimeout() {
	var (
		original endpoint.Endpoint
		wrapped  endpoint.Endpoint
	)
	original = func(ctx context.Context, request interface{}) (response interface{}, err error) {
		select {
		case <-ctx.Done():
			return nil, errors.New("timeout")
		case <-time.After(100000 * time.Microsecond):
			return nil, nil
		}
	}

	wrapped = mw.Timeout(time.Microsecond)(original)
	_, err := wrapped(context.Background(), nil)
	fmt.Println(err)

	// Output:
	// timeout
}
