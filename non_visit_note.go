package elation

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type NonVisitNoteServicer interface {
	Create(ctx context.Context, create *NonVisitNoteCreate) (*NonVisitNote, *http.Response, error)
	Find(ctx context.Context, opts *FindNonVisitNotesOptions) (*Response[[]*NonVisitNote], *http.Response, error)
	Get(ctx context.Context, id int64) (*NonVisitNote, *http.Response, error)
}

var _ NonVisitNoteServicer = (*NonVisitNoteService)(nil)

type NonVisitNoteService struct {
	client *HTTPClient
}

type NonVisitNoteCreate struct {
	Bullets      []*NonVisitNoteBullet `json:"bullets"`               //: [{}], 				                                                                                                     // Required
	ChartDate    time.Time             `json:"chart_date"`            //: "2010-06-10T11:05:08Z", 	                                                                                         // Required
	DocumentDate time.Time             `json:"document_date"`         //: "2010-06-10T11:05:08Z", 	                                                                                         // Required
	Patient      int64                 `json:"patient"`               //: 1638401, 				                                                                                                   // Required
	SignedBy     int64                 `json:"signed_by,omitempty"`   //: 131074,
	SignedDate   *time.Time            `json:"signed_date,omitempty"` //: "2010-06-10T11:05:08Z",
	Type         string                `json:"type,omitempty"`        //: "nonvisit", ["email", "nonvisit", "phone"]
	Tags         []*NonVisitNoteTag    `json:"tags,omitempty"`        //: [{}],
}

type NonVisitNote struct {
	ID           int64                   `json:"id"`
	Bullets      []*NonVisitNoteBullet   `json:"bullets"`       //: [{}],
	ChartDate    time.Time               `json:"chart_date"`    //: "2016-10-15T23:32:43Z",
	CreatedDate  time.Time               `json:"created_date"`  //: "2016-10-15T23:32:43Z",
	DeletedDate  *time.Time              `json:"deleted_date"`  //: null,
	DocumentDate time.Time               `json:"document_date"` //: "2016-10-15T23:32:43Z",
	NoteDocument []*NonVisitNoteDocument `json:"note_document"` //: [{}],
	NoteItem     []*NonVisitNoteItem     `json:"note_item"`     //: [{}]
	Notes        []*NonVisitNoteNote     `json:"notes"`         //: [{}],
	Patient      int64                   `json:"patient"`       //: 64058687489,
	SignedBy     int64                   `json:"signed_by"`     //: 131074,
	SignedDate   *time.Time              `json:"signed_date"`   //: "2016-10-15T23:34:18Z",
	Tags         []*NonVisitNoteTag      `json:"tags"`          //: [{}],
	Type         string                  `json:"type"`          //: "nonvisit",                       // ["email", "nonvisit", "phone"]
}

type NonVisitNoteBullet struct {
	ID          int64     `json:"id"`
	Author      int64     `json:"author"`       //: 64,
	Text        string    `json:"text"`         //: "note is here",
	Version     int64     `json:"version"`      //: 2,
	UpdatedDate time.Time `json:"updated_date"` //: "2022-05-15T13:50:09Z"
}

type NonVisitNoteDocument struct {
	ID       int64                        `json:"id"`
	Document NonVisitNoteDocumentDocument `json:"document"`
	Summary  string                       `json:"summary"` //: "Ordervdsv"
}

type NonVisitNoteDocumentDocument struct {
	ID                int64      `json:"id"`                 //: 140758538780701,
	AuthoringPractice int64      `json:"authoring_practice"` //: 65540,
	ChartDate         time.Time  `json:"chart_date"`         //: "2022-05-06T15:26:44Z",
	CreatedDate       time.Time  `json:"created_date"`       //: "2022-05-06T15:26:44Z",
	DeletedDate       *time.Time `json:"deleted_date"`       //: null,
	DocumentDate      time.Time  `json:"document_date"`      //: "2022-05-06T15:26:43Z",
	DocumentType      string     `json:"document_type"`      //: 29,
	LastModified      time.Time  `json:"last_modified"`      //: "2022-05-06T15:26:45Z",
	Patient           int64      `json:"patient"`            //: 64058687489,
	SignDate          *time.Time `json:"sign_date"`          //: "2022-05-06T15:26:44Z",
	SignedBy          int64      `json:"signed_by"`          //: 4
}

type NonVisitNoteNote struct {
	ID          int64     `json:"id"`           //: 1407585389143534,
	Author      int64     `json:"author"`       //: 64,
	Text        string    `json:"text"`         //: "note is here",
	Version     int64     `json:"version"`      //: 2,
	UpdatedDate time.Time `json:"updated_date"` //: "2022-05-15T13:50:09Z"
}

type NonVisitNoteItem struct {
	ID   int64                 `json:"id"`   //: 1407585389143534,
	Item *NonVisitNoteItemItem `json:"item"` //: {}
}

type NonVisitNoteItemItem struct {
	ID             int64      `json:"id"`              //: 140000000000000,
	CreatedDate    *time.Time `json:"created_date"`    //: null,
	DeletedDate    *time.Time `json:"deleted_date"`    //: null,
	Patient        int64      `json:"patient"`         //: 1638401,
	Type           string     `json:"type"`            //: "PatientHistoryItem",
	IsConfidential bool       `json:"is_confidential"` //: false
}

type NonVisitNoteTag struct {
	ID               int64      `json:"id"`                 //: 58719011162,
	Code             string     `json:"code"`               //: "408289007",
	CodeType         int64      `json:"code_type"`          //: 2,
	ConceptName      string     `json:"concept_name"`       //: "wt-mgmt-referral",
	Context          string     `json:"context"`            //: null,
	CreatedDate      time.Time  `json:"created_date"`       //: "2015-07-23T19:30:51Z",
	DeletedDate      *time.Time `json:"deleted_date"`       //: "2019-05-29T00:06:03Z",
	Description      string     `json:"description"`        //: "Referral to weight management program",
	PracticeCreated  int64      `json:"practice_created"`   //: null,
	SnomedResultCode string     `json:"snomed_result_code"` //: null,
	Value            string     `json:"value"`              //: "CQM: Ref to Wt Mgt Program"
}

type FindNonVisitNotesOptions struct {
	*Pagination

	Patient int64 `url:"patient,omitempty"`
}

func (s *NonVisitNoteService) Create(ctx context.Context, create *NonVisitNoteCreate) (*NonVisitNote, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "create non visit note", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	nvn := &NonVisitNote{}

	res, err := s.client.request(ctx, http.MethodPost, "/non_visit_notes", nil, create, &nvn)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return nvn, res, nil
}

func (s *NonVisitNoteService) Find(ctx context.Context, opts *FindNonVisitNotesOptions) (*Response[[]*NonVisitNote], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find non visit notes", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*NonVisitNote]{}

	res, err := s.client.request(ctx, http.MethodGet, "/non_visit_notes", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *NonVisitNoteService) Get(ctx context.Context, id int64) (*NonVisitNote, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get non visit note", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.non_visit_note_id", id)))
	defer span.End()

	out := &NonVisitNote{}

	res, err := s.client.request(ctx, http.MethodGet, "/non_visit_notes/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
