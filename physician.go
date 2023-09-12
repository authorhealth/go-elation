package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type PhysicianServicer interface {
	Find(ctx context.Context, opts *FindPhysiciansOptions) ([]*Physician, *http.Response, error)
	Get(ctx context.Context, id int64) (*Physician, *http.Response, error)
}

var _ PhysicianServicer = (*PhysicianService)(nil)

type PhysicianService struct {
	client *Client
}

type Physician struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Npi          string `json:"npi"`
	License      string `json:"license"`
	LicenseState string `json:"license_state"`
	Credentials  string `json:"credentials"`
	Specialty    string `json:"specialty"`
	Email        string `json:"email"`
	UserID       int    `json:"user_id"`
	Practice     int    `json:"practice"`
	IsActive     bool   `json:"is_active"`
	Metadata     any    `json:"metadata"`
}

type FindPhysiciansOptions struct {
	FirstName string `url:"first_name,omitempty"`
	LastName  string `url:"last_name,omitempty"`
	NPI       string `url:"npi,omitempty"`
}

func (s *PhysicianService) Find(ctx context.Context, opts *FindPhysiciansOptions) ([]*Physician, *http.Response, error) {
	out := &Response[[]*Physician]{}

	res, err := s.client.request(ctx, http.MethodGet, "/physicians", opts, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out.Results, res, nil
}

func (s *PhysicianService) Get(ctx context.Context, id int64) (*Physician, *http.Response, error) {
	out := &Physician{}

	res, err := s.client.request(ctx, http.MethodGet, "/physicians/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
