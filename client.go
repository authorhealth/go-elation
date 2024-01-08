package elation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	querystring "github.com/google/go-querystring/query"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultPaginationLimit = 25

	ResourceAppointments string = "appointments"
	ResourceMedications  string = "medications"
	ResourcePatients     string = "patients"
	ResourcePhysicians   string = "physicians"
	ResourceProblems     string = "problems"

	WebhookEventActionSaved   string = "saved"
	WebhookEventActionDeleted string = "deleted"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	tracer     trace.Tracer

	AppointmentSvc          *AppointmentService
	ClinicalDocumentSvc     *ClinicalDocumentService
	ContactSvc              *ContactService
	InsuranceCompanySvc     *InsuranceCompanyService
	InsuranceEligibilitySvc *InsuranceEligibilityService
	InsurancePlanSvc        *InsurancePlanService
	LetterSvc               *LetterService
	MedicationSvc           *MedicationService
	NonVisitNoteSvc         *NonVisitNoteService
	PatientSvc              *PatientService
	PhysicianSvc            *PhysicianService
	PracticeSvc             *PracticeService
	ProblemSvc              *ProblemService
	ServiceLocationSvc      *ServiceLocationService
	SubscriptionSvc         *SubscriptionService
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
		tracer:     otel.GetTracerProvider().Tracer("github.com/authorhealth/go-elation"),
	}

	client.AppointmentSvc = &AppointmentService{client}
	client.ClinicalDocumentSvc = &ClinicalDocumentService{client}
	client.ContactSvc = &ContactService{client}
	client.InsuranceCompanySvc = &InsuranceCompanyService{client}
	client.InsuranceEligibilitySvc = &InsuranceEligibilityService{client}
	client.InsurancePlanSvc = &InsurancePlanService{client}
	client.LetterSvc = &LetterService{client}
	client.MedicationSvc = &MedicationService{client}
	client.NonVisitNoteSvc = &NonVisitNoteService{client}
	client.PatientSvc = &PatientService{client}
	client.PhysicianSvc = &PhysicianService{client}
	client.PracticeSvc = &PracticeService{client}
	client.ProblemSvc = &ProblemService{client}
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

func (r *Response[ResultsT]) HasPrevious() bool {
	return len(r.Previous) > 0
}

func (r *Response[ResultsT]) HasNext() bool {
	return len(r.Next) > 0
}

func (r *Response[ResultsT]) PaginationNext() *Pagination {
	return &Pagination{
		Limit:  parseLimit(r.Next),
		Offset: parseOffset(r.Next),
	}
}

func (r *Response[ResultsT]) PaginationPrevious() *Pagination {
	return &Pagination{
		Limit:  parseLimit(r.Previous),
		Offset: parseOffset(r.Previous),
	}
}

func (r *Response[ResultsT]) PaginationNextWithLimit(limit int) *Pagination {
	return &Pagination{
		Limit:  limit,
		Offset: parseOffset(r.Next),
	}
}

type Error struct {
	StatusCode int
	Body       string
}

type Pagination struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("API error (status code %d)", e.StatusCode)
}

func (c *Client) request(ctx context.Context, method string, path string, query any, body any, out any) (*http.Response, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(semconv.HTTPRequestMethodKey.String(method))

	q, err := querystring.Values(query)
	if err != nil {
		return nil, fmt.Errorf("encoding URL query: %w", err)
	}

	u := c.baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	span.SetAttributes(semconv.URLFull(u))

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

	span.SetAttributes(semconv.HTTPResponseStatusCode(res.StatusCode))

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	//nolint
	_ = res.Body.Close()

	res.Body = io.NopCloser(bytes.NewBuffer(resBody))

	if res.StatusCode >= http.StatusBadRequest {
		return res, &Error{
			StatusCode: res.StatusCode,
			Body:       string(resBody),
		}
	}

	if out != nil {
		err = json.Unmarshal(resBody, out)
		if err != nil {
			return res, fmt.Errorf("unmarshaling results: %w", err)
		}
	}

	return res, nil
}

func parseLimit(v string) int {
	u, err := url.Parse(v)
	if err != nil {
		return defaultPaginationLimit
	}

	offset, err := strconv.Atoi(u.Query().Get("limit"))
	if err != nil {
		return defaultPaginationLimit
	}

	return offset
}

func parseOffset(v string) int {
	u, err := url.Parse(v)
	if err != nil {
		return 0
	}

	offset, err := strconv.Atoi(u.Query().Get("offset"))
	if err != nil {
		return 0
	}

	return offset
}
