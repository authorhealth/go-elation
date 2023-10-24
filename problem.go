package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ProblemServicer interface {
	Find(ctx context.Context, opts *FindPatientProblemsOptions) (*Response[[]*PatientProblem], *http.Response, error)
	Get(ctx context.Context, id int64) (*PatientProblem, *http.Response, error)
}

var _ ProblemServicer = (*ProblemService)(nil)

type ProblemService struct {
	client *Client
}

type PatientProblem struct {
	ID           int64               `json:"id"`
	Description  string              `json:"description"`
	Status       string              `json:"status"`
	Synopsis     string              `json:"synopsis"`
	StartDate    string              `json:"start_date"`
	ResolvedDate string              `json:"resolved_date"`
	Dx           []*PatientProblemDX `json:"dx"`
	Patient      int64               `json:"patient"`
	CreatedDate  time.Time           `json:"created_date"`
	DeletedDate  *time.Time          `json:"deleted_date"`
}

type PatientProblemDX struct {
	Icd9   []string `json:"icd9"`
	Icd10  []string `json:"icd10"`
	Snomed string   `json:"snomed"`
}

type FindPatientProblemsOptions struct {
	*Pagination

	Patient int64 `url:"patient,omitempty"`
}

func (s *ProblemService) Find(ctx context.Context, opts *FindPatientProblemsOptions) (*Response[[]*PatientProblem], *http.Response, error) {
	out := &Response[[]*PatientProblem]{}

	res, err := s.client.request(ctx, http.MethodGet, "/problems", opts, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *ProblemService) Get(ctx context.Context, id int64) (*PatientProblem, *http.Response, error) {
	out := &PatientProblem{}

	res, err := s.client.request(ctx, http.MethodGet, "/problems/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
