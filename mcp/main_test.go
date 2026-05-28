package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/authorhealth/go-elation"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestRegisterTools_Comprehensive(t *testing.T) {
	assert := assert.New(t)
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

	assert.Equal(len(expected), len(tools))

	for _, name := range expected {
		assert.Contains(tools, name)
	}
}

func TestRegisterTools_SafeByDefault(t *testing.T) {
	assert := assert.New(t)
	s := server.NewMCPServer("test", "0.0.0")

	registerTools(s, nil, false)

	tools := s.ListTools()

	// Verify that all unsafe tools are hidden
	for _, name := range unsafeToolNames {
		assert.NotContains(tools, name)
	}

	// Verify at least one known safe tool exists
	assert.Contains(tools, "patients_get")
}

func TestRegisterTools_InputSchemasPublished(t *testing.T) {
	assert := assert.New(t)
	s := server.NewMCPServer("test", "0.0.0")

	registerTools(s, nil, true)

	tools := s.ListTools()

	getTool := tools["patients_get"]
	assert.NotNil(getTool)
	getSchema := toolInputSchema(t, getTool.Tool)
	assert.Contains(toolInputSchemaProperties(t, getSchema), "id")
	assert.Contains(toolInputSchemaRequired(getSchema), "id")

	findTool := tools["patients_find"]
	assert.NotNil(findTool)
	findSchema := toolInputSchema(t, findTool.Tool)
	assert.Contains(toolInputSchemaProperties(t, findSchema), "options")
	assert.NotContains(toolInputSchemaRequired(findSchema), "options")

	createTool := tools["patients_create"]
	assert.NotNil(createTool)
	createSchema := toolInputSchema(t, createTool.Tool)
	assert.Contains(toolInputSchemaProperties(t, createSchema), "body")
	assert.Contains(toolInputSchemaRequired(createSchema), "body")

	updateTool := tools["patients_update"]
	assert.NotNil(updateTool)
	updateSchema := toolInputSchema(t, updateTool.Tool)
	assert.Contains(toolInputSchemaProperties(t, updateSchema), "id")
	assert.Contains(toolInputSchemaProperties(t, updateSchema), "body")
	assert.Contains(toolInputSchemaRequired(updateSchema), "id")
	assert.Contains(toolInputSchemaRequired(updateSchema), "body")

	deleteTool := tools["patients_delete"]
	assert.NotNil(deleteTool)
	deleteSchema := toolInputSchema(t, deleteTool.Tool)
	assert.Contains(toolInputSchemaProperties(t, deleteSchema), "id")
	assert.Contains(toolInputSchemaRequired(deleteSchema), "id")
}

func TestRegisterTools_FindOptionsSchemaUsesURLTagsAndOmitempty(t *testing.T) {
	assert := assert.New(t)
	s := server.NewMCPServer("test", "0.0.0")

	registerTools(s, nil, true)

	tools := s.ListTools()
	medsFind := tools["medications_find"]
	assert.NotNil(medsFind)

	schema := toolInputSchema(t, medsFind.Tool)
	rootProps := toolInputSchemaProperties(t, schema)
	optionsProp, ok := rootProps["options"]
	assert.True(ok)

	optionsSchema, ok := optionsProp.(map[string]any)
	assert.True(ok)
	optionsProps := toolInputSchemaProperties(t, optionsSchema)

	// URL-tag names should be exposed to MCP clients, not Go field names.
	assert.Contains(optionsProps, "patient")
	assert.Contains(optionsProps, "practice")
	assert.Contains(optionsProps, "cursor")
	assert.NotContains(optionsProps, "Patient")
	assert.NotContains(optionsProps, "Practice")
	assert.NotContains(optionsProps, "Cursor")

	// `url:",omitempty"` should not be represented as required.
	assert.NotContains(toolInputSchemaRequired(optionsSchema), "patient")
	assert.NotContains(toolInputSchemaRequired(optionsSchema), "practice")
	assert.NotContains(toolInputSchemaRequired(optionsSchema), "cursor")
}

func TestRequireInt64(t *testing.T) {
	t.Run("success from multiple numeric representations", func(t *testing.T) {
		assert := assert.New(t)
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
			t.Run(tc.name, func(t *testing.T) {
				got, err := requireInt64(toolReq(map[string]any{"id": tc.value}), "id")
				assert.NoError(err)
				assert.Equal(tc.want, got)
			})
		}
	})

	t.Run("missing required argument", func(t *testing.T) {
		assert := assert.New(t)
		_, err := requireInt64(toolReq(map[string]any{}), "id")
		assert.Error(err)
		assert.Contains(err.Error(), `missing required argument "id"`)
	})

	t.Run("invalid type", func(t *testing.T) {
		assert := assert.New(t)
		_, err := requireInt64(toolReq(map[string]any{"id": true}), "id")
		assert.Error(err)
		assert.Contains(err.Error(), "must be a number or numeric string")
	})

	t.Run("reject non-integral float64", func(t *testing.T) {
		assert := assert.New(t)
		_, err := requireInt64(toolReq(map[string]any{"id": float64(42.9)}), "id")
		assert.Error(err)
		assert.Contains(err.Error(), "must be an integer")
	})

	t.Run("reject non-integral float32", func(t *testing.T) {
		assert := assert.New(t)
		_, err := requireInt64(toolReq(map[string]any{"id": float32(42.5)}), "id")
		assert.Error(err)
		assert.Contains(err.Error(), "must be an integer")
	})
}

func TestDecodeObjectArg(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("required object present", func(t *testing.T) {
		assert := assert.New(t)
		got, err := decodeObjectArg[payload](toolReq(map[string]any{
			"body": map[string]any{"name": "Ada", "age": 35},
		}), "body", true)
		assert.NoError(err)
		assert.NotNil(got)
		assert.Equal("Ada", got.Name)
		assert.Equal(35, got.Age)
	})

	t.Run("optional object missing returns nil", func(t *testing.T) {
		assert := assert.New(t)
		got, err := decodeObjectArg[payload](toolReq(map[string]any{}), "options", false)
		assert.NoError(err)
		assert.Nil(got)
	})

	t.Run("required object missing", func(t *testing.T) {
		assert := assert.New(t)
		_, err := decodeObjectArg[payload](toolReq(map[string]any{}), "body", true)
		assert.Error(err)
		assert.Contains(err.Error(), `missing required argument "body"`)
	})

	t.Run("invalid object payload", func(t *testing.T) {
		assert := assert.New(t)
		_, err := decodeObjectArg[payload](toolReq(map[string]any{"body": "not-an-object"}), "body", true)
		assert.Error(err)
		assert.Contains(err.Error(), `unmarshaling "body"`)
	})
}

func TestToolErrorResult(t *testing.T) {
	t.Run("elation api error", func(t *testing.T) {
		assert := assert.New(t)
		res := toolErrorResult(&elation.Error{StatusCode: 400, Body: `{"error":"bad request"}`})
		assert.True(res.IsError)
		text := resultText(t, res)
		assert.Contains(text, "status=400")
	})

	t.Run("generic error", func(t *testing.T) {
		assert := assert.New(t)
		res := toolErrorResult(fmt.Errorf("boom"))
		assert.True(res.IsError)
		text := resultText(t, res)
		assert.Equal("boom", text)
	})
}

func TestRegisteredHandler_ValidatesInputBeforeClientCall(t *testing.T) {
	assert := assert.New(t)
	s := server.NewMCPServer("test", "0.0.0")
	registerTools(s, nil, true)

	tool := s.GetTool("patients_get")
	assert.NotNil(tool)

	res, err := tool.Handler(t.Context(), toolReq(map[string]any{}))
	assert.NoError(err)
	assert.True(res.IsError)
	text := resultText(t, res)
	assert.Contains(text, `missing required argument "id"`)
}

func TestParseServerOptions(t *testing.T) {
	t.Run("defaults to safe tools only", func(t *testing.T) {
		assert := assert.New(t)
		opts, err := parseServerOptions(nil)
		assert.NoError(err)
		assert.False(opts.allowUnsafeTools)
	})

	t.Run("enables unsafe tools flag", func(t *testing.T) {
		assert := assert.New(t)
		opts, err := parseServerOptions([]string{"--allow-unsafe-tools"})
		assert.NoError(err)
		assert.True(opts.allowUnsafeTools)
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

func toolInputSchema(t *testing.T, tool mcp.Tool) map[string]any {
	t.Helper()

	b, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("marshaling tool: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatalf("unmarshaling tool json: %v", err)
	}

	rawSchema, ok := payload["inputSchema"]
	if !ok {
		t.Fatal("tool json missing inputSchema")
	}

	schema, ok := rawSchema.(map[string]any)
	if !ok {
		t.Fatalf("inputSchema has unexpected type %T", rawSchema)
	}

	return schema
}

func toolInputSchemaProperties(t *testing.T, schema map[string]any) map[string]any {
	t.Helper()

	rawProps, ok := schema["properties"]
	if !ok {
		t.Fatal("input schema missing properties")
	}

	props, ok := rawProps.(map[string]any)
	if !ok {
		t.Fatalf("properties has unexpected type %T", rawProps)
	}

	return props
}

func toolInputSchemaRequired(schema map[string]any) []string {
	rawRequired, ok := schema["required"]
	if !ok {
		return nil
	}

	requiredValues, ok := rawRequired.([]any)
	if !ok {
		return nil
	}

	required := make([]string, 0, len(requiredValues))
	for _, value := range requiredValues {
		str, ok := value.(string)
		if ok {
			required = append(required, str)
		}
	}

	return required
}
