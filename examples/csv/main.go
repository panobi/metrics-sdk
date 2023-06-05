package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/civil"
	panobi "github.com/panobi/metrics-sdk"
)

func main() {
	//
	// We need the name of the CSV file.
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
	// Open the file and read it. Each line is assumed to be a list of
	// comma-separated values, with the columns in the following order.
	//
	// metricID, date, value
	//

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	items := make(map[string][]panobi.MetricItem, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Split(line, ",")

		if len(cols) != 3 {
			log.Fatalf(`Expected three columns in "%s"`, line)
		}

		metricID := strings.TrimSpace(cols[0])
		date, err := civil.ParseDate(strings.TrimSpace(cols[1]))
		if err != nil {
			log.Fatalf("Unable to parse date value %s: %s", cols[1], err.Error())
		}

		// value may be int or float
		value, err := strconv.ParseFloat(strings.TrimSpace(cols[2]), 64)
		if err != nil {
			log.Fatalf("Unable to parse metric value %s: %s", cols[2], err.Error())
		}

		item := panobi.MetricItem{
			Date:  date,
			Value: value,
		}

		_, ok := items[metricID]
		if ok {
			items[metricID] = append(items[metricID], item)
			// when we reach max batch size, send the items to Panobi and then start a new batch
			if len(items[metricID]) == panobi.MaxItems {
				err := client.SendMetricItems(metricID, items[metricID])
				if err != nil {
					log.Fatalf("Error sending items for metricID %s: %s", metricID, err.Error())
				}

				log.Printf("Successfully sent %d event(s) for metricID %s", panobi.MaxItems, metricID)
				items[metricID] = make([]panobi.MetricItem, 0)
			}
		} else {
			items[metricID] = []panobi.MetricItem{item}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print("Error reading input:", err)
	}

	//
	// Send any remaining items to Panobi.
	//

	for metricID, i := range items {
		if len(i) > 0 {
			err := client.SendMetricItems(metricID, i)
			if err != nil {
				log.Fatalf("Error sending items for metricID %s: %s", metricID, err.Error())
			}

			log.Printf("Successfully sent %d item(s) for metricID %s", len(i), metricID)
		}
	}
}
