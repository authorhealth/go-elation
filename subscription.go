package elation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type SubscriptionServicer interface {
	Find(ctx context.Context) ([]*Subscription, *http.Response, error)
	Subscribe(ctx context.Context, opts *SubscribeOptions) (*Subscription, *http.Response, error)
	Delete(ctx context.Context, opts *DeleteSubscriptionOptions) (*http.Response, error)
}

var _ SubscriptionServicer = (*SubscriptionService)(nil)

type SubscriptionService struct {
	client *Client
}

type Subscription struct {
	ID            int64                 `json:"id"`
	Resource      string                `json:"resource"`
	Target        string                `json:"target"`
	CreatedDate   SubscriptionJSONDate  `json:"created_date"`
	DeletedDate   *SubscriptionJSONDate `json:"deleted_date"`
	SigningPubKey string                `json:"signing_pub_key"`
}

type SubscriptionJSONDate time.Time

func (s *SubscriptionJSONDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02T15:04:05", strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	*s = SubscriptionJSONDate(t)

	return nil
}

func (s *SubscriptionService) Find(ctx context.Context) ([]*Subscription, *http.Response, error) {
	out := &Response[[]*Subscription]{}

	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodGet, "/app/subscriptions/", nil, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out.Results, res, nil
}

type SubscribeOptions struct {
	Resource   string          `json:"resource"`
	Target     string          `json:"target"`
	Properties json.RawMessage `json:"properties"`
}

func (s *SubscriptionService) Subscribe(ctx context.Context, opts *SubscribeOptions) (*Subscription, *http.Response, error) {
	out := &Subscription{}

	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodPost, "/app/subscriptions/", nil, opts, out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

type DeleteSubscriptionOptions struct {
	ID int64 `url:"id"`
}

func (s *SubscriptionService) Delete(ctx context.Context, opts *DeleteSubscriptionOptions) (*http.Response, error) {
	// The trailing slash in the path is required.
	res, err := s.client.request(ctx, http.MethodDelete, "/app/subscriptions/"+strconv.Itoa(int(opts.ID))+"/", nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}
