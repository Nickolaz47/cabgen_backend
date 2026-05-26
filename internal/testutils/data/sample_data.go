package data

import (
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

const validUUID = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"

var baseSampleCreateBody = map[string]any{
	"name":              "Sample-SARS-CoV-2",
	"collection_date":   "2026-05-20T00:00:00Z",
	"run_number":        "RUN-2026-XYZ",
	"run_date":          "2026-05-25T00:00:00Z",
	"city":              "Maricá",
	"origin_code":       "BR-RJ-01",
	"gender":            "Male",
	"date_of_birth":     "1990-01-01T00:00:00Z",
	"country_code":      "BRA",
	"user_id":           validUUID,
	"origin_id":         validUUID,
	"sample_source_id":  validUUID,
	"microorganism_id":  validUUID,
	"sequencer_id":      validUUID,
	"laboratory_id":     validUUID,
	"health_service_id": validUUID,
}

var CreateSampleTests = []Body{
	{"Missing name", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "name")
		return b
	}()), `{"error":"Name is required."}`},
	{"Name too short", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		b["name"] = "Sa"
		return b
	}()), `{"error":"Name must be at least 3 characters long."}`},
	{"Name too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		b["name"] = strings.Repeat("A", 101)
		return b
	}()), `{"error":"Name must be at most 100 characters long."}`},
	{"Missing collection_date", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "collection_date")
		return b
	}()), `{"error":"The collection date is required."}`},
	{"Missing run_date", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "run_date")
		return b
	}()), `{"error":"The run date is required."}`},
	{"Missing run_number", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "run_number")
		return b
	}()), `{"error":"The run number is required."}`},
	{"Run number too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		b["run_number"] = strings.Repeat("1", 51)
		return b
	}()), `{"error":"The run number must have a maximum of 50 characters."}`},
	{"Origin code too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		b["origin_code"] = strings.Repeat("A", 256)
		return b
	}()), `{"error":"The origin code must have a maximum of 255 characters."}`},
	{"Missing country_code", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "country_code")
		return b
	}()), `{"error":"Country code is required."}`},
	{"Country code invalid length", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		b["country_code"] = "BR"
		return b
	}()), `{"error":"The country code must be 3 characters long."}`},
	{"Missing user_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "user_id")
		return b
	}()), `{"error":"The user ID is required."}`},
	{"Missing origin_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "origin_id")
		return b
	}()), `{"error":"The origin ID is required."}`},
	{"Missing sample_source_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "sample_source_id")
		return b
	}()), `{"error":"The sample source ID is required."}`},
	{"Missing microorganism_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "microorganism_id")
		return b
	}()), `{"error":"The microorganism ID is required."}`},
	{"Missing sequencer_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "sequencer_id")
		return b
	}()), `{"error":"The sequencer ID is required."}`},
	{"Missing laboratory_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "laboratory_id")
		return b
	}()), `{"error":"The laboratory ID is required."}`},
	{"Missing health_service_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleCreateBody)
		delete(b, "health_service_id")
		return b
	}()), `{"error":"The health service ID is required."}`},
}

var baseSampleUpdateBody = map[string]any{
	"name":        "Updated-Sample-Name",
	"run_number":  "RUN-UPDATED-01",
	"origin_code": "BR-RJ-NEW",
	"gender":      "Female",
}

var UpdateSampleTests = []Body{
	{"Name too short on update", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleUpdateBody)
		b["name"] = "Ab"
		return b
	}()), `{"error":"Name must be at least 3 characters long."}`},

	{"Name too long on update", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleUpdateBody)
		b["name"] = strings.Repeat("A", 101)
		return b
	}()), `{"error":"Name must be at most 100 characters long."}`},

	{"Run number too long on update", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleUpdateBody)
		b["run_number"] = strings.Repeat("1", 51)
		return b
	}()), `{"error":"The run number must have a maximum of 50 characters."}`},

	{"Origin code too long on update", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleUpdateBody)
		b["origin_code"] = strings.Repeat("A", 256)
		return b
	}()), `{"error":"The origin code must have a maximum of 255 characters."}`},

	{"Invalid country_code len on update", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleUpdateBody)
		b["country_code"] = "BR"
		return b
	}()), `{"error":"The country code must be 3 characters long."}`},
}

var baseSampleAttachmentBody = map[string]any{
	"fastq1": "/app/uploads/r1.fastq.gz",
	"fastq2": "/app/uploads/r2.fastq.gz",
	"fasta":  nil,
}

var AttachmentSampleTests = []Body{
	{"Fastq1 path too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleAttachmentBody)
		b["fastq1"] = strings.Repeat("A", 256)
		return b
	}()), `{"error":"The fastq1 name must have a maximum of 255 characters."}`},

	{"Fastq2 path too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleAttachmentBody)
		b["fastq2"] = strings.Repeat("A", 256)
		return b
	}()), `{"error":"The fastq2 name must have a maximum of 255 characters."}`},

	{"Fasta path too long", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseSampleAttachmentBody)
		b["fasta"] = strings.Repeat("A", 256)
		return b
	}()), `{"error":"The fasta name must have a maximum of 255 characters."}`},
}
