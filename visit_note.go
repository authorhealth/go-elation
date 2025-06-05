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

type VisitNoteServicer interface {
	Create(ctx context.Context, create *VisitNoteCreate) (*VisitNote, *http.Response, error)
	Delete(ctx context.Context, id int64) (*http.Response, error)
	Find(ctx context.Context, opts *FindVisitNotesOptions) (*Response[[]*VisitNote], *http.Response, error)
	Get(ctx context.Context, id int64) (*VisitNote, *http.Response, error)
}

var _ VisitNoteServicer = (*VisitNoteService)(nil)

type VisitNoteService struct {
	client *HTTPClient
}

type VisitNote struct {
	ID                  int64                 `json:"id"`
	Bullets             []*VisitNoteBullet    `json:"bullets"`               //: [{}],
	Checklists          *VisitNoteChecklists  `json:"checklists"`            //: {},
	Edits               []*VisitNoteEdit      `json:"edits"`                 //: [{}],
	Signatures          []*VisitNoteSignature `json:"signatures"`            //: [{}],
	Type                string                `json:"type"`                  //: "Office Visit Note",              // string(50).  This list is managed in the practice"s setting page
	Template            string                `json:"template"`              //: "SOAP", ["Simple", "SOAP", "Complete H&P (1 col)", "Complete H&P (2 col)", "Complete H&P (2 col A/P)", "Pre-Op"]
	AmendmentRequest    any                   `json:"amendment_request"`     //: null,
	Patient             int64                 `json:"patient"`               //: 1638401,
	Physician           int64                 `json:"physician"`             //: 131074,
	Practice            int64                 `json:"practice"`              //: 34323,
	DocumentDate        time.Time             `json:"document_date"`         //: "2010-06-10T11:05:08Z",
	ChartDate           time.Time             `json:"chart_date"`            //: "2010-06-10T11:05:08Z",
	ClinicalSummaryLink string                `json:"clinical_summary_link"` //: "http://localhost/api/2.0/visit_notes/140758496444440/clinical_summary",
	VisitSummaryLink    string                `json:"visit_summary_link"`    //: "http://localhost/api/2.0/visit_notes/140758496444440/clinical_summary",
	SignedDate          time.Time             `json:"signed_date"`           //: "2010-06-10T11:05:08Z",
	SignedBy            int64                 `json:"signed_by"`             //: 131074,
	CreatedDate         time.Time             `json:"created_date"`          //: "2010-06-10T11:05:08Z",
	LastModified        time.Time             `json:"last_modified"`         //: "2022-05-20T11:06:12.507775Z",
	DeletedDate         *time.Time            `json:"deleted_date"`          //: null,
	Tags                []any                 `json:"tags"`                  //: [],
	Confidential        bool                  `json:"confidential"`          //: false
}

type VisitNoteCreate struct {
	Bullets      []*VisitNoteBullet    `json:"bullets"`                //: [{}], 				                                                                                                     // Required
	ChartDate    time.Time             `json:"chart_date"`             //: "2010-06-10T11:05:08Z", 	                                                                                         // Required
	DocumentDate time.Time             `json:"document_date"`          //: "2010-06-10T11:05:08Z", 	                                                                                         // Required
	Patient      int64                 `json:"patient"`                //: 1638401, 				                                                                                                   // Required
	Template     string                `json:"template"`               //: "SOAP", ["Simple", "SOAP", "Complete H&P (1 col)", "Complete H&P (2 col)", "Complete H&P (2 col A/P)", "Pre-Op"]   // Required
	Physician    int64                 `json:"physician"`              //: 131074, 				          // Required
	Type         string                `json:"type,omitempty"`         //: "Office Visit Note",
	Confidential bool                  `json:"confidential,omitempty"` //: false,
	SignedBy     int64                 `json:"signed_by,omitempty"`    //: 131074,
	SignedDate   *time.Time            `json:"signed_date,omitempty"`  //: "2010-06-10T11:05:08Z",
	Signatures   []*VisitNoteSignature `json:"signatures,omitempty"`   //: [{}],
}

type VisitNoteBullet struct {
	Category       string                 `json:"category"`           //: "Problem", ["Problem", "Past", "Family", "Social", "Instr", "PE", "ROS", "Med", "Data", "Assessment", "Test", "Tx", "Narrative", "Followup", "Reason", "Plan", "Objective", "Hpi", "Allergies", "Habits", "Assessplan", "Consultant", "Attending", "Dateprocedure", "Surgical", "Orders", "Referenced", "Procedure"],
	Text           string                 `json:"text"`               //: "Dizziness" string(500),
	Version        int64                  `json:"version"`            //: 1,
	Sequence       int64                  `json:"sequence"`           //: 0,
	Author         int64                  `json:"author"`             //: 10,
	ReplacedByEdit any                    `json:"replaced_by_edit"`   //: null,
	ReplacedBy     any                    `json:"replaced_by"`        //: null,
	Edit           any                    `json:"edit"`               //: null,
	DeletedDate    *time.Time             `json:"deleted_date"`       //: null,
	NoteDocument   *VisitNoteNoteDocument `json:"note_document"`      //: null,
	NoteItem       *VisitNoteNoteItem     `json:"note_item"`          //: null,
	Handout        *int64                 `json:"handout"`            //: null,
	Children       []*VisitNoteChild      `json:"children,omitempty"` //: [{}]          // Not visible in the Find Visit Notes endpoint
}

type VisitNoteChild struct {
	Category       string                 `json:"category"`         //: "Problem", ["Problem", "Past", "Family", "Social", "Instr", "PE", "ROS", "Med", "Data", "Assessment", "Test", "Tx", "Narrative", "Followup", "Reason", "Plan", "Objective", "Hpi", "Allergies", "Habits", "Assessplan", "Consultant", "Attending", "Dateprocedure", "Surgical", "Orders", "Referenced", "Procedure"],
	Text           string                 `json:"text"`             //: "Dizziness", string(500),
	Version        int64                  `json:"version"`          //: 1,
	Sequence       int64                  `json:"sequence"`         //: 0,
	Author         int64                  `json:"author"`           //: 10,
	UpdatedDate    *string                `json:"updated_date"`     //: "2022-05-15T13:50:09" (Missing TZ offset, but can assume PT)
	ReplacedByEdit any                    `json:"replaced_by_edit"` //: null,
	ReplacedBy     any                    `json:"replaced_by"`      //: null,
	Edit           any                    `json:"edit"`             //: null,
	DeletedDate    *time.Time             `json:"deleted_date"`     //: null,
	NoteItem       *VisitNoteNoteItem     `json:"note_item"`        //: null,
	NoteDocument   *VisitNoteNoteDocument `json:"note_document"`    //: null,
	Handout        *int64                 `json:"handout"`          //: 140758529736869
}

type VisitNoteNoteItem struct {
	ID   int64          `json:"id"`   //: 338362508,
	Item *VisitNoteItem `json:"item"` //: {}
}

type VisitNoteItem struct {
	ID             int64      `json:"id"`              //: 240582699
	CreatedDate    time.Time  `json:"created_date"`    //: "2013-04-09T18:41:23Z",
	DeletedDate    *time.Time `json:"deleted_date"`    //: null,
	Patient        int64      `json:"patient"`         //: 234160129,
	Type           string     `json:"type"`            //: "PatientProblem",
	IsConfidential bool       `json:"is_confidential"` //: false,
	ItemType       string     `json:"itemType"`        //: "PatientProblem"
}

type VisitNoteNoteDocument struct {
	ID       int64              `json:"id"`       //: 99483787,
	Document *VisitNoteDocument `json:"document"` //: {}
	Summary  *string            `json:"summary"`  //: null
}

type VisitNoteDocument struct {
	AuthoringPractice int64      `json:"authoring_practice"` // 65540,
	ChartDate         time.Time  `json:"chart_date"`         // "2010-06-10T17:57:35Z",
	CreatedDate       time.Time  `json:"created_date"`       // "2010-06-10T17:57:35Z",
	DeletedDate       *time.Time `json:"deleted_date"`       // null,
	DocumentDate      time.Time  `json:"document_date"`      // "2010-06-10T17:57:35Z",
	DocumentType      int64      `json:"document_type"`      // 33,
	ID                int64      `json:"id"`                 // 2687009,
	LastModified      *time.Time `json:"last_modified"`      // null,
	Patient           int64      `json:"patient"`            // 1638401,
	SignDate          time.Time  `json:"sign_date"`          // "2010-06-10T17:57:35Z",
	SignedBy          int64      `json:"signed_by"`          // 4
}

type VisitNoteChecklists struct {
	PE  []*VisitNoteChecklistItem `json:"PE"`  //: [{}]   // Physical Exam
	ROS []*VisitNoteChecklistItem `json:"ROS"` //: [{}]   // Review of Systems
}

type VisitNoteChecklistItem struct {
	Name     string `json:"name"`     //: "General",
	Value    string `json:"value"`    //: "well nourished",
	Sequence int64  `json:"sequence"` // 0
}

type VisitNoteEdit struct {
	VisitNote          int64     `json:"visit_note"`           //: 140755736526872,
	CreateUser         int64     `json:"create_user"`          //: 1273,
	CreatedDate        time.Time `json:"created_date"`         //: "2021-07-21T08:14:40Z",
	PreviousNoteType   any       `json:"previous_note_type"`   //: null,
	PreviousNoteTime   any       `json:"previous_note_time"`   //: null,
	NewNoteType        any       `json:"new_note_type"`        //: null,
	NewNoteTime        any       `json:"new_note_time"`        //: null,
	PreviousPrevWeight any       `json:"previous_prev_weight"` //: null,
	NewPrevWeight      any       `json:"new_prev_weight"`      //: null,
	PreviousPrevBMI    any       `json:"previous_prev_bmi"`    //: null,
	NewPrevBMI         any       `json:"new_prev_bmi"`         //: null,
	PreviousPrevTime   any       `json:"previous_prev_time"`   //: null,
	NewPrevTime        any       `json:"new_prev_time"`        //: null,
	PreviousBMI        any       `json:"previous_bmi"`         //: null,
	NewBMI             any       `json:"new_bmi"`              //: null,
}

type VisitNoteSignature struct {
	User       int64     `json:"user"`        //: 12,
	UserName   string    `json:"user_name"`   //: "Douglas Ross, MD",
	SignedDate time.Time `json:"signed_date"` //: "2022-10-31T12:30:00Z",
	Role       string    `json:"role"`        //: "cosigner",
	Comments   *string   `json:"comments"`    //: null
}

func (v *VisitNoteService) Create(ctx context.Context, create *VisitNoteCreate) (*VisitNote, *http.Response, error) {
	ctx, span := v.client.tracer.Start(ctx, "create visit note", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	vn := &VisitNote{}

	res, err := v.client.request(ctx, http.MethodPost, "/visit_notes", nil, create, &vn)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return vn, res, nil
}

type FindVisitNotesOptions struct {
	*Pagination

	Patient   int64 `url:"patient,omitempty"`
	Physician int64 `url:"physician,omitempty"`
	Practice  int64 `url:"practice,omitempty"`

	LastModifiedGT  time.Time `url:"last_modified__gt,omitempty"`
	LastModifiedGTE time.Time `url:"last_modified__gte,omitempty"`
	LastModifiedLT  time.Time `url:"last_modified__lt,omitempty"`
	LastModifiedLTE time.Time `url:"last_modified__lte,omitempty"`

	FromSignedDate time.Time `url:"from_signed_date,omitempty"`
	ToSignedDate   time.Time `url:"to_signed_date,omitempty"`
	Unsigned       bool      `url:"unsigned,omitempty"`
}

func (v *VisitNoteService) Find(ctx context.Context, opts *FindVisitNotesOptions) (*Response[[]*VisitNote], *http.Response, error) {
	ctx, span := v.client.tracer.Start(ctx, "find visit notes", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*VisitNote]{}

	res, err := v.client.request(ctx, http.MethodGet, "/visit_notes", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (v *VisitNoteService) Delete(ctx context.Context, id int64) (*http.Response, error) {
	ctx, span := v.client.tracer.Start(ctx, "delete visit note", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.visit_note_id", id)))
	defer span.End()

	res, err := v.client.request(ctx, http.MethodDelete, "/visit_notes/"+strconv.FormatInt(id, 10), nil, nil, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}

func (v *VisitNoteService) Get(ctx context.Context, id int64) (*VisitNote, *http.Response, error) {
	ctx, span := v.client.tracer.Start(ctx, "get visit note", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.visit_note_id", id)))
	defer span.End()

	out := &VisitNote{}

	res, err := v.client.request(ctx, http.MethodGet, "/visit_notes/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
