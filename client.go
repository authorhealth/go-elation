package elation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	querystring "github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	httpClient *http.Client
	baseURL    string

	AppointmentSvc     *AppointmentService
	ServiceLocationSvc *ServiceLocationService
	SubscriptionSvc    *SubscriptionService
}

func NewClient(httpClient *http.Client, tokenURL, clientID, clientSecret, baseURL string) *Client {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	client := &Client{
		httpClient: config.Client(ctx),
		baseURL:    baseURL,
	}

	client.AppointmentSvc = &AppointmentService{client}
	client.ServiceLocationSvc = &ServiceLocationService{client}
	client.SubscriptionSvc = &SubscriptionService{client}

	return client
}

type Response[ResultsT any] struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  ResultsT `json:"results"`
}

type ErrorResponse struct {
	StatusCode int
	Detail     string `json:"detail"`
}

func (e *ErrorResponse) Error() string {
	return e.Detail
}

func (c *Client) request(ctx context.Context, method string, path string, query any, body any, out any) (*http.Response, error) {
	q, err := querystring.Values(query)
	if err != nil {
		return nil, fmt.Errorf("encoding URL query: %w", err)
	}

	u := c.baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	reader := bytes.NewReader(nil)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}

		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return nil, fmt.Errorf("making new HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing HTTP request: %w", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response bodu: %w", err)
	}
	//nolint
	_ = res.Body.Close()

	res.Body = io.NopCloser(bytes.NewBuffer(resBody))

	if res.StatusCode > http.StatusIMUsed {
		errRes := &ErrorResponse{}
		err = json.Unmarshal(resBody, errRes)
		if err != nil {
			return res, fmt.Errorf("unmarshaling response body (error): %w", err)
		}

		errRes.StatusCode = res.StatusCode

		return res, fmt.Errorf("API error: %w", errRes)
	}

	if out != nil {
		err = json.Unmarshal(resBody, out)
		if err != nil {
			return res, fmt.Errorf("unmarshaling results: %w", err)
		}
	}

	return res, nil
}
