package models_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditToResponse(t *testing.T) {
	audit := testmodels.NewAudit(models.AuditEventLogin, "10.0.0.1",
		`{"email":"user@example.com"}`, 200)

	expected := models.AuditResponse{
		ID:        audit.ID,
		Event:     audit.Event,
		Source:    audit.Source,
		Status:    audit.Status,
		Metadata:  *audit.Metadata,
		CreatedAt: audit.CreatedAt.Format(time.RFC3339),
		Username:  audit.User.Username,
	}
	result := audit.ToResponse()

	assert.Equal(t, expected, result)
}

func TestAuditToResponseWithoutUser(t *testing.T) {
	audit := testmodels.NewAudit(models.AuditEventLoginFailed, "10.0.0.1",
		"{}", 401)
	audit.User = nil

	result := audit.ToResponse()

	assert.Empty(t, result.Username)
}

func TestAuditFilterBinding(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected models.AuditFilter
		wantErr  bool
	}{
		{
			name: "Success - Full filter",
			query: "event=auth.login&source=10.0.0.1&status=404&date=2026-01-02&user=" +
				"123e4567-e89b-12d3-a456-426614174000",
			expected: models.AuditFilter{
				Event:  "auth.login",
				Source: "10.0.0.1",
				Status: 404,
				Date: func() *time.Time {
					d := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
					return &d
				}(),
				UserID: func() *uuid.UUID {
					id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
					return &id
				}(),
			},
		},
		{
			name:     "Success - Empty query",
			query:    "",
			expected: models.AuditFilter{},
		},
		{
			name:    "Error - Invalid date",
			query:   "date=02/01/2026",
			wantErr: true,
		},
		{
			name:    "Error - Invalid user ID",
			query:   "user=abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := testutils.BindFilter[models.AuditFilter](tt.query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.expected, filter)
			}
		})
	}
}

func TestAuditEventsInSync(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "Audit.go", nil,
		parser.SkipObjectResolution)
	assert.NoError(t, err)

	constNames := make(map[string]bool)
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range valueSpec.Names {
				if strings.HasPrefix(name.Name, "AuditEvent") {
					constNames[name.Name] = true
				}
			}
		}
	}

	sliceNames := make(map[string]bool)
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok || len(valueSpec.Names) != 1 ||
				valueSpec.Names[0].Name != "AuditEvents" {
				continue
			}

			composite, ok := valueSpec.Values[0].(*ast.CompositeLit)
			assert.True(t, ok, "AuditEvents should be a slice literal")

			for _, element := range composite.Elts {
				ident, ok := element.(*ast.Ident)
				assert.True(t, ok,
					"AuditEvents elements must be constant identifiers")
				sliceNames[ident.Name] = true
			}
		}
	}

	assert.NotEmpty(t, constNames, "audit constants should not be empty")
	assert.NotEmpty(t, sliceNames, "AuditEvents slice should not be empty")

	var missingInSlice []string
	for name := range constNames {
		if !sliceNames[name] {
			missingInSlice = append(missingInSlice, name)
		}
	}
	sort.Strings(missingInSlice)

	assert.Empty(t, missingInSlice,
		"every AuditEvent constant must be present in AuditEvents: %s",
		strings.Join(missingInSlice, ", "))

	var unknownInSlice []string
	for name := range sliceNames {
		if !constNames[name] {
			unknownInSlice = append(unknownInSlice, name)
		}
	}
	sort.Strings(unknownInSlice)

	assert.Empty(t, unknownInSlice,
		"AuditEvents must not reference unknown constants: %s",
		strings.Join(unknownInSlice, ", "))
}

func TestAuditEventSelectOptions(t *testing.T) {
	opts := models.AuditEventSelectOptions()

	assert.Len(t, opts, len(models.AuditEvents))

	seen := make(map[string]bool)
	for _, opt := range opts {
		assert.Equal(t, opt.Value, opt.Label)

		assert.False(t, seen[opt.Value],
			"duplicated event in select options: %s", opt.Value)
		seen[opt.Value] = true
	}
}
