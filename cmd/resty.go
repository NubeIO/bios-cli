package commander

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

func (bt *BuildTool) handleRestyHTTPRequest(params interface{}) (interface{}, error) {
	// Convert params to a map[string]interface{}
	paramMap, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid params for Resty HTTP request")
	}

	// Create a new Resty client
	client := resty.New()

	// Extract the URL
	url, _ := paramMap["url"].(string)

	// Prepare the request
	req := client.R()

	// Set headers, if any
	if headers, ok := paramMap["header"].(map[string]interface{}); ok {
		for key, value := range headers {
			req.SetHeader(key, fmt.Sprintf("%v", value))
		}
	}

	// Set body, if any
	if body, ok := paramMap["body"].(map[string]interface{}); ok {
		req.SetBody(body)
	}

	// Set basic auth, if any
	if auth, ok := paramMap["auth"].(map[string]interface{}); ok {
		if basic, ok := auth["basic"].(map[string]interface{}); ok {
			username, _ := basic["username"].(string)
			password, _ := basic["password"].(string)
			req.SetBasicAuth(username, password)
		}
	}

	// Execute the request based on the method
	var resp *resty.Response
	var err error
	method, _ := paramMap["method"].(string)
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = req.Get(url)
	case "POST":
		resp, err = req.Post(url)
	case "PUT":
		resp, err = req.Put(url)
	case "DELETE":
		resp, err = req.Delete(url)
	case "PATCH":
		resp, err = req.Patch(url)
	case "HEAD":
		resp, err = req.Head(url)
	case "OPTIONS":
		resp, err = req.Options(url)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, fmt.Errorf("resty HTTP request failed: %v", err)
	}

	fmt.Printf("Response status code: %d\n", resp.StatusCode())
	fmt.Printf("Response body: %s\n", resp.String())

	return nil, nil
}
