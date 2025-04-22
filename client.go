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
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultPaginationLimit = 25
)

type Client interface {
	Allergies() AllergyServicer
	Appointments() AppointmentServicer
	Bill() BillServicer
	ClinicalDocuments() ClinicalDocumentServicer
	Contacts() ContactServicer
	DiscontinuedMedications() DiscontinuedMedicationServicer
	HistoryDownloadFills() HistoryDownloadFillServicer
	InsuranceCompanies() InsuranceCompanyServicer
	InsuranceEligibility() InsuranceEligibilityServicer
	InsurancePlans() InsurancePlanServicer
	Letters() LetterServicer
	Medications() MedicationServicer
	MessageThreads() MessageThreadServicer
	NonVisitNotes() NonVisitNoteServicer
	Patients() PatientServicer
	Pharmacies() PharmacyServicer
	Physicians() PhysicianServicer
	Practices() PracticeServicer
	PrescriptionFills() PrescriptionFillServicer
	Problems() ProblemServicer
	RecurringEventGroups() RecurringEventGroupServicer
	ServiceLocations() ServiceLocationServicer
	Subscriptions() SubscriptionServicer
	ThreadMembers() ThreadMemberServicer
	VisitNote() VisitNoteServicer
}

type HTTPClient struct {
	httpClient *http.Client
	baseURL    string
	tracer     trace.Tracer

	AllergySvc                 *AllergyService
	AppointmentSvc             *AppointmentService
	BillSvc                    *BillService
	ClinicalDocumentSvc        *ClinicalDocumentService
	ContactSvc                 *ContactService
	DiscontinuedMedicationSvc  *DiscontinuedMedicationService
	HistoryDownloadFillSvc     *HistoryDownloadFillService
	InsuranceCompanySvc        *InsuranceCompanyService
	InsuranceEligibilitySvc    *InsuranceEligibilityService
	InsurancePlanSvc           *InsurancePlanService
	LetterSvc                  *LetterService
	MedicationSvc              *MedicationService
	MessageThreadSvc           *MessageThreadService
	NonVisitNoteSvc            *NonVisitNoteService
	PatientSvc                 *PatientService
	PharmacySvc                *PharmacyService
	PhysicianSvc               *PhysicianService
	PracticeSvc                *PracticeService
	PrescriptionFillSvc        *PrescriptionFillService
	ProblemSvc                 *ProblemService
	RecurringEventGroupService *RecurringEventGroupService
	ServiceLocationSvc         *ServiceLocationService
	SubscriptionSvc            *SubscriptionService
	ThreadMemberSvc            *ThreadMemberService
	VisitNoteSvc               *VisitNoteService
}

var _ Client = (*HTTPClient)(nil)

func NewHTTPClient(httpClient *http.Client, tokenURL, clientID, clientSecret, baseURL string) *HTTPClient {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport,
		// Ensure that the trace context is not propagated by using an empty composite propagator.
		otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator()))

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	client := &HTTPClient{
		httpClient: config.Client(ctx),
		baseURL:    baseURL,
		tracer:     otel.GetTracerProvider().Tracer("github.com/authorhealth/go-elation"),
	}

	client.AllergySvc = &AllergyService{client}
	client.AppointmentSvc = &AppointmentService{client}
	client.BillSvc = &BillService{client}
	client.ClinicalDocumentSvc = &ClinicalDocumentService{client}
	client.ContactSvc = &ContactService{client}
	client.DiscontinuedMedicationSvc = &DiscontinuedMedicationService{client}
	client.HistoryDownloadFillSvc = &HistoryDownloadFillService{client}
	client.InsuranceCompanySvc = &InsuranceCompanyService{client}
	client.InsuranceEligibilitySvc = &InsuranceEligibilityService{client}
	client.InsurancePlanSvc = &InsurancePlanService{client}
	client.LetterSvc = &LetterService{client}
	client.MedicationSvc = &MedicationService{client}
	client.MessageThreadSvc = &MessageThreadService{client}
	client.NonVisitNoteSvc = &NonVisitNoteService{client}
	client.PatientSvc = &PatientService{client}
	client.PharmacySvc = &PharmacyService{client}
	client.PhysicianSvc = &PhysicianService{client}
	client.PracticeSvc = &PracticeService{client}
	client.PrescriptionFillSvc = &PrescriptionFillService{client}
	client.ProblemSvc = &ProblemService{client}
	client.RecurringEventGroupService = &RecurringEventGroupService{client}
	client.ServiceLocationSvc = &ServiceLocationService{client}
	client.SubscriptionSvc = &SubscriptionService{client}
	client.ThreadMemberSvc = &ThreadMemberService{client}
	client.VisitNoteSvc = &VisitNoteService{client}

	return client
}

func (c *HTTPClient) Allergies() AllergyServicer {
	return c.AllergySvc
}

func (c *HTTPClient) Appointments() AppointmentServicer {
	return c.AppointmentSvc
}

func (c *HTTPClient) Bill() BillServicer {
	return c.BillSvc
}

func (c *HTTPClient) ClinicalDocuments() ClinicalDocumentServicer {
	return c.ClinicalDocumentSvc
}

func (c *HTTPClient) Contacts() ContactServicer {
	return c.ContactSvc
}

func (c *HTTPClient) DiscontinuedMedications() DiscontinuedMedicationServicer {
	return c.DiscontinuedMedicationSvc
}

func (c *HTTPClient) HistoryDownloadFills() HistoryDownloadFillServicer {
	return c.HistoryDownloadFillSvc
}

func (c *HTTPClient) InsuranceCompanies() InsuranceCompanyServicer {
	return c.InsuranceCompanySvc
}

func (c *HTTPClient) InsuranceEligibility() InsuranceEligibilityServicer {
	return c.InsuranceEligibilitySvc
}

func (c *HTTPClient) InsurancePlans() InsurancePlanServicer {
	return c.InsurancePlanSvc
}

func (c *HTTPClient) Letters() LetterServicer {
	return c.LetterSvc
}

func (c *HTTPClient) Medications() MedicationServicer {
	return c.MedicationSvc
}

func (c *HTTPClient) MessageThreads() MessageThreadServicer {
	return c.MessageThreadSvc
}

func (c *HTTPClient) NonVisitNotes() NonVisitNoteServicer {
	return c.NonVisitNoteSvc
}

func (c *HTTPClient) Patients() PatientServicer {
	return c.PatientSvc
}

func (c *HTTPClient) Pharmacies() PharmacyServicer {
	return c.PharmacySvc
}

func (c *HTTPClient) Physicians() PhysicianServicer {
	return c.PhysicianSvc
}

func (c *HTTPClient) Practices() PracticeServicer {
	return c.PracticeSvc
}

func (c *HTTPClient) PrescriptionFills() PrescriptionFillServicer {
	return c.PrescriptionFillSvc
}

func (c *HTTPClient) Problems() ProblemServicer {
	return c.ProblemSvc
}

func (c *HTTPClient) RecurringEventGroups() RecurringEventGroupServicer {
	return c.RecurringEventGroupService
}

func (c *HTTPClient) ServiceLocations() ServiceLocationServicer {
	return c.ServiceLocationSvc
}

func (c *HTTPClient) Subscriptions() SubscriptionServicer {
	return c.SubscriptionSvc
}

func (c *HTTPClient) ThreadMembers() ThreadMemberServicer {
	return c.ThreadMemberSvc
}

func (c *HTTPClient) VisitNote() VisitNoteServicer {
	return c.VisitNoteSvc
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

func (c *HTTPClient) request(ctx context.Context, method string, path string, query any, body any, out any) (*http.Response, error) {
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

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
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
