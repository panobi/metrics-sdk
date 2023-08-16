# Panobi Metrics SDK

## Overview

This SDK lets you push metrics data from your private data source to your Panobi workspace, so you can observe growth patterns and make better-informed product decisions.

It also serves as an example for using the integration API in other languages.

## Who is it for?

If you use a private or in-house data warehouse system, this SDK will allow you to connect your data.

(If you use a third-party data warehouse system, please check out our other data integrations in your Panobi settings.)

## How does it work?

The SDK is based on metrics and items. Metrics are created in the Panobi UI and have a unique identifier, which is a string.

There are two kinds of metrics in Panobi. **Timeseries** metrics show on the Panobi Timeline page and require a calendar day as the X-axis, along with a single numeric (float or integer) value. The day is effectively a unique key for a metric. Timeseries data can be sent one item at a time or in batches of up to 1000 items. Panobi will only store new items.

Other chart types like bar, column, area, and table support arbitrary numbers of columns of different types.

This SDK uses separate API endpoints to send data for timeseries metrics and other chart types, so you'll need to know which kind of metric you're sending.

## How to use this SDK

There are three main ways to use the SDK

1. Use the Go libraries in the SDK to write a program in Go to send data to Panobi
2. Export data into CSV or JSON formats using the tools of your choice (such as SQL exports, spreadsheets, or any programming language), and use Go to compile and run our example programs to send the data to Panobi
3. Use the [curl example](#openapi) or construct a similar request in any programming language, without needing to use Go

## Compatibility

The [API specification](openapi.yaml) was generated for [OpenAPI 3.1.0](https://spec.openapis.org/oas/v3.1.0).

The source files were written against [Go 1.20](https://go.dev/doc/go1.20). They may also work with older versions.

## Getting started

The Metrics SDK must be enabled from your Settings -> Integrations page in your Panobi workspace. On this page you can copy your signing key, which you will need to authenticate this integration.

Publishing data for a metric requires the Metric ID. To obtain this ID, create a metric on the Metrics or Timeline page on your Panobi workspace and choose "Metrics SDK" as the data source. On this page you can copy your Metric ID for use with this SDK. Your selection for the "Time series" toggle when creating the metric will determine how you should send the data using the SDK - timeseries metrics are sent differently from other types of metrics.

Once you have a signing key and at least one metric ID, you're ready to start running the provided [example programs](#running-the-example-programs), which demonstrate how to construct metrics items and send them to Panobi.

If you're using a language other than Golang, or you'd rather write your own commands, then take a look at how to send metrics to us via [OpenAPI](#openapi).

## Running the example programs

The example programs expect the signing key in the form of an environment variable.

```console
export METRICS_SDK_SIGNING_KEY=<your signing key>
```

Make sure to store your signing key in a secure location; do not commit it to source control.

### Simple

The simple example uses hard coded data to demonstrate how to write a Go program to send data for a timeseries metric to Panobi.

```console
cd examples/simple
go run main.go
```

Roughly, it works as follows.

1. Reads the environment variable and parses your key.
2. Creates a client with the parsed key.
3. Constructs a time series item for a hard-coded metric ID.
4. Sends the item to Panobi.

Once the item has been successfully sent, it should show up in your timeline in Panobi.

### CSV

This example program demonstrates how to send more than one item at a time. It will read a file of comma-separated values, where each line represents one item for one metric. If you're able to export the results of a data warehouse query to a CSV you may be able to use this program as-is to upload your data.

```console
cd examples/csv

# for timeseries metrics
go run main.go -t ./metrics.csv

# for other chart types
go run main.go ./metrics.csv
```

Each row is in the following format:

```
MetricID,Date,Value
```

The following are examples of valid rows:

```
<your metric id>,2023-08-01,1000
<your metric id>,2023-08-02,1000.5
```

For timeseries metrics, the header row is optional and must be `MetricID,Date,Value` if present. Only new rows will be uploaded, existing rows will not be modified.

For other chart types, a header row is required to set the column names. `MetricID` must be one of the columns. Before any data is sent, all existing data will be deleted.

### JSON

This example program works like the CSV example, but reads events in JSON format. Instead of one line per value, the JSON is structured to group items per metric ID.

```console
cd examples/json

# for timeseries metrics
go run main.go -t ./metrics.json

# for other chart types
go run main.go ./metrics.json
```

The following is an example of a valid row for a timeseries metric.

```json
[
  {
    "metricID": "<your metric id>",
    "items": [{ "date": "2023-08-01", "value": 1000 }]
  }
]
```

## OpenAPI

In an effort to be language agnostic, we've provided an [OpenAPI specification](openapi.yaml) that you can use to send data directly to Panobi.

Once you've built a request according to the specification, you need to sign it so that Panobi knows it's from you. The following little shell script demonstrates how to do this via curl.
All modern programming languages should have equivalent libraries allowing you to sign an hmac payload using your signing key in a similar fashion.

```shell
#!/usr/bin/env bash

# We read the request body from a file, so that we can generate a signature for
# it. The first argument to this script is the name of that file.
input=$(<"$1")

# Split the signing key into its component parts.
arr=(${METRICS_SDK_SIGNING_KEY//-/ })
wid=${arr[0]}    # Workspace ID
eid=${arr[1]}    # External ID
secret=${arr[2]} # Secret

# Get the milliseconds since Unix epoch. We'll use this as a timestamp for the
# signature.
ts=$(date +%s)000

# Hash the timestamp and the request body using the secret part of the signing
# key.
msg="v0:${ts}:${input}"
sig=$(echo -n "${msg}" | openssl dgst -sha256 -hmac "${secret}")

# Post the headers and the request to Panobi using good ol' curl.
curl -v \
    -X POST \
    -H "X-Panobi-Signature: v0=""${sig}" \
    -H "X-Panobi-Request-Timestamp: ""${ts}" \
    -H "Content-Type: application/json" \
    -d "${input}" \
    https://panobi.com/integrations/metrics-sdk/timeseries/"${wid}"/"${eid}"
```

## License

This SDK is provided under the terms of the [Apache License 2.0](LICENSE).

## About Panobi

Panobi is the platform for growth observability: helping companies see, understand, and drive their growth.
