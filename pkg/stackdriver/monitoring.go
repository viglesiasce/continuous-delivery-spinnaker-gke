package stackdriver

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/monitoring/v3"
)

func projectResource(projectID string) string {
	return "projects/" + projectID
}

// formatResource marshals a response objects as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}
func MetricType(serviceName string) string {
	return fmt.Sprintf("custom.googleapis.com/%s", serviceName)
}

// createCustomMetric creates a custom metric specified by the metric type.
func CreateCustomMetric(s *monitoring.Service, projectID, serviceName string, unit string) error {
	//ld := monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "An arbitrary measurement"}
	md := monitoring.MetricDescriptor{
		Type: MetricType(serviceName),
		//Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "DOUBLE",
		Unit:        unit,
		Description: serviceName,
		DisplayName: serviceName,
	}
	resp, err := s.Projects.MetricDescriptors.Create(projectResource(projectID), &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

// getCustomMetric reads the custom metric created.
func GetCustomMetric(s *monitoring.Service, projectID, serviceName string) (*monitoring.ListMetricDescriptorsResponse, error) {
	resp, err := s.Projects.MetricDescriptors.List(projectResource(projectID)).
		Filter(fmt.Sprintf("metric.type=\"%s\"", MetricType(serviceName))).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not get custom metric: %v", err)
	}

	log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// writeTimeSeriesValue writes a value for the custom metric created
func WriteTimeSeriesValue(s *monitoring.Service, projectID, serviceName string, value float64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: MetricType(serviceName),
			// Labels: map[string]string{
			// 	"environment": "STAGING",
			// },
		},
		// Resource: &monitoring.MonitoredResource{
		// 	Labels: map[string]string{
		// 		"instance_id": "test-instance",
		// 		"zone":        "us-central1-f",
		// 	},
		// 	Type: "gce_instance",
		// },
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					DoubleValue: &value,
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := s.Projects.TimeSeries.Create(projectResource(projectID), &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}
