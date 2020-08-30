package alphavantage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

const aplhavantageURL = "https://www.alphavantage.co"

// HTTPClient allows to plug a custom http.Client as long as it can Do(http.Request)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func buildURL(apiKey string, function string, symbol string) string {
	return fmt.Sprintf("%s/query?function=%s&symbol=%s&apikey=%s", aplhavantageURL, function, symbol, apiKey)
}

func panicParseInt64ish(v string) int64 {
	if v == "None" {
		return 0
	}
	res, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		panic(errors.Wrapf(err, "Cannot parse '%s'", v))
	}
	return res
}

func makeRequest(ctx context.Context, httpClient HTTPClient, url string, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "Error creating http.Request")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Error during HTTP call")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		_, err = io.Copy(&buf, res.Body)
		if err != nil {
			return errors.Wrap(err, "Error reading result.Body")
		}

		log.Printf("[DEBUG] %s\n", buf.String())
		return errors.Errorf("Alphavantage HTTP Status %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}
