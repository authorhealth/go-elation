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
	Find(ctx context.Context, opts *FindVisitNotesOptions) (*Response[[]*VisitNote], *http.Response, error)
	Get(ctx context.Context, id int64) (*VisitNote, *http.Response, error)
}

var _ VisitNoteServicer = (*VisitNoteService)(nil)

type VisitNoteService struct {
	client *Client
}

type VisitNote struct {
	ID                  int                   `json:"id"`
	Bullets             []*VisitNoteBullet    `json:"bullets"`
	Checklists          *VisitNoteChecklists  `json:"checklists"`
	Edits               []*VisitNoteEdit      `json:"edits"`
	Signatures          []*VisitNoteSignature `json:"signatures"`
	Type                string                `json:"type"`
	Template            string                `json:"template"`
	AmendmentRequest    any                   `json:"amendment_request"`
	Patient             int                   `json:"patient"`
	Physician           int                   `json:"physician"`
	Practice            int                   `json:"practice"`
	DocumentDate        time.Time             `json:"document_date"`
	ChartDate           time.Time             `json:"chart_date"`
	ClinicalSummaryLink string                `json:"clinical_summary_link"`
	VisitSummaryLink    string                `json:"visit_summary_link"`
	SignedDate          time.Time             `json:"signed_date"`
	SignedBy            int                   `json:"signed_by"`
	CreatedDate         time.Time             `json:"created_date"`
	LastModified        time.Time             `json:"last_modified"`
	DeletedDate         *time.Time            `json:"deleted_date"`
	Tags                []any                 `json:"tags"`
	Confidential        bool                  `json:"confidential"`
}

type VisitNoteItemItem struct {
	ID             int        `json:"id"`
	CreatedDate    time.Time  `json:"created_date"`
	DeletedDate    *time.Time `json:"deleted_date"`
	Patient        int        `json:"patient"`
	Type           string     `json:"type"`
	IsConfidential bool       `json:"is_confidential"`
	ItemType       string     `json:"itemType"`
}

type VisitNoteItem struct {
	ID   int                `json:"id"`
	Item *VisitNoteItemItem `json:"item"`
}

type VisitNoteDocumentDocument struct {
	AuthoringPractice int        `json:"authoring_practice"`
	ChartDate         time.Time  `json:"chart_date"`
	CreatedDate       time.Time  `json:"created_date"`
	DeletedDate       *time.Time `json:"deleted_date"`
	DocumentDate      time.Time  `json:"document_date"`
	DocumentType      int        `json:"document_type"`
	ID                int        `json:"id"`
	LastModified      time.Time  `json:"last_modified"`
	Patient           int        `json:"patient"`
	SignDate          time.Time  `json:"sign_date"`
	SignedBy          int        `json:"signed_by"`
}

type VisitNoteDocument struct {
	ID       int                        `json:"id"`
	Document *VisitNoteDocumentDocument `json:"document"`
	Summary  any                        `json:"summary"`
}

type VisitNoteChild struct {
	Category       string             `json:"category"`
	Text           string             `json:"text"`
	Version        int                `json:"version"`
	Sequence       int                `json:"sequence"`
	Author         int                `json:"author"`
	UpdatedDate    string             `json:"updated_date"`
	ReplacedByEdit any                `json:"replaced_by_edit"`
	ReplacedBy     any                `json:"replaced_by"`
	Edit           any                `json:"edit"`
	DeletedDate    *time.Time         `json:"deleted_date"`
	NoteItem       *VisitNoteItem     `json:"note_item"`
	NoteDocument   *VisitNoteDocument `json:"note_document"`
	Handout        int64              `json:"handout"`
}

type VisitNoteBullet struct {
	Category       string            `json:"category"`
	Text           string            `json:"text"`
	Version        int               `json:"version"`
	Sequence       int               `json:"sequence"`
	Author         int               `json:"author"`
	UpdatedDate    string            `json:"updated_date"`
	ReplacedByEdit any               `json:"replaced_by_edit"`
	ReplacedBy     any               `json:"replaced_by"`
	Edit           any               `json:"edit"`
	DeletedDate    *time.Time        `json:"deleted_date"`
	NoteDocument   any               `json:"note_document"`
	NoteItem       any               `json:"note_item"`
	Handout        any               `json:"handout"`
	Children       []*VisitNoteChild `json:"children"`
}

type VisitNoteChecklistPE struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Sequence int    `json:"sequence"`
}

type VisitNoteChecklistROS struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Sequence int    `json:"sequence"`
}

type VisitNoteChecklists struct {
	Pe  []*VisitNoteChecklistPE  `json:"PE"`
	Ros []*VisitNoteChecklistROS `json:"ROS"`
}

type VisitNoteEdit struct {
	VisitNote          int64     `json:"visit_note"`
	CreateUser         int       `json:"create_user"`
	CreatedDate        time.Time `json:"created_date"`
	PreviousNoteType   any       `json:"previous_note_type"`
	PreviousNoteTime   any       `json:"previous_note_time"`
	NewNoteType        any       `json:"new_note_type"`
	NewNoteTime        any       `json:"new_note_time"`
	PreviousPrevWeight any       `json:"previous_prev_weight"`
	NewPrevWeight      any       `json:"new_prev_weight"`
	PreviousPrevBmi    any       `json:"previous_prev_bmi"`
	NewPrevBmi         any       `json:"new_prev_bmi"`
	PreviousPrevTime   any       `json:"previous_prev_time"`
	NewPrevTime        any       `json:"new_prev_time"`
	PreviousBmi        any       `json:"previous_bmi"`
	NewBmi             any       `json:"new_bmi"`
}

type VisitNoteSignature struct {
	User       int       `json:"user"`
	UserName   string    `json:"user_name"`
	SignedDate time.Time `json:"signed_date"`
	Role       string    `json:"role"`
	Comments   any       `json:"comments"`
}

type FindVisitNotesOptions struct {
	*Pagination

	Patient         []int64   `url:"patient,omitempty"`
	Physician       []int64   `url:"physician,omitempty"`
	Practice        []int64   `url:"practice,omitempty"`
	Unsigned        bool      `url:"unsigned,omitempty"`
	FromSignedDate  time.Time `url:"from_signed_date,omitempty"`
	ToSignedDate    time.Time `url:"to_signed_date,omitempty"`
	LastModifiedGT  time.Time `url:"last_modified_gt,omitempty"`
	LastModifiedGTE time.Time `url:"last_modified_gte,omitempty"`
	LastModifiedLT  time.Time `url:"last_modified_lt,omitempty"`
	LastModifiedLTE time.Time `url:"last_modified_lte,omitempty"`
}

func (s *VisitNoteService) Find(ctx context.Context, opts *FindVisitNotesOptions) (*Response[[]*VisitNote], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find vist notes", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*VisitNote]{}

	res, err := s.client.request(ctx, http.MethodGet, "/visit_notes", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *VisitNoteService) Get(ctx context.Context, id int64) (*VisitNote, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get visit note", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.visit note_id", id)))
	defer span.End()

	out := &VisitNote{}

	res, err := s.client.request(ctx, http.MethodGet, "/visit_notes/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
