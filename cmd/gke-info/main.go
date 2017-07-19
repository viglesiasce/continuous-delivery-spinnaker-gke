package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/go-kit/kit/log"
)

var (
	version = os.Getenv("VERSION")
)

func main() {
	// Create Local logger
	localLogger := log.NewLogfmtLogger(os.Stderr)
	serviceComponent := os.Getenv("COMPONENT")
	backendURL := os.Getenv("BACKEND_URL")

	var common CommonService
	common = commonService{backendURL: backendURL}

	createCommonEndpoints(common)
	if serviceComponent == "frontend" {
		createFrontendEndpoints(common)
	} else if serviceComponent == "backend" {
		createBackendEndpoints(common)
	} else {
		panic("Unknown component: " + serviceComponent)
	}

	localLogger.Log("msg", "HTTP", "addr", ":8080")
	localLogger.Log("err", http.ListenAndServe(":8080", nil))
}
