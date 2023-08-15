package panobi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type apiURI string

const (
	// TODO CHANGE to prod
	TimeseriesURI apiURI = "https://dev.panobi.com/integrations/metrics-sdk/timeseries"
	ChartDataURI  apiURI = "https://dev.panobi.com/integrations/metrics-sdk/chart-data"
	DeleteURI     apiURI = "https://dev.panobi.com/integrations/metrics-sdk/delete"

	attempts          int = 3
	backoffInitial    int = 1
	backoffMultiplier int = 2
)

type transport struct {
	c  *http.Client
	ki KeyInfo
}

func createTransport(ki KeyInfo) *transport {
	return &transport{
		c:  &http.Client{},
		ki: ki,
	}
}

func (t *transport) post(uri apiURI, input []byte) ([]byte, error) {
	si, err := CalculateSignature(input, t.ki, nil)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/%s/%s",
		uri,
		url.PathEscape(t.ki.WorkspaceID),
		url.PathEscape(t.ki.ExternalID))
	backoff := backoffInitial
	i := 1

	for {
		b, err := func() ([]byte, error) {
			req, err := http.NewRequest("POST", url, bytes.NewReader(input))
			if err != nil {
				return nil, err
			}

			req.Header = t.getHeaders(si)

			resp, err := t.c.Do(req)
			if err != nil {
				fmt.Println("Returning early")
				return nil, err
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Fprintln(os.Stderr, "Error closing body:", err)
				}
			}()

			body, err := io.ReadAll(resp.Body)
			switch code := resp.StatusCode; {
			case code >= 200 && code < 300:
				return body, err
			case code == 408 || code == 429:
				if i == attempts {
					return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, body)
				}
				time.Sleep(getRetryAfter(resp, backoff))
				backoff = backoff * backoffMultiplier
				i++
				return nil, nil
			default:
				return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, body)
			}
		}()

		if err != nil {
			return nil, err
		}

		if b != nil {
			return b, nil
		}
	}
}

func (t *transport) getHeaders(si SignatureInfo) http.Header {
	headers := make(http.Header)

	headers.Set("Content-Type", "application/json")
	headers.Set("X-Panobi-Signature", si.S)
	headers.Set("X-Panobi-Request-Timestamp", si.TS)
	headers.Set("X-Request-ID", uuid.NewString())

	return headers
}

func getRetryAfter(resp *http.Response, defaultRetryAfter int) time.Duration {
	headerVal := strings.TrimSpace(resp.Header.Get("Retry-After"))
	if headerVal == "" {
		return time.Duration(defaultRetryAfter) * time.Second
	}

	retryAfter, err := strconv.ParseInt(headerVal, 10, 64)
	if err != nil {
		return time.Duration(defaultRetryAfter) * time.Second
	}

	return time.Duration(retryAfter) * time.Second
}
