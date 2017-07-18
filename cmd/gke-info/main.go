package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	"context"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-kit/kit/log"
)

var (
	version = os.Getenv("VERSION")
)

func main() {
	// Create Local logger
	localLogger := log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	projectID := metadata.ProjectID
	serviceName := "gke-info"
	serviceComponent := os.Getenv("COMPONENT")
	backendURL := os.Getenv("BACKEND_URL")
	sdc, err := NewStackDriverClient(ctx, projectID, serviceName+"-"+serviceComponent, version)
	if err != nil {
		panic("Unable to create stackdriver clients: " + err.Error())
	}

	var common CommonService
	common = commonService{backendURL: backendURL, sdc: sdc}
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
