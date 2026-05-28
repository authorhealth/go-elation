# go-elation MCP server

This module exposes the `github.com/authorhealth/go-elation` client APIs as MCP tools over **stdio**.

## Setup

1. Configure environment variables:
   - `ELATION_BASE_URL`
   - `ELATION_CLIENT_ID`
   - `ELATION_CLIENT_SECRET`
   - `ELATION_TOKEN_URL`
2. You can copy defaults from `.env-default`.

## Run

```bash
cd mcp
go run .
```

The process speaks MCP on stdin/stdout and is intended to be launched by an MCP host (for example, VS Code/Copilot Chat MCP configuration).

By default, only **safe** (non-mutative) tools are exposed.

To also expose mutative tools (`create`, `update`, `delete`, and subscription writes), start with:

```bash
cd mcp
go run . --allow-unsafe-tools
```

## Tool naming and argument conventions

Tools use `<resource>_<action>` naming, for example:

- `patients_get`
- `patients_find`
- `patients_create`
- `patients_update`
- `patients_delete`

Common argument shapes:

- `id`: primary resource identifier
- `body`: request payload object for `create`/`update` actions
- `options`: query/filter object for `find`/`list` actions
- scoped IDs for nested resources (for example `patient_id`, `patient_insurance_id`)

## Response format

- Successful calls return one text content item containing JSON for the underlying API response object.
- Delete calls return:
  - `ok`: boolean
  - `status_code`: upstream HTTP status code (when available)

## Error behavior

- Tool errors are returned as MCP **tool results** (`isError=true`) instead of protocol-level failures.
- `elation.Error` values are surfaced with API status code and response body to aid troubleshooting.

## Coverage

The server registers tools across the full current `elation.Client` service surface, including:

- allergy/allergy documentation
- appointments
- billing
- clinical documents
- contacts
- medications, discontinued medications, prescription/history fills
- insurance companies, plans, policies, and eligibility
- letters and notes (visit and non-visit)
- patients, physicians, practices, pharmacies
- messaging threads and thread members
- recurring event groups and service locations
- subscriptions
