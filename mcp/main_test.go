package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/authorhealth/go-elation"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestRegisterTools_Comprehensive(t *testing.T) {
	s := server.NewMCPServer("test", "0.0.0")

	registerTools(s, nil, true)

	tools := s.ListTools()
	expected := []string{
		"allergies_find", "allergies_get", "allergy_documentation_find", "allergy_documentation_get",
		"appointments_create", "appointments_find", "appointments_get", "appointments_update", "appointments_delete",
		"bills_create", "bills_find", "bills_get",
		"clinical_documents_find", "clinical_documents_get",
		"contacts_get", "contacts_list",
		"discontinued_medications_create", "discontinued_medications_find", "discontinued_medications_get",
		"history_download_fills_find", "history_download_fills_get",
		"insurance_companies_create", "insurance_companies_find", "insurance_companies_get", "insurance_companies_update", "insurance_companies_delete",
		"insurance_eligibility_create", "insurance_eligibility_get", "insurance_eligibility_get_full_report",
		"insurance_plans_create", "insurance_plans_find", "insurance_plans_get", "insurance_plans_update", "insurance_plans_delete",
		"insurance_policies_create", "insurance_policies_find", "insurance_policies_get", "insurance_policies_update", "insurance_policies_delete",
		"letters_find", "letters_get",
		"medications_create", "medications_find", "medications_get",
		"message_threads_find", "message_threads_get",
		"non_visit_notes_create", "non_visit_notes_find", "non_visit_notes_get",
		"patients_create", "patients_find", "patients_get", "patients_update", "patients_delete",
		"pharmacies_get",
		"physicians_find", "physicians_get",
		"practices_find", "practices_get",
		"prescription_fills_find", "prescription_fills_get",
		"problems_find", "problems_get",
		"recurring_event_groups_create", "recurring_event_groups_find", "recurring_event_groups_get", "recurring_event_groups_update", "recurring_event_groups_delete",
		"service_locations_find",
		"subscriptions_find", "subscriptions_subscribe", "subscriptions_delete",
		"thread_members_find", "thread_members_get",
		"visit_notes_create", "visit_notes_delete", "visit_notes_find", "visit_notes_get",
	}

	if len(tools) != len(expected) {
		t.Fatalf("expected %d tools, got %d", len(expected), len(tools))
	}

	for _, name := range expected {
		if _, ok := tools[name]; !ok {
			t.Fatalf("expected tool %q to be registered", name)
		}
	}
}

func TestRegisterTools_SafeByDefault(t *testing.T) {
	s := server.NewMCPServer("test", "0.0.0")

	registerTools(s, nil, false)

	tools := s.ListTools()
	if got := len(tools); got != 51 {
		t.Fatalf("expected 51 safe tools, got %d", got)
	}

	for _, name := range unsafeToolNames {
		if _, ok := tools[name]; ok {
			t.Fatalf("expected unsafe tool %q to be hidden by default", name)
		}
	}
}

func TestRequireInt64(t *testing.T) {
	t.Run("success from multiple numeric representations", func(t *testing.T) {
		cases := []struct {
			name  string
			value any
			want  int64
		}{
			{name: "float64", value: float64(42), want: 42},
			{name: "float32", value: float32(42), want: 42},
			{name: "int", value: int(42), want: 42},
			{name: "int32", value: int32(42), want: 42},
			{name: "int64", value: int64(42), want: 42},
			{name: "json number", value: json.Number("42"), want: 42},
			{name: "string", value: "42", want: 42},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				got, err := requireInt64(toolReq(map[string]any{"id": tc.value}), "id")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tc.want {
					t.Fatalf("expected %d, got %d", tc.want, got)
				}
			})
		}
	})

	t.Run("missing required argument", func(t *testing.T) {
		_, err := requireInt64(toolReq(map[string]any{}), "id")
		if err == nil || !strings.Contains(err.Error(), `missing required argument "id"`) {
			t.Fatalf("expected missing argument error, got: %v", err)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := requireInt64(toolReq(map[string]any{"id": true}), "id")
		if err == nil || !strings.Contains(err.Error(), "must be a number or numeric string") {
			t.Fatalf("expected invalid type error, got: %v", err)
		}
	})
}

func TestDecodeObjectArg(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("required object present", func(t *testing.T) {
		got, err := decodeObjectArg[payload](toolReq(map[string]any{
			"body": map[string]any{"name": "Ada", "age": 35},
		}), "body", true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil || got.Name != "Ada" || got.Age != 35 {
			t.Fatalf("unexpected decoded payload: %+v", got)
		}
	})

	t.Run("optional object missing returns nil", func(t *testing.T) {
		got, err := decodeObjectArg[payload](toolReq(map[string]any{}), "options", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil options, got %+v", got)
		}
	})

	t.Run("required object missing", func(t *testing.T) {
		_, err := decodeObjectArg[payload](toolReq(map[string]any{}), "body", true)
		if err == nil || !strings.Contains(err.Error(), `missing required argument "body"`) {
			t.Fatalf("expected missing required argument error, got: %v", err)
		}
	})

	t.Run("invalid object payload", func(t *testing.T) {
		_, err := decodeObjectArg[payload](toolReq(map[string]any{"body": "not-an-object"}), "body", true)
		if err == nil || !strings.Contains(err.Error(), `unmarshaling "body"`) {
			t.Fatalf("expected unmarshal error, got: %v", err)
		}
	})
}

func TestToolErrorResult(t *testing.T) {
	t.Run("elation api error", func(t *testing.T) {
		res := toolErrorResult(&elation.Error{StatusCode: 400, Body: `{"error":"bad request"}`})
		if !res.IsError {
			t.Fatal("expected IsError=true")
		}
		text := resultText(t, res)
		if !strings.Contains(text, "status=400") {
			t.Fatalf("expected status in tool error text, got %q", text)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		res := toolErrorResult(fmt.Errorf("boom"))
		if !res.IsError {
			t.Fatal("expected IsError=true")
		}
		if text := resultText(t, res); text != "boom" {
			t.Fatalf("expected boom, got %q", text)
		}
	})
}

func TestRegisteredHandler_ValidatesInputBeforeClientCall(t *testing.T) {
	s := server.NewMCPServer("test", "0.0.0")
	registerTools(s, nil, true)

	tool := s.GetTool("patients_get")
	if tool == nil {
		t.Fatal("expected patients_get tool to exist")
	}

	res, err := tool.Handler(t.Context(), toolReq(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected protocol-level error: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected tool-level error result")
	}
	if text := resultText(t, res); !strings.Contains(text, `missing required argument "id"`) {
		t.Fatalf("unexpected error text: %q", text)
	}
}

func TestParseServerOptions(t *testing.T) {
	t.Run("defaults to safe tools only", func(t *testing.T) {
		opts, err := parseServerOptions(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.allowUnsafeTools {
			t.Fatal("expected allowUnsafeTools=false by default")
		}
	})

	t.Run("enables unsafe tools flag", func(t *testing.T) {
		opts, err := parseServerOptions([]string{"--allow-unsafe-tools"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !opts.allowUnsafeTools {
			t.Fatal("expected allowUnsafeTools=true")
		}
	})
}

func toolReq(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: args,
		},
	}
}

func resultText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()

	if len(res.Content) == 0 {
		t.Fatal("expected at least one content entry")
	}

	textContent, ok := mcp.AsTextContent(res.Content[0])
	if !ok {
		t.Fatalf("expected text content, got %T", res.Content[0])
	}

	return textContent.Text
}
