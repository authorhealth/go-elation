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

type ClinicalDocumentServicer interface {
	Find(ctx context.Context, opts *FindClinicalDocumentsOptions) (*Response[[]*ClinicalDocument], *http.Response, error)
	Get(ctx context.Context, id int64) (*ClinicalDocument, *http.Response, error)
}

var _ ClinicalDocumentServicer = (*ClinicalDocumentService)(nil)

type ClinicalDocumentService struct {
	client *HTTPClient
}

type ClinicalDocument struct {
	ID                    int64                    `json:"id"`
	AuthoringPractice     int64                    `json:"authoring_practice"`     //: 125323,
	Patient               int64                    `json:"patient"`                //: 322134,
	DataFormat            string                   `json:"data_format"`            //: "ccda",
	XMLFile               *ClinicalDocumentXMLFile `json:"xml_file"`               //: {},
	DemographicsImported  bool                     `json:"demographics_imported"`  //: false,
	AllergiesImported     bool                     `json:"allergies_imported"`     //: false,
	EncountersImported    bool                     `json:"encounters_imported"`    //: false,
	ImmunizationsImported bool                     `json:"immunizations_imported"` //: false,
	LabsImported          bool                     `json:"labs_imported"`          //: false,
	MedicationsImported   bool                     `json:"medications_imported"`   //: false,
	ProblemsImported      bool                     `json:"problems_imported"`      //: false,
	ProceduresImported    bool                     `json:"procedures_imported"`    //: false,
	VitalsImported        bool                     `json:"vitals_imported"`        //: false,
	AuthorName            string                   `json:"author_name"`            //: "Author Name",
	CreatedDate           time.Time                `json:"created_date"`           //: "2016-05-02T13:30:07Z",
	DeletedDate           *time.Time               `json:"deleted_date"`           //: null
}

type ClinicalDocumentXMLFile struct {
	ContentType      string `json:"content_type"`      //: "text/xml",
	OriginalFilename string `json:"original_filename"` //: "full_ccda.xml"
}

type FindClinicalDocumentsOptions struct {
	*Pagination

	Patient int64 `url:"patient,omitempty"`
}

func (s *ClinicalDocumentService) Find(ctx context.Context, opts *FindClinicalDocumentsOptions) (*Response[[]*ClinicalDocument], *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "find clinical documents", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	out := &Response[[]*ClinicalDocument]{}

	res, err := s.client.request(ctx, http.MethodGet, "/clinical_documents", opts, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}

func (s *ClinicalDocumentService) Get(ctx context.Context, id int64) (*ClinicalDocument, *http.Response, error) {
	ctx, span := s.client.tracer.Start(ctx, "get clinical document", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attribute.Int64("elation.clinical_document_id", id)))
	defer span.End()

	out := &ClinicalDocument{}

	res, err := s.client.request(ctx, http.MethodGet, "/clinical_documents/"+strconv.FormatInt(id, 10), nil, nil, &out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error making request")
		return nil, res, fmt.Errorf("making request: %w", err)
	}

	return out, res, nil
}
