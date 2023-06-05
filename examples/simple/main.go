package main

import (
	"log"
	"os"

	"cloud.google.com/go/civil"
	panobi "github.com/panobi/metrics-sdk"
)

func main() {
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
	// Push a single item.
	//
	// The Date should be a day value
	// The Value can be an int or float
	// The metricID should be retreived from Panobi's Metrics configuration page
	//

	item := panobi.MetricItem{
		Date:  civil.Date{Year: 2023, Month: 7, Day: 1},
		Value: 1000,
	}
	metricID := "XRnrRBTedmWzy8RQ6pqh2d"

	if err := client.SendMetricItem(metricID, item); err != nil {
		log.Fatal("Error sending item:", err)
	}
}
