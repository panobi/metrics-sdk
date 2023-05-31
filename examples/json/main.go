package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	panobi "github.com/panobi/metrics-sdk"
)

func main() {
	//
	// We need the name of the JSON file.
	//

	if len(os.Args) < 2 || os.Args[1] == "-t" && len(os.Args) < 3 {
		log.Fatalf("Usage: %s [-t] <filename>\n", os.Args[0])
	}

	//
	// You can find your key in your Panobi workspace's integration settings.
	// It is safer to load it from an environment variable than to paste it
	// directly into this code; do not put secrets in GitHub.
	//

	k, err := panobi.ParseKey(os.Getenv("METRICS_SDK_SIGNING_KEY"))
	if err != nil {
		log.Fatal("Error parsing key:", err)
	}

	//
	// Create a client with the signing key information.
	//

	client := panobi.CreateClient(k)
	defer client.Close()

	//
	// Open the file and read it. The JSON structure is assumed to be an array of RequestMetricsSDKTimeseries or RequestChartData (from openapi.yaml)
	//

	var timeseries bool
	var path string
	if os.Args[1] == "-t" {
		timeseries = true
		path = os.Args[2]
	} else {
		path = os.Args[1]
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	if timeseries {
		sendTimeseriesData(file, client)
	} else {
		sendChartData(file, client)
	}
}

// sends data for non-timeseries metrics, without deleting existing data. only new rows are stored
func sendTimeseriesData(file *os.File, client *panobi.Client) {
	var metrics []panobi.MetricItems
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		log.Fatal("Error decoding json:", err)
	}

	//
	// send to Panobi in batches
	//
	for _, metric := range metrics {

		itemCount := len(metric.Items)
		for i := 0; i < itemCount; i += panobi.MaxItems {
			rangeEnd := itemCount
			if i+panobi.MaxItems < itemCount {
				rangeEnd = i + panobi.MaxItems
			}
			err := client.SendMetricItems(metric.MetricID, metric.Items[i:rangeEnd])
			if err != nil {
				log.Fatalf("Error sending items for metricID %s: %s", metric.MetricID, err.Error())
			}

			log.Printf("Successfully sent %d item(s) for metricID %s", len(metric.Items[i:rangeEnd]), metric.MetricID)
		}
	}
}

// sends data for non-timeseries metrics, deleting existing data first
func sendChartData(file *os.File, client *panobi.Client) {
	var metrics []panobi.RequestChartData
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		log.Fatal("Error decoding json:", err)
	}

	for _, metric := range metrics {
		// delete existing data for this metric first
		err := client.DeleteMetricData(metric.MetricID)
		if err != nil {
			log.Fatalf("Error deleting existing data for metricID %s: %s", metric.MetricID, err.Error())
		}
		itemCount := len(metric.Items)
		for i := 0; i < itemCount; i += panobi.MaxItems {
			rangeEnd := itemCount
			if i+panobi.MaxItems < itemCount {
				rangeEnd = i + panobi.MaxItems
			}
			err := client.SendMetricChartData(metric.MetricID, metric.Items[i:rangeEnd])
			if err != nil {
				log.Fatalf("Error sending items for metricID %s: %s", metric.MetricID, err.Error())
			}

			log.Printf("Successfully sent %d item(s) for metricID %s", len(metric.Items[i:rangeEnd]), metric.MetricID)
		}
	}
}
