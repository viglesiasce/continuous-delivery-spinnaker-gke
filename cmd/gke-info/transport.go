package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeVersionEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.Version(), nil
	}
}

func makeMetaDataEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.MetaData(), nil
	}
}

func makeHealthEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.Health(), nil
	}
}

func makeErrorEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.Error(), nil
	}
}

func makeHomeEndpoint(svc CommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.Home(), nil
	}
}

func decodeNoParamsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponseJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeResponseRaw(_ context.Context, w http.ResponseWriter, response interface{}) error {
	fmt.Fprintf(w, "%s", response)
	return nil
}
