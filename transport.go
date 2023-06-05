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
)

const (
	// TODO CHANGE
	itemsURI string = "http://localhost:8080/integrations/metrics-sdk/items"

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

func (t *transport) post(input []byte) ([]byte, error) {
	si, err := CalculateSignature(input, t.ki, nil)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/%s/%s",
		itemsURI,
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
				return nil, err
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					fmt.Fprintln(os.Stderr, "Error closing body:", err)
				}
			}()

			switch code := resp.StatusCode; {
			case code >= 200 && code < 300:
				return io.ReadAll(resp.Body)
			case code == 408 || code == 429:
				if i == attempts {
					return nil, fmt.Errorf("http error: %d", resp.StatusCode)
				}
				time.Sleep(getRetryAfter(resp, backoff))
				backoff = backoff * backoffMultiplier
				i++
				return nil, nil
			default:
				return nil, fmt.Errorf("http error: %d", resp.StatusCode)
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
