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
	// Open the file and read it.
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
	scanner := bufio.NewScanner(file)

	if timeseries {
		sendTimeseriesData(scanner, client)
	} else {
		sendChartData(scanner, client)
	}
}

func sendTimeseriesData(scanner *bufio.Scanner, client *panobi.Client) {
	items := make(map[string][]panobi.MetricItem, 0)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		i++
		// skip header row if present
		if i == 1 && strings.HasPrefix(line, "MetricID") {
			continue
		}

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

func sendChartData(scanner *bufio.Scanner, client *panobi.Client) {
	items := make(map[string][]panobi.ChartData, 0)

	i := 0
	columns := make(map[int]string, 0)
	clearedMetricIDs := make(map[string]bool, 0)
	metricIDColumn := -1
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ",")

		i++
		if i == 1 {
			for index, col := range row {
				columns[index] = strings.TrimSpace(col)
				if columns[index] == "MetricID" {
					metricIDColumn = index
				}
			}
			if metricIDColumn == -1 {
				log.Fatalf(`Expected a column named MetricID in header row "%s"`, line)
			}

			if len(columns) < 2 {
				log.Fatalf(`Expected at least two columns in header row "%s"`, line)
			}
			continue
		}

		if len(row) != len(columns) {
			log.Fatalf(`Column count %d does not match header row %d in "%s"`, len(row), len(columns), line)
		}

		item := panobi.ChartData{}
		var metricID string
		for index, name := range columns {
			var value interface{}
			value = strings.TrimSpace(row[index])
			if name == "MetricID" {
				metricID = value.(string)
			} else {
				// if the value successfully parses as an int or a float, we'll send it as such
				// otherwise it's sent as a string
				numericValue, err := strconv.ParseFloat(value.(string), 64)
				if err == nil {
					value = numericValue
				}
				item[name] = value
			}
		}
		_, ok := items[metricID]
		if ok {
			items[metricID] = append(items[metricID], item)

			// when we reach max batch size, send the items to Panobi and then start a new batch
			if len(items[metricID]) == panobi.MaxItems {
				// before sending the first batch of data for a metric, delete existing chart data
				if !clearedMetricIDs[metricID] {
					err := client.DeleteMetricData(metricID)
					if err != nil {
						log.Fatalf("Error deleting existing data for metricID %s: %s", metricID, err.Error())
					}
					clearedMetricIDs[metricID] = true
				}
				err := client.SendMetricChartData(metricID, items[metricID])
				if err != nil {
					log.Fatalf("Error sending items for metricID %s: %s", metricID, err.Error())
				}

				log.Printf("Successfully sent %d event(s) for metricID %s", panobi.MaxItems, metricID)
				items[metricID] = make([]panobi.ChartData, 0)
			}
		} else {
			items[metricID] = []panobi.ChartData{item}
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
			// before sending the first batch of data for a metric, delete existing chart data
			if !clearedMetricIDs[metricID] {
				err := client.DeleteMetricData(metricID)
				if err != nil {
					log.Fatalf("Error deleting existing data for metricID %s: %s", metricID, err.Error())
				}
				clearedMetricIDs[metricID] = true
			}

			err := client.SendMetricChartData(metricID, i)
			if err != nil {
				log.Fatalf("Error sending items for metricID %s: %s", metricID, err.Error())
			}

			log.Printf("Successfully sent %d item(s) for metricID %s", len(i), metricID)
		}
	}
}
