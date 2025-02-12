package httpClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultContentType = "application/json"
)

// Client base http client struct
type Client struct {
	baseURL        *url.URL
	httpClient     *http.Client
	defaultHeaders map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
}

// NewClient creates new Client
func NewClient(baseURL string, defaultHeaders map[string]string) *Client {
	url, _ := url.Parse(baseURL)
	if defaultHeaders == nil {
		defaultHeaders = make(map[string]string)
	}
	return &Client{httpClient: &http.Client{}, baseURL: url, defaultHeaders: defaultHeaders}
}

// Request creates httpRequest and returns response
func (c *Client) Request(method string, resourceURI string, body interface{}, headers map[string]string) (Response, error) {
	url, _ := c.baseURL.Parse(resourceURI)

	responseBody := Response{}
	request, err := c.createNewRequest(method, url.String(), body, headers)

	if err != nil {
		return responseBody, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return responseBody, err
	}

	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseBody, err
	}

	if response.StatusCode > 399 {
		return responseBody, fmt.Errorf("Error: %v", string(b))
	}

	responseBody = Response{
		StatusCode: response.StatusCode,
		Body:       b,
	}
	return responseBody, nil
}

func (c *Client) createNewRequest(method string, url string, body interface{}, headers map[string]string) (*http.Request, error) {
	bodyBuffer, err := c.parseBody(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, err
	}
	requestHeaders := c.mergeWithDefaultHeaders(headers)
	c.addHeadersToRequest(request, requestHeaders)

	if body != nil {
		contentTypeHeader := c.getValueOfHeader(headers, "Content-Type")
		if contentTypeHeader == "" {
			request.Header.Set("Content-Type", defaultContentType)
		} else {
			request.Header.Set("Content-Type", contentTypeHeader)
		}
	}

	return request, nil
}

func (c *Client) addHeadersToRequest(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c *Client) parseBody(body interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter

	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func (c *Client) getValueOfHeader(headers map[string]string, header string) string {
	for currentHeader, currentValue := range headers {
		if strings.ToLower(currentHeader) == strings.ToLower(header) {
			return currentValue
		}
	}
	return ""
}

func (c *Client) mergeWithDefaultHeaders(incomingHeaders map[string]string) map[string]string {
	if incomingHeaders == nil {
		return c.defaultHeaders
	}
	mergedHeaders := make(map[string]string, len(incomingHeaders))
	for k, v := range c.defaultHeaders {
		mergedHeaders[k] = v
	}
	for k, v := range incomingHeaders {
		mergedHeaders[k] = v
	}
	return mergedHeaders
}
