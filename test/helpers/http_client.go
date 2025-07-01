package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type HTTPClient struct {
	baseURL string
	client  *http.Client
	t       *testing.T
}

func NewTestHTTPClient(t *testing.T, baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		client:  &http.Client{},
		t:       t,
	}
}

type Request struct {
	Method  string
	Path    string
	Body    any
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func (c *HTTPClient) Do(req Request) *Response {
	var bodyReader io.Reader

	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		require.NoError(c.t, err, "Failed to marshal request body")
		bodyReader = bytes.NewReader(bodyBytes)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, req.Path)
	httpReq, err := http.NewRequest(req.Method, url, bodyReader)
	require.NoError(c.t, err, "Failed to create HTTP request")

	if req.Body != nil && (req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH") {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := c.client.Do(httpReq)
	require.NoError(c.t, err, "Failed to execute HTTP request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(c.t, err, "Failed to read response body")

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}
}

func (c *HTTPClient) GET(path string, headers ...map[string]string) *Response {
	req := Request{
		Method: "GET",
		Path:   path,
	}
	if len(headers) > 0 {
		req.Headers = headers[0]
	}
	return c.Do(req)
}

func (c *HTTPClient) POST(path string, body any, headers ...map[string]string) *Response {
	req := Request{
		Method: "POST",
		Path:   path,
		Body:   body,
	}
	if len(headers) > 0 {
		req.Headers = headers[0]
	}
	return c.Do(req)
}

func (c *HTTPClient) PUT(path string, body any, headers ...map[string]string) *Response {
	req := Request{
		Method: "PUT",
		Path:   path,
		Body:   body,
	}
	if len(headers) > 0 {
		req.Headers = headers[0]
	}
	return c.Do(req)
}

func (c *HTTPClient) DELETE(path string, headers ...map[string]string) *Response {
	req := Request{
		Method: "DELETE",
		Path:   path,
	}
	if len(headers) > 0 {
		req.Headers = headers[0]
	}
	return c.Do(req)
}

func (r *Response) ParseJSON(target any) error {
	return json.Unmarshal(r.Body, target)
}

func (r *Response) GetBodyAsString() string {
	return string(r.Body)
}
