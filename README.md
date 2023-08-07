# Panobi Metrics SDK

## Overview

This SDK lets you push metrics data from your private data source to your Panobi workspace, so you can observe growth patterns and make better-informed product decisions.

It also serves as an example for using the integration API in other languages.

## Who is it for?

If you use a private or in-house data warehouse system, this SDK will allow you to connect your data.

(If you use a third-party data warehouse system, please check out our other data integrations in your Panobi settings.)

## How does it work?

The SDK is based on metrics and items. Metrics are created in the Panobi UI and have a unique identifier, which is a string. Items are pairs of a calendar day and a numeric (float or integer) value. Metrics items can be sent one at a time or in batches of up to 1000 items. Panobi will only store new items.

## Compatibility

The [API specification](openapi.yaml) was generated for [OpenAPI 3.1.0](https://spec.openapis.org/oas/v3.1.0).

The source files were written against [Go 1.20](https://go.dev/doc/go1.20). They may also work with older versions.

## Getting started

The Metrics SDK must be enabled from your Settings -> Integrations page in your Panobi workspace. On this page you can copy your signing key, which you will need to authenticate this integration.

Publishing data for a metric requires the Metric ID. To obtain this ID, create a metric in the Metrics or Timeline page on your Panobi workspace and choose "Metrics SDK" as the data source. On this page you can copy your Metric ID for use with this SDK.

Once you have a signing key and at least one metric ID, you're ready to start running the provided [example programs](#running-the-example-programs), which demonstrate how to construct metrics items and send them to Panobi.

If you're using a language other than Golang, or you'd rather write your own commands, then take a look at how to send metrics to us via [OpenAPI](#openapi).

## Running the example programs

The example programs expect the signing key in the form of an environment variable.

```console
export METRICS_SDK_SIGNING_KEY=<your signing key>
```

Make sure to store your signing key in a secure location; do not commit it to source control.

### Simple

The simple example is a good place to start.

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
go run main.go ./metrics.csv
```

Each row is in the following format:

```
MetricID, Date, Value
```

The following are examples of valid rows:

```
<your metric id>,2023-08-01,1000
<your metric id>,2023-08-02,1000.5
```

### JSON

This example program works like the CSV example, but reads events in JSON format. Instead of one line per value, the JSON is structured to group items per metric ID.

```console
cd examples/json
go run main.go ./metrics.json
```

The following is an example of a valid row:

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
    https://panobi.com/integrations/metrics-sdk/items/"${wid}"/"${eid}"
```

## License

This SDK is provided under the terms of the [Apache License 2.0](LICENSE).

## About Panobi

Panobi is the platform for growth observability: helping companies see, understand, and drive their growth.
