package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/viglesiasce/gke-info/pkg/stackdriver"

	"context"

	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"

	"net/http"
	"net/http/httputil"

	"cloud.google.com/go/errors"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/trace"
)

type stackDriverClient struct {
	projectID         string
	serviceName       string
	version           string
	monitoringService *monitoring.Service
	cloudLogger       *logging.Logger
	traceClient       *trace.Client
	errorsClient      *errors.Client
}

func NewStackDriverClient(ctx context.Context, projectID string, serviceName string, version string) (*stackDriverClient, error) {

	monitoringService, err := newMonitoringService(ctx, projectID, stackdriver.MetricType(serviceName))
	if err != nil {
		return nil, fmt.Errorf("Failed creating monitoring service: %v", err)
	}
	cloudLogger, err := newCloudLogger(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed creating cloud logger: %v", err)
	}
	traceClient, err := newTraceClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("Failed creating trace client: %v", err)
	}
	errorsClient, err := newErrorsClient(ctx, projectID, serviceName, version)
	if err != nil {
		return nil, fmt.Errorf("Failed creating errors client: %v", err)
	}
	return &stackDriverClient{projectID, serviceName, version, monitoringService, cloudLogger, traceClient, errorsClient}, nil
}

func newTraceClient(ctx context.Context, projectID string) (*trace.Client, error) {
	traceClient, err := trace.NewClient(context.Background(), projectID)
	if err != nil {
		panic(fmt.Sprintf("Failed to create trace client: %v", err))
	}
	policy, err := trace.NewLimitedSampler(1.0, 5)
	if err != nil {
		panic("Failed to create sampler: " + err.Error())
	}
	traceClient.SetSamplingPolicy(policy)
	return traceClient, err
}

func newErrorsClient(ctx context.Context, projectID string, serviceName string, version string) (*errors.Client, error) {
	errors, err := errors.NewClient(ctx, projectID, serviceName, version, true)
	if err != nil {
		return nil, fmt.Errorf("Failed to create errors client: %v", err)
	}
	return errors, nil
}

func newCloudLogger(ctx context.Context, projectID string) (*logging.Logger, error) {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logging client: %v", err))
	}
	logging := client.Logger("gceme-log")
	return logging, nil
}

func newMonitoringService(ctx context.Context, projectID string, metricType string) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	monitoringService, err := monitoring.New(hc)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring service: %v", err)
	}
	if err := stackdriver.CreateCustomMetric(monitoringService, projectID, metricType, "s"); err != nil {
		return nil, fmt.Errorf("failed to create custom metric: %v", err)
	}
	for {
		resp, err := stackdriver.GetCustomMetric(monitoringService, projectID, metricType)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve custom metric: %v", err)
		}
		if len(resp.MetricDescriptors) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return monitoringService, nil
}

type stackDriverMiddleware struct {
	context context.Context
	sdc     *stackDriverClient
	logger  log.Logger
	next    commonService
}

func logRequest(method string, r *http.Request, begin time.Time, mw stackDriverMiddleware) {
	elapsed := time.Since(begin)
	// Log request
	raw, _ := httputil.DumpRequest(r, true)
	_ = mw.logger.Log(
		"method", method,
		"took", elapsed,
		"request", raw,
	)
}

func writeMetric(endpoint string, begin time.Time, mw stackDriverMiddleware) {
	if err := stackdriver.WriteTimeSeriesValue(mw.sdc.monitoringService, mw.sdc.projectID,
		mw.sdc.serviceName+"/"+endpoint, time.Since(begin).Seconds()); err != nil {
		fmt.Printf("Failed to send metric value: %v", err)
	}
}

func (mw stackDriverMiddleware) Version(r *http.Request) (version string) {
	// Catch any panics and report them to Error Reporting service
	endpoint := "version"
	defer mw.sdc.errorsClient.Catch(mw.context)
	defer writeMetric(endpoint, time.Now(), mw)
	defer logRequest(endpoint, r, time.Now(), mw)
	version = mw.next.Version(r)
	return
}

func (mw stackDriverMiddleware) MetaData(r *http.Request) (instance *Instance) {
	// Catch any panics and report them to Error Reporting service
	endpoint := "metadata"
	defer mw.sdc.errorsClient.Catch(mw.context)
	defer writeMetric(endpoint, time.Now(), mw)
	defer logRequest(endpoint, r, time.Now(), mw)
	instance = mw.next.MetaData(r)
	return
}

func (mw stackDriverMiddleware) Health(r *http.Request) (status string) {
	// Catch any panics and report them to Error Reporting service
	endpoint := "health"
	defer mw.sdc.errorsClient.Catch(mw.context)
	defer writeMetric(endpoint, time.Now(), mw)
	defer logRequest(endpoint, r, time.Now(), mw)
	status = mw.next.Health(r)
	return
}

func (mw stackDriverMiddleware) Error(r *http.Request) (err error) {
	// Catch any panics and report them to Error Reporting service
	endpoint := "error"
	defer mw.sdc.errorsClient.Catch(mw.context)
	defer writeMetric(endpoint, time.Now(), mw)
	defer logRequest(endpoint, r, time.Now(), mw)
	err = mw.next.Error(r)
	return
}

func (mw stackDriverMiddleware) Home(r *http.Request) (doc string) {
	// Catch any panics and report them to Error Reporting service
	endpoint := "home"
	defer mw.sdc.errorsClient.Catch(mw.context)
	defer writeMetric(endpoint, time.Now(), mw)
	defer logRequest(endpoint, r, time.Now(), mw)
	doc = mw.next.Home(r)
	return
}
