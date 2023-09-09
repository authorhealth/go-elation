package elation

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ServiceLocationServicer interface {
	Find(ctx context.Context) ([]*ServiceLocation, *http.Response, error)
}

var _ ServiceLocationServicer = (*ServiceLocationService)(nil)

type ServiceLocationService struct {
	client *Client
}

type ServiceLocation struct {
	AddressLine1   string     `json:"address_line1"`
	AddressLine2   string     `json:"address_line2"`
	City           string     `json:"city"`
	CreatedDate    time.Time  `json:"created_date"`
	DeletedDate    *time.Time `json:"deleted_date"`
	Email          string     `json:"email"`
	Fax            string     `json:"fax"`
	ID             int64      `json:"id"`
	IsPrimary      bool       `json:"is_primary"`
	Name           string     `json:"name"`
	Phone          string     `json:"phone"`
	PlaceOfService string     `json:"place_of_service"`
	Practice       int64      `json:"practice"`
	State          string     `json:"state"`
	Status         string     `json:"status"`
	Zip            string     `json:"zip"`
}

func (s *ServiceLocationService) Find(ctx context.Context) ([]*ServiceLocation, *http.Response, error) {
	out := &Response[[]*ServiceLocation]{}

	res, err := s.client.request(ctx, http.MethodGet, "/service_locations", nil, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out.Results, res, nil
}
