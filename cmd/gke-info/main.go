package main

import (
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	version = "v6.1.5"
)

func main() {
	// Create Local logger
	localLogger := log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	projectID := "vic-goog"
	serviceName := "gke-info"
	sdc, err := NewStackDriverClient(ctx, projectID, serviceName, version)
	if err != nil {
		panic("Unable to create stackdriver clients: " + err.Error())
	}

	var common CommonService
	common = commonService{}
	common = stackDriverMiddleware{ctx, sdc, localLogger, common.(commonService)}
	metaDataHandler := httptransport.NewServer(
		makeMetaDataEndpoint(common),
		decodeNoParamsRequest,
		encodeResponse,
	)

	versionHandler := httptransport.NewServer(
		makeVersionEndpoint(common),
		decodeNoParamsRequest,
		encodeResponse,
	)

	healthHandler := httptransport.NewServer(
		makeHealthEndpoint(common),
		decodeNoParamsRequest,
		encodeResponse,
	)

	errorHandler := httptransport.NewServer(
		makeErrorEndpoint(common),
		decodeNoParamsRequest,
		encodeResponse,
	)

	http.Handle("/metadata", sdc.traceClient.HTTPHandler(metaDataHandler))
	http.Handle("/version", sdc.traceClient.HTTPHandler(versionHandler))
	http.Handle("/health", sdc.traceClient.HTTPHandler(healthHandler))
	http.Handle("/error", sdc.traceClient.HTTPHandler(errorHandler))
	localLogger.Log("msg", "HTTP", "addr", ":8080")
	localLogger.Log("err", http.ListenAndServe(":8080", nil))
}
