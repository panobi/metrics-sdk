package panobi

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// maximum items to send in one request
	MaxItems           int           = 1000
	bufferedSendPeriod time.Duration = 10 * time.Second
	errMaxNumberSize   string        = "%s cannot be larger than %d %s"
)

// Client for pushing metrics items to your Panobi workspace.
type Client struct {
	t *transport
}

// Creates a new client with the given key information.
func CreateClient(k KeyInfo) *Client {
	c := &Client{
		t: createTransport(k),
	}

	return c
}

func (client *Client) Close() {
}

// Sends a single metric item to your Panobi workspace.
func (client *Client) SendMetricItem(metricID string, item MetricItem) error {
	return client.SendMetricItems(metricID, []MetricItem{item})
}

// Sends multiple metric items to your Panobi workspace.
func (client *Client) SendMetricItems(metricID string, items []MetricItem) error {
	if len(items) > MaxItems {
		return fmt.Errorf(errMaxNumberSize, "batch", MaxItems, "MetricItems")
	}

	b, err := json.Marshal(&MetricItems{
		MetricID: metricID,
		Items:    items,
	})
	if err != nil {
		return err
	}

	_, err = client.t.post(TimeseriesURI, b)
	return err
}

// Sends metric chart data rows to your Panobi workspace.
func (client *Client) SendMetricChartData(metricID string, items []ChartData) error {
	if len(items) > MaxItems {
		return fmt.Errorf(errMaxNumberSize, "batch", MaxItems, "ChartData")
	}

	b, err := json.Marshal(&RequestChartData{
		MetricID: metricID,
		Items:    items,
	})
	if err != nil {
		return err
	}

	_, err = client.t.post(ChartDataURI, b)
	return err
}

// Delete all stored rows for a metric (timeseries or non-timeseries)
func (client *Client) DeleteMetricData(metricID string) error {
	b, err := json.Marshal(&RequestMetricDataDelete{
		MetricID: metricID,
	})
	if err != nil {
		return err
	}

	_, err = client.t.post(DeleteURI, b)
	return err
}
