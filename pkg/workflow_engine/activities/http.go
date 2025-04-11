package activities

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"strconv"
	"time"

	workflowengine "github.com/forkbombeu/didimo/pkg/workflow_engine"
)

type HTTPActivity struct{}

func (HTTPActivity) Name() string {
	return "HTTPRequest"
}

func (a *HTTPActivity) Execute(ctx context.Context, input workflowengine.ActivityInput) (workflowengine.ActivityResult, error) {
	var result workflowengine.ActivityResult

	method := input.Config["method"]
	url := input.Config["url"]
	if queryParams, ok := input.Payload["query_params"].(map[string]any); ok {
		parsedURL, err := URL.Parse(url)
		if err != nil {
			return workflowengine.Fail(&result, fmt.Sprintf("failed to parse URL: %v", err))
		}

		// Add query parameters
		query := parsedURL.Query()
		for key, value := range queryParams {
			if strValue, ok := value.(string); ok {
				query.Add(key, strValue)
			}
		}

		parsedURL.RawQuery = query.Encode()
		url = parsedURL.String() // Update the URL with query parameters
	}
	if method == "" || url == "" {
		return workflowengine.Fail(&result, "missing 'method' or 'url' in config")
	}
	timeout := 10 * time.Second
	if tStr, ok := input.Config["timeout"]; ok {
		if t, err := strconv.Atoi(tStr); err == nil {
			timeout = time.Duration(t) * time.Second
		}
	}
	headers := map[string]string{}
	if rawHeaders, ok := input.Payload["headers"].(map[string]any); ok {
		for k, v := range rawHeaders {
			if vs, ok := v.(string); ok {
				headers[k] = vs
			}
		}
	}
	var body io.Reader
	if input.Payload["body"] != nil {
		jsonBody, err := json.Marshal(input.Payload["body"])
		if err != nil {
			return workflowengine.Fail(&result, fmt.Sprintf("failed to marshal body: %v", err))
		}
		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return workflowengine.Fail(&result, fmt.Sprintf("failed to create request: %v", err))
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return workflowengine.Fail(&result, fmt.Sprintf("failed to perform request: %v", err))
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return workflowengine.Fail(&result, fmt.Sprintf("failed to read response body: %v", err))
	}

	// Handle non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return workflowengine.Fail(&result,
			fmt.Sprintf("received status code %d: %s", resp.StatusCode, string(respBody)),
		)
	}

	var output any
	if err := json.Unmarshal(respBody, &output); err != nil {
		// if not JSON, return as string
		output = string(respBody)
	}

	result.Output = map[string]any{
		"status":  resp.StatusCode,
		"headers": resp.Header,
		"body":    output,
	}

	return result, nil
}
