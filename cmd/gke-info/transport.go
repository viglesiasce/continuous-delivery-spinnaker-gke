package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeVersionEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		return svc.Version(r), nil
	}
}

func makeMetaDataEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		return svc.MetaData(r), nil
	}
}

func makeHealthEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		return svc.Health(r), nil
	}
}

func makeErrorEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		return svc.Error(r), nil
	}
}

func makeHomeEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*http.Request)
		return svc.Home(r), nil
	}
}

func decodeNoParamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponseJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeResponseRaw(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Fprintf(w, "%s", response)
	return nil
}

func createFrontendEndpoints(common CommonService) {
	homeHandler := httptransport.NewServer(
		makeHomeEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseRaw,
	)
	http.Handle("/", homeHandler)
}

func createBackendEndpoints(common CommonService) {
	metaDataHandler := httptransport.NewServer(
		makeMetaDataEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/metadata", metaDataHandler)
}

func createCommonEndpoints(common CommonService) {
	versionHandler := httptransport.NewServer(
		makeVersionEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/version", versionHandler)

	healthHandler := httptransport.NewServer(
		makeHealthEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/health", healthHandler)

	errorHandler := httptransport.NewServer(
		makeErrorEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/error", errorHandler)
}
