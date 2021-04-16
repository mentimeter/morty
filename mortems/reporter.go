package mortems

import (
	"context"
	"fmt"
	"time"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ReportingService

type ReportingService interface {
	ReportSeverity(string) error
	ReportDetect(detect time.Duration, severity string) error
	ReportResolve(resolve time.Duration, severity string) error
	ReportDowntime(downtime time.Duration, severity string) error
}

type DatadogReporter struct {
	ctx       context.Context
	apiClient *datadog.APIClient
}

func NewDatadogReporter() ReportingService {
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)

	return &DatadogReporter{
		ctx:       ctx,
		apiClient: apiClient,
	}
}

func (d *DatadogReporter) ReportSeverity(severity string) error {
	return d.sendMetric("post_mortems.severity", 1, []string{})
}

func (d *DatadogReporter) ReportDetect(detect time.Duration, severity string) error {
	return d.sendMetric("post_mortems.detect", float64(detect.Minutes()), []string{fmt.Sprintf("severity:%s", severity)})
}

func (d *DatadogReporter) ReportResolve(resolve time.Duration, severity string) error {
	return d.sendMetric("post_mortems.resolve", float64(resolve.Minutes()), []string{fmt.Sprintf("severity:%s", severity)})
}

func (d *DatadogReporter) ReportDowntime(downtime time.Duration, severity string) error {
	return d.sendMetric("post_mortems.downtime", float64(downtime.Minutes()), []string{fmt.Sprintf("severity:%s", severity)})
}

func (d *DatadogReporter) sendMetric(name string, value float64, tags []string) error {
	series := datadog.NewSeries(name, [][]float64{{float64(time.Now().Unix()), float64(value)}})
	series.SetType("count")
	series.SetTags(tags)
	body := *datadog.NewMetricsPayload([]datadog.Series{*series})
	_, _, err := d.apiClient.MetricsApi.SubmitMetrics(d.ctx).Body(body).Execute()
	return err
}
