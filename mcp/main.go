package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/authorhealth/go-elation"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultHTTPTimeout = 15 * time.Second
	serverName         = "go-elation-mcp"
	serverVersion      = "0.2.0"
)

var unsafeToolNames = []string{
	"appointments_create", "appointments_update", "appointments_delete",
	"bills_create",
	"discontinued_medications_create",
	"insurance_companies_create", "insurance_companies_update", "insurance_companies_delete",
	"insurance_eligibility_create",
	"insurance_plans_create", "insurance_plans_update", "insurance_plans_delete",
	"insurance_policies_create", "insurance_policies_update", "insurance_policies_delete",
	"medications_create",
	"non_visit_notes_create",
	"patients_create", "patients_update", "patients_delete",
	"recurring_event_groups_create", "recurring_event_groups_update", "recurring_event_groups_delete",
	"subscriptions_subscribe", "subscriptions_delete",
	"visit_notes_create", "visit_notes_delete",
}

type serverOptions struct {
	allowUnsafeTools bool
}

func main() {
	opts, err := parseServerOptions(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	client := elation.NewHTTPClient(
		&http.Client{Timeout: defaultHTTPTimeout},
		os.Getenv("ELATION_TOKEN_URL"),
		os.Getenv("ELATION_CLIENT_ID"),
		os.Getenv("ELATION_CLIENT_SECRET"),
		os.Getenv("ELATION_BASE_URL"),
	)

	s := server.NewMCPServer(serverName, serverVersion)
	registerTools(s, client, opts.allowUnsafeTools)

	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}

func parseServerOptions(args []string) (serverOptions, error) {
	var opts serverOptions

	fs := flag.NewFlagSet(serverName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.BoolVar(&opts.allowUnsafeTools, "allow-unsafe-tools", false, "Expose mutative tools (create/update/delete)")

	if err := fs.Parse(args); err != nil {
		return serverOptions{}, err
	}

	return opts, nil
}

func registerTools(s *server.MCPServer, client elation.Client, allowUnsafeTools bool) {
	registerFindTool[elation.FindAllergiesOptions](s, "allergies_find", "Find allergies", func(ctx context.Context, opts *elation.FindAllergiesOptions) (any, error) {
		out, _, err := client.Allergies().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "allergies_get", "Get an allergy by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Allergies().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindAllergiesDocumentationOptions](s, "allergy_documentation_find", "Find allergy documentation", func(ctx context.Context, opts *elation.FindAllergiesDocumentationOptions) (any, error) {
		out, _, err := client.AllergyDocumentation().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "allergy_documentation_get", "Get allergy documentation by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.AllergyDocumentation().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.AppointmentCreate](s, "appointments_create", "Create an appointment", func(ctx context.Context, body *elation.AppointmentCreate) (any, error) {
		out, _, err := client.Appointments().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindAppointmentsOptions](s, "appointments_find", "Find appointments", func(ctx context.Context, opts *elation.FindAppointmentsOptions) (any, error) {
		out, _, err := client.Appointments().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "appointments_get", "Get an appointment by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Appointments().Get(ctx, id)
		return out, err
	})
	registerUpdateTool[elation.AppointmentUpdate](s, "appointments_update", "Update an appointment", func(ctx context.Context, id int64, body *elation.AppointmentUpdate) (any, error) {
		out, _, err := client.Appointments().Update(ctx, id, body)
		return out, err
	})
	registerDeleteTool(s, "appointments_delete", "Delete an appointment", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.Appointments().Delete(ctx, id)
	})

	registerCreateTool[elation.BillCreate](s, "bills_create", "Create a bill", func(ctx context.Context, body *elation.BillCreate) (any, error) {
		out, _, err := client.Bill().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindBillOptions](s, "bills_find", "Find bills", func(ctx context.Context, opts *elation.FindBillOptions) (any, error) {
		out, _, err := client.Bill().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "bills_get", "Get a bill by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Bill().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindClinicalDocumentsOptions](s, "clinical_documents_find", "Find clinical documents", func(ctx context.Context, opts *elation.FindClinicalDocumentsOptions) (any, error) {
		out, _, err := client.ClinicalDocuments().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "clinical_documents_get", "Get a clinical document by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.ClinicalDocuments().Get(ctx, id)
		return out, err
	})

	registerGetTool(s, "contacts_get", "Get a contact by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Contacts().Get(ctx, id)
		return out, err
	})
	registerFindTool[elation.ListContactsOptions](s, "contacts_list", "List contacts", func(ctx context.Context, opts *elation.ListContactsOptions) (any, error) {
		out, _, err := client.Contacts().List(ctx, opts)
		return out, err
	})

	registerCreateTool[elation.DiscontinuedMedicationCreate](s, "discontinued_medications_create", "Create discontinued medication", func(ctx context.Context, body *elation.DiscontinuedMedicationCreate) (any, error) {
		out, _, err := client.DiscontinuedMedications().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindDiscontinuedMedicationsOptions](s, "discontinued_medications_find", "Find discontinued medications", func(ctx context.Context, opts *elation.FindDiscontinuedMedicationsOptions) (any, error) {
		out, _, err := client.DiscontinuedMedications().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "discontinued_medications_get", "Get discontinued medication by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.DiscontinuedMedications().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindHistoryDownloadFillsOptions](s, "history_download_fills_find", "Find history download fills", func(ctx context.Context, opts *elation.FindHistoryDownloadFillsOptions) (any, error) {
		out, _, err := client.HistoryDownloadFills().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "history_download_fills_get", "Get history download fill by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.HistoryDownloadFills().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.InsuranceCompanyCreate](s, "insurance_companies_create", "Create an insurance company", func(ctx context.Context, body *elation.InsuranceCompanyCreate) (any, error) {
		out, _, err := client.InsuranceCompanies().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindInsuranceCompaniesOptions](s, "insurance_companies_find", "Find insurance companies", func(ctx context.Context, opts *elation.FindInsuranceCompaniesOptions) (any, error) {
		out, _, err := client.InsuranceCompanies().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "insurance_companies_get", "Get insurance company by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.InsuranceCompanies().Get(ctx, id)
		return out, err
	})
	registerUpdateTool[elation.InsuranceCompanyUpdate](s, "insurance_companies_update", "Update an insurance company", func(ctx context.Context, id int64, body *elation.InsuranceCompanyUpdate) (any, error) {
		out, _, err := client.InsuranceCompanies().Update(ctx, id, body)
		return out, err
	})
	registerDeleteTool(s, "insurance_companies_delete", "Delete an insurance company", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.InsuranceCompanies().Delete(ctx, id)
	})

	registerScopedCreateTool[elation.InsuranceEligibilityCreate](s, "insurance_eligibility_create", "Create insurance eligibility for a patient insurance", "patient_insurance_id", func(ctx context.Context, patientInsuranceID int64, body *elation.InsuranceEligibilityCreate) (any, error) {
		out, _, err := client.InsuranceEligibility().Create(ctx, patientInsuranceID, body)
		return out, err
	})
	registerScopedGetTool(s, "insurance_eligibility_get", "Get insurance eligibility for a patient insurance", "patient_insurance_id", func(ctx context.Context, patientInsuranceID int64) (any, error) {
		out, _, err := client.InsuranceEligibility().Get(ctx, patientInsuranceID)
		return out, err
	})
	registerScopedGetTool(s, "insurance_eligibility_get_full_report", "Get insurance eligibility full report for a patient insurance", "patient_insurance_id", func(ctx context.Context, patientInsuranceID int64) (any, error) {
		out, _, err := client.InsuranceEligibility().GetFullReport(ctx, patientInsuranceID)
		return out, err
	})

	registerCreateTool[elation.InsurancePlanCreate](s, "insurance_plans_create", "Create an insurance plan", func(ctx context.Context, body *elation.InsurancePlanCreate) (any, error) {
		out, _, err := client.InsurancePlans().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindInsurancePlansOptions](s, "insurance_plans_find", "Find insurance plans", func(ctx context.Context, opts *elation.FindInsurancePlansOptions) (any, error) {
		out, _, err := client.InsurancePlans().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "insurance_plans_get", "Get insurance plan by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.InsurancePlans().Get(ctx, id)
		return out, err
	})
	registerUpdateTool[elation.InsurancePlanUpdate](s, "insurance_plans_update", "Update an insurance plan", func(ctx context.Context, id int64, body *elation.InsurancePlanUpdate) (any, error) {
		out, _, err := client.InsurancePlans().Update(ctx, id, body)
		return out, err
	})
	registerDeleteTool(s, "insurance_plans_delete", "Delete an insurance plan", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.InsurancePlans().Delete(ctx, id)
	})

	registerParentCreateTool[elation.InsurancePolicyCreate](s, "insurance_policies_create", "Create a patient insurance policy", "patient_id", func(ctx context.Context, patientID int64, body *elation.InsurancePolicyCreate) (any, error) {
		out, _, err := client.InsurancePolicies().Create(ctx, patientID, body)
		return out, err
	})
	registerParentFindTool[elation.FindInsurancePoliciesOptions](s, "insurance_policies_find", "Find insurance policies for a patient", "patient_id", func(ctx context.Context, patientID int64, opts *elation.FindInsurancePoliciesOptions) (any, error) {
		out, _, err := client.InsurancePolicies().Find(ctx, patientID, opts)
		return out, err
	})
	registerParentGetTool(s, "insurance_policies_get", "Get an insurance policy by patient and policy ID", "patient_id", func(ctx context.Context, patientID int64, id int64) (any, error) {
		out, _, err := client.InsurancePolicies().Get(ctx, patientID, id)
		return out, err
	})
	registerParentUpdateTool[elation.InsurancePolicyUpdate](s, "insurance_policies_update", "Update an insurance policy by patient and policy ID", "patient_id", func(ctx context.Context, patientID int64, id int64, body *elation.InsurancePolicyUpdate) (any, error) {
		out, _, err := client.InsurancePolicies().Update(ctx, patientID, id, body)
		return out, err
	})
	registerParentDeleteTool(s, "insurance_policies_delete", "Delete an insurance policy by patient and policy ID", "patient_id", func(ctx context.Context, patientID int64, id int64) (*http.Response, error) {
		return client.InsurancePolicies().Delete(ctx, patientID, id)
	})

	registerFindTool[elation.FindLettersOptions](s, "letters_find", "Find letters", func(ctx context.Context, opts *elation.FindLettersOptions) (any, error) {
		out, _, err := client.Letters().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "letters_get", "Get a letter by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Letters().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.PatientMedicationCreate](s, "medications_create", "Create a medication", func(ctx context.Context, body *elation.PatientMedicationCreate) (any, error) {
		out, _, err := client.Medications().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindPatientMedicationsOptions](s, "medications_find", "Find medications", func(ctx context.Context, opts *elation.FindPatientMedicationsOptions) (any, error) {
		out, _, err := client.Medications().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "medications_get", "Get a medication by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Medications().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindMessageThreadsOptions](s, "message_threads_find", "Find message threads", func(ctx context.Context, opts *elation.FindMessageThreadsOptions) (any, error) {
		out, _, err := client.MessageThreads().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "message_threads_get", "Get a message thread by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.MessageThreads().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.NonVisitNoteCreate](s, "non_visit_notes_create", "Create a non-visit note", func(ctx context.Context, body *elation.NonVisitNoteCreate) (any, error) {
		out, _, err := client.NonVisitNotes().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindNonVisitNotesOptions](s, "non_visit_notes_find", "Find non-visit notes", func(ctx context.Context, opts *elation.FindNonVisitNotesOptions) (any, error) {
		out, _, err := client.NonVisitNotes().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "non_visit_notes_get", "Get a non-visit note by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.NonVisitNotes().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.PatientCreate](s, "patients_create", "Create a patient", func(ctx context.Context, body *elation.PatientCreate) (any, error) {
		out, _, err := client.Patients().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindPatientsOptions](s, "patients_find", "Find patients", func(ctx context.Context, opts *elation.FindPatientsOptions) (any, error) {
		out, _, err := client.Patients().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "patients_get", "Get a patient by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Patients().Get(ctx, id)
		return out, err
	})
	registerUpdateTool[elation.PatientUpdate](s, "patients_update", "Update a patient", func(ctx context.Context, id int64, body *elation.PatientUpdate) (any, error) {
		out, _, err := client.Patients().Update(ctx, id, body)
		return out, err
	})
	registerDeleteTool(s, "patients_delete", "Delete a patient", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.Patients().Delete(ctx, id)
	})

	registerGetStringTool(s, "pharmacies_get", "Get a pharmacy by NCPDP ID", "ncpdpid", func(ctx context.Context, ncpdpid string) (any, error) {
		out, _, err := client.Pharmacies().Get(ctx, ncpdpid)
		return out, err
	})

	registerFindTool[elation.FindPhysiciansOptions](s, "physicians_find", "Find physicians", func(ctx context.Context, opts *elation.FindPhysiciansOptions) (any, error) {
		out, _, err := client.Physicians().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "physicians_get", "Get a physician by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Physicians().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindPracticesOptions](s, "practices_find", "Find practices", func(ctx context.Context, opts *elation.FindPracticesOptions) (any, error) {
		out, _, err := client.Practices().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "practices_get", "Get a practice by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Practices().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindPrescriptionFillsOptions](s, "prescription_fills_find", "Find prescription fills", func(ctx context.Context, opts *elation.FindPrescriptionFillsOptions) (any, error) {
		out, _, err := client.PrescriptionFills().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "prescription_fills_get", "Get a prescription fill by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.PrescriptionFills().Get(ctx, id)
		return out, err
	})

	registerFindTool[elation.FindPatientProblemsOptions](s, "problems_find", "Find patient problems", func(ctx context.Context, opts *elation.FindPatientProblemsOptions) (any, error) {
		out, _, err := client.Problems().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "problems_get", "Get a patient problem by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.Problems().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.RecurringEventGroupCreate](s, "recurring_event_groups_create", "Create a recurring event group", func(ctx context.Context, body *elation.RecurringEventGroupCreate) (any, error) {
		out, _, err := client.RecurringEventGroups().Create(ctx, body)
		return out, err
	})
	registerFindTool[elation.FindRecurringEventGroupsOptions](s, "recurring_event_groups_find", "Find recurring event groups", func(ctx context.Context, opts *elation.FindRecurringEventGroupsOptions) (any, error) {
		out, _, err := client.RecurringEventGroups().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "recurring_event_groups_get", "Get a recurring event group by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.RecurringEventGroups().Get(ctx, id)
		return out, err
	})
	registerUpdateTool[elation.RecurringEventGroupUpdate](s, "recurring_event_groups_update", "Update a recurring event group", func(ctx context.Context, id int64, body *elation.RecurringEventGroupUpdate) (any, error) {
		out, _, err := client.RecurringEventGroups().Update(ctx, id, body)
		return out, err
	})
	registerDeleteTool(s, "recurring_event_groups_delete", "Delete a recurring event group", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.RecurringEventGroups().Delete(ctx, id)
	})

	registerFindTool[elation.FindServiceLocationOptions](s, "service_locations_find", "Find service locations", func(ctx context.Context, opts *elation.FindServiceLocationOptions) (any, error) {
		out, _, err := client.ServiceLocations().Find(ctx, opts)
		return out, err
	})

	registerNoArgsTool(s, "subscriptions_find", "List subscriptions", func(ctx context.Context) (any, error) {
		out, _, err := client.Subscriptions().Find(ctx)
		return out, err
	})
	registerCreateTool[elation.Subscribe](s, "subscriptions_subscribe", "Create a subscription", func(ctx context.Context, body *elation.Subscribe) (any, error) {
		out, _, err := client.Subscriptions().Subscribe(ctx, body)
		return out, err
	})
	registerDeleteTool(s, "subscriptions_delete", "Delete a subscription by ID", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.Subscriptions().Delete(ctx, id)
	})

	registerFindTool[elation.FindThreadMembersOptions](s, "thread_members_find", "Find thread members", func(ctx context.Context, opts *elation.FindThreadMembersOptions) (any, error) {
		out, _, err := client.ThreadMembers().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "thread_members_get", "Get a thread member by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.ThreadMembers().Get(ctx, id)
		return out, err
	})

	registerCreateTool[elation.VisitNoteCreate](s, "visit_notes_create", "Create a visit note", func(ctx context.Context, body *elation.VisitNoteCreate) (any, error) {
		out, _, err := client.VisitNote().Create(ctx, body)
		return out, err
	})
	registerDeleteTool(s, "visit_notes_delete", "Delete a visit note", func(ctx context.Context, id int64) (*http.Response, error) {
		return client.VisitNote().Delete(ctx, id)
	})
	registerFindTool[elation.FindVisitNotesOptions](s, "visit_notes_find", "Find visit notes", func(ctx context.Context, opts *elation.FindVisitNotesOptions) (any, error) {
		out, _, err := client.VisitNote().Find(ctx, opts)
		return out, err
	})
	registerGetTool(s, "visit_notes_get", "Get a visit note by ID", func(ctx context.Context, id int64) (any, error) {
		out, _, err := client.VisitNote().Get(ctx, id)
		return out, err
	})

	if !allowUnsafeTools {
		s.DeleteTools(unsafeToolNames...)
	}
}

func registerNoArgsTool(s *server.MCPServer, name, description string, fn func(context.Context) (any, error)) {
	s.AddTool(
		mcp.NewTool(name, mcp.WithDescription(description)),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			out, err := fn(ctx)
			if err != nil {
				return toolErrorResult(err), nil
			}

			return jsonResult(out)
		},
	)
}

func registerGetTool(s *server.MCPServer, name, description string, fn func(context.Context, int64) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, id)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerGetStringTool(s *server.MCPServer, name, description, key string, fn func(context.Context, string) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithString(key, mcp.Description("Resource identifier"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			value, err := req.RequireString(key)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, value)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerScopedGetTool(s *server.MCPServer, name, description, scopeKey string, fn func(context.Context, int64) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(scopeKey, mcp.Description("Scope ID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scopeID, err := requireInt64(req, scopeKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, scopeID)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerFindTool[OptionsT any](s *server.MCPServer, name, description string, fn func(context.Context, *OptionsT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithObject("options", mcp.Description("Query options object")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			opts, err := decodeObjectArg[OptionsT](req, "options", false)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, opts)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerCreateTool[BodyT any](s *server.MCPServer, name, description string, fn func(context.Context, *BodyT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithObject("body", mcp.Description("Request body object"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body, err := decodeObjectArg[BodyT](req, "body", true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, body)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerUpdateTool[BodyT any](s *server.MCPServer, name, description string, fn func(context.Context, int64, *BodyT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
			mcp.WithObject("body", mcp.Description("Update body object"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body, err := decodeObjectArg[BodyT](req, "body", true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, id, body)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerDeleteTool(s *server.MCPServer, name, description string, fn func(context.Context, int64) (*http.Response, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			res, callErr := fn(ctx, id)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(map[string]any{
				"ok":          true,
				"status_code": statusCode(res),
			})
		},
	)
}

func registerScopedCreateTool[BodyT any](s *server.MCPServer, name, description, scopeKey string, fn func(context.Context, int64, *BodyT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(scopeKey, mcp.Description("Scope ID"), mcp.Required()),
			mcp.WithObject("body", mcp.Description("Request body object"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			scopeID, err := requireInt64(req, scopeKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body, err := decodeObjectArg[BodyT](req, "body", true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, scopeID, body)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerParentFindTool[OptionsT any](s *server.MCPServer, name, description, parentKey string, fn func(context.Context, int64, *OptionsT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(parentKey, mcp.Description("Parent ID"), mcp.Required()),
			mcp.WithObject("options", mcp.Description("Query options object")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			parentID, err := requireInt64(req, parentKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts, err := decodeObjectArg[OptionsT](req, "options", false)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, parentID, opts)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerParentCreateTool[BodyT any](s *server.MCPServer, name, description, parentKey string, fn func(context.Context, int64, *BodyT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(parentKey, mcp.Description("Parent ID"), mcp.Required()),
			mcp.WithObject("body", mcp.Description("Request body object"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			parentID, err := requireInt64(req, parentKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body, err := decodeObjectArg[BodyT](req, "body", true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, parentID, body)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerParentGetTool(s *server.MCPServer, name, description, parentKey string, fn func(context.Context, int64, int64) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(parentKey, mcp.Description("Parent ID"), mcp.Required()),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			parentID, err := requireInt64(req, parentKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, parentID, id)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerParentUpdateTool[BodyT any](s *server.MCPServer, name, description, parentKey string, fn func(context.Context, int64, int64, *BodyT) (any, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(parentKey, mcp.Description("Parent ID"), mcp.Required()),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
			mcp.WithObject("body", mcp.Description("Update body object"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			parentID, err := requireInt64(req, parentKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			body, err := decodeObjectArg[BodyT](req, "body", true)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			out, callErr := fn(ctx, parentID, id, body)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(out)
		},
	)
}

func registerParentDeleteTool(s *server.MCPServer, name, description, parentKey string, fn func(context.Context, int64, int64) (*http.Response, error)) {
	s.AddTool(
		mcp.NewTool(name,
			mcp.WithDescription(description),
			mcp.WithNumber(parentKey, mcp.Description("Parent ID"), mcp.Required()),
			mcp.WithNumber("id", mcp.Description("Resource ID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			parentID, err := requireInt64(req, parentKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			id, err := requireInt64(req, "id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			res, callErr := fn(ctx, parentID, id)
			if callErr != nil {
				return toolErrorResult(callErr), nil
			}

			return jsonResult(map[string]any{
				"ok":          true,
				"status_code": statusCode(res),
			})
		},
	)
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshaling result: %w", err)
	}

	return mcp.NewToolResultText(string(b)), nil
}

func toolErrorResult(err error) *mcp.CallToolResult {
	apiErr := &elation.Error{}
	if errors.As(err, &apiErr) {
		return mcp.NewToolResultError(fmt.Sprintf("elation API error: status=%d body=%s", apiErr.StatusCode, apiErr.Body))
	}

	return mcp.NewToolResultError(err.Error())
}

func decodeObjectArg[T any](req mcp.CallToolRequest, key string, required bool) (*T, error) {
	args := req.GetArguments()
	raw, found := args[key]
	if !found || raw == nil {
		if required {
			return nil, fmt.Errorf("missing required argument %q", key)
		}
		return nil, nil
	}

	b, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshaling %q: %w", key, err)
	}

	out := new(T)
	if err := json.Unmarshal(b, out); err != nil {
		return nil, fmt.Errorf("unmarshaling %q: %w", key, err)
	}

	return out, nil
}

func requireInt64(req mcp.CallToolRequest, key string) (int64, error) {
	args := req.GetArguments()
	raw, found := args[key]
	if !found || raw == nil {
		return 0, fmt.Errorf("missing required argument %q", key)
	}

	switch v := raw.(type) {
	case float64:
		// Reject non-integral floats (e.g., 42.9)
		if v != float64(int64(v)) {
			return 0, fmt.Errorf("argument %q must be an integer (got non-integral float)", key)
		}
		// Reject NaN and Infinity
		if v != v || v > float64(^uint64(0)>>1) || v < -float64(^uint64(0)>>1) {
			return 0, fmt.Errorf("argument %q: invalid float value", key)
		}
		return int64(v), nil
	case float32:
		// Reject non-integral floats (e.g., 42.9)
		if v != float32(int64(v)) {
			return 0, fmt.Errorf("argument %q must be an integer (got non-integral float)", key)
		}
		// Reject NaN and Infinity
		if v != v || v > float32(^uint32(0)>>1) || v < -float32(^uint32(0)>>1) {
			return 0, fmt.Errorf("argument %q: invalid float value", key)
		}
		return int64(v), nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("parsing %q: %w", key, err)
		}
		return n, nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parsing %q: %w", key, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("argument %q must be a number or numeric string", key)
	}
}

func statusCode(res *http.Response) int {
	if res == nil {
		return 0
	}
	return res.StatusCode
}
