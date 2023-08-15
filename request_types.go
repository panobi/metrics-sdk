package panobi

import (
	"cloud.google.com/go/civil"
)

// represents a single data point in a timeseries metric
type MetricItem struct {
	Date  civil.Date `json:"date"`
	Value float64    `json:"value"`
}

// used to send a batch of data points for a timeseries metric
type MetricItems struct {
	MetricID string       `json:"metricID"`
	Items    []MetricItem `json:"items"`
}

// represents a single row of data for a non-timeseries metric
type ChartData map[string]interface{}

// used to send a batch of data points for a non-timeseries metric
type RequestChartData struct {
	MetricID string      `json:"metricID"`
	Items    []ChartData `json:"items"`
}

// Used to delete data rows for a metric (timeseries or non-timeseries)
type RequestMetricDataDelete struct {
	MetricID string `json:"metricID"`
}
