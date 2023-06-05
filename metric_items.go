package panobi

import (
	"cloud.google.com/go/civil"
)

type MetricItem struct {
	Date  civil.Date `json:"date"`
	Value float64    `json:"value"`
}

type MetricItems struct {
	MetricID string       `json:"metricID"`
	Items    []MetricItem `json:"items"`
}
