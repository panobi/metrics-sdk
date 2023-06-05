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
type client struct {
	t *transport
}

// Creates a new client with the given key information.
func CreateClient(k KeyInfo) *client {
	c := &client{
		t: createTransport(k),
	}

	return c
}

func (client *client) Close() {
}

// Sends a single metric item to your Panobi workspace.
func (client *client) SendMetricItem(metricID string, item MetricItem) error {
	return client.SendMetricItems(metricID, []MetricItem{item})
}

// Sends multiple metric items to your Panobi workspace.
func (client *client) SendMetricItems(metricID string, items []MetricItem) error {
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

	_, err = client.t.post(b)
	return err
}
