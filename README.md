# Panobi Metrics SDK

## Overview

This SDK lets you push metrics data to your Panobi workspace. It also serves as an example for using the integration API in other languages.

## Who is it for?

Anyone with a private or in-house data warehouse system who cannot use our other metrics data source integrations.

## How does it work?

The SDK is based on metrics and items. Metrics are created in the Panobi UI and have a unique intendifier, which is a string. Items are pairs of a calendar day and a numeric (float or integer) value. Metrics items can be sent one at a time or in batches of up to 1000 items. Panobi will only store new items.

## Compatibility

The [API specification](openapi.yaml) was generated for [OpenAPI 3.1.0](https://spec.openapis.org/oas/v3.1.0).

The source files were written against [Go 1.20](https://go.dev/doc/go1.20). They may also work with older versions.

## Getting started

The quickest way to get started is to run the provided [example programs](#running-the-example-programs), which demonstrate how to construct metrics items and send them to Panobi.

If you're using a language other than Golang, or you'd rather write your own commands, then take a look at how to send metrics to us via [OpenAPI](#openapi).

## Running the example programs

You will need your signing key, which you can copy from the integration settings in your Panobi workspace. The example programs expect the signing key in the form of an environment variable.

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

1. Reads the enviroment variable and parses your key.
2. Creates a client with the parsed key.
3. Constructs a timeseries item for a hard coded metric ID
4. Sends the item to Panobi.

Once the item has been successfully sent, it should show up in the Timeline view of your Panobi workspace.

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
XRnrRBTedmWzy8RQ6pqh2d,2023-07-01,1000
XRnrRBTedmWzy8RQ6pqh2d,2023-07-02,1000.5
```

### JSON

This example program works like the CSV example, but reads events in JSON format. Instead of one line per value, the JSON is structured to group items per metricID.

```console
cd examples/json
go run main.go ./metrics.json
```

The following is an example of a valid row:

```json
[
    {
        "metricID": "XRnrRBTedmWzy8RQ6pqh2d",
        "items": [
            {"date": "2023-07-01", "value": 1000}
        ]
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

The platform designed for growth teams.

Panobi helps growth teams increase their velocity, deliver results, and amplify customer insights across the company.
