package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	panobi "github.com/panobi/metrics-sdk"
)

func main() {
	//
	// We need the name of the JSON file.
	//

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>\n", os.Args[0])
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
	// Open the file and read it. The JSON structure is assumed to be an array of RequestMetricsSDKItems (from openapi.yaml)
	//

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	var metrics []panobi.MetricItems
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(bytes, &metrics)
	if err != nil {
		log.Fatal("Error decoding json:", err)
	}

	for _, metric := range metrics {
		//
		// send to Panobi in batches
		//
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
