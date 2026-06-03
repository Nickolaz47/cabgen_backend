package data

import (
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
)

var baseAdminAnalysisCreateBody = map[string]any{
	"type":      "wgs",
	"sample_id": validUUID,
	"user_id":   validUUID,
}

var AdminAnalysisCreateTests = []Body{
	{"Missing type", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAdminAnalysisCreateBody)
		delete(b, "type")
		return b
	}()), `{"error":"The analysis type is required."}`},
	{"Missing sample_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAdminAnalysisCreateBody)
		delete(b, "sample_id")
		return b
	}()), `{"error":"The sample ID is required."}`},
	{"Missing user_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAdminAnalysisCreateBody)
		delete(b, "user_id")
		return b
	}()), `{"error":"The user ID is required."}`},
}

var baseAnalysisCreateBody = map[string]any{
	"type":      "wgs",
	"sample_id": validUUID,
}

var AnalysisCreateTests = []Body{
	{"Missing type", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAnalysisCreateBody)
		delete(b, "type")
		return b
	}()), `{"error":"The analysis type is required."}`},
	{"Missing sample_id", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAnalysisCreateBody)
		delete(b, "sample_id")
		return b
	}()), `{"error":"The sample ID is required."}`},
}

var baseAdminAnalysisUpdateBody = map[string]any{
	"status":           "completed",
	"metrics":          map[string]any{"coverage": 98.5, "reads": 1500000},
	"fastqc1":          "/app/uploads/fastqc1_report.html",
	"fastqc2":          "/app/uploads/fastqc2_report.html",
	"results_zip_path": "/app/results/analysis_results.zip",
	"error_message":    nil,
}

var AdminAnalysisUpdateTests = []Body{
	{"Invalid status", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAdminAnalysisUpdateBody)
		b["status"] = "invalid_status"
		return b
	}()), `{"error":"The analysis status is invalid."}`},
}

var baseAnalysisTSVDownloadBody = map[string]any{
	"ids": []string{validUUID, validUUID},
}

var AnalysisTSVDownloadTests = []Body{
	{"Missing ids", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAnalysisTSVDownloadBody)
		delete(b, "ids")
		return b
	}()), `{"error":"IDs are required to download the result tables."}`},
	{"Empty ids list", testutils.ToJSON(func() map[string]any {
		b := testutils.CopyMap(baseAnalysisTSVDownloadBody)
		b["ids"] = []string{}
		return b
	}()), `{"error":"At least one ID is required to download the result table."}`},
}
