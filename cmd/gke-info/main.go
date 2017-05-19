package main

import (
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	version = "v7.0.0"
)

func createFrontendEndpoints(common CommonService, sdc *stackDriverClient) {
	homeHandler := httptransport.NewServer(
		makeHomeEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseRaw,
	)
	http.Handle("/", sdc.traceClient.HTTPHandler(homeHandler))
	http.Handle("/home", sdc.traceClient.HTTPHandler(homeHandler))
}

func createBackendEndpoints(common CommonService, sdc *stackDriverClient) {
	metaDataHandler := httptransport.NewServer(
		makeMetaDataEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/metadata", sdc.traceClient.HTTPHandler(metaDataHandler))
}

func createCommonEndpoints(common CommonService, sdc *stackDriverClient) {
	versionHandler := httptransport.NewServer(
		makeVersionEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/version", sdc.traceClient.HTTPHandler(versionHandler))

	healthHandler := httptransport.NewServer(
		makeHealthEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/health", sdc.traceClient.HTTPHandler(healthHandler))

	errorHandler := httptransport.NewServer(
		makeErrorEndpoint(common),
		decodeNoParamsRequest,
		encodeResponseJSON,
	)
	http.Handle("/error", sdc.traceClient.HTTPHandler(errorHandler))
}

func main() {
	// Create Local logger
	localLogger := log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	projectID := "vic-goog"
	serviceName := "gke-info"
	serviceComponent := os.Getenv("COMPONENT")
	sdc, err := NewStackDriverClient(ctx, projectID, serviceName+"-"+serviceComponent, version)
	if err != nil {
		panic("Unable to create stackdriver clients: " + err.Error())
	}

	var common CommonService
	common = commonService{}
	common = stackDriverMiddleware{ctx, sdc, localLogger, common.(commonService)}

	createCommonEndpoints(common, sdc)
	if serviceComponent == "frontend" {
		createFrontendEndpoints(common, sdc)
	} else if serviceComponent == "backend" {
		createBackendEndpoints(common, sdc)
	} else {
		panic("Unknown component: " + serviceComponent)
	}

	localLogger.Log("msg", "HTTP", "addr", ":8080")
	localLogger.Log("err", http.ListenAndServe(":8080", nil))
}
