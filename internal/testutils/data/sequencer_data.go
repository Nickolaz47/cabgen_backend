package data

import "github.com/CABGenOrg/cabgen_backend/internal/testutils"

var baseSequencerCreateBody = map[string]any{
	"brand":     "Illumina",
	"model":     "MiSeq",
	"is_active": true,
}

var CreateSequencerTests = []Body{
	{"Missing brand", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["brand"] = ""; return b }()), `{"error":"Sequencer brand is required."}`},
	{"Missing model", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["model"] = ""; return b }()), `{"error":"Sequencer model is required."}`},
	{"Brand too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["brand"] = "I"; return b }()), `{"error":"Sequencer brand must contain at least 3 characters."}`},
	{"Model too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["model"] = "Mi"; return b }()), `{"error":"Sequencer model must contain at least 3 characters."}`},
}

var UpdateSequencerTests = []Body{
	{"Brand too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["brand"] = "I"; return b }()), `{"error":"Sequencer brand must contain at least 3 characters."}`},
	{"Model too short", testutils.ToJSON(func() map[string]any { b := testutils.CopyMap(baseSequencerCreateBody); b["model"] = "Mi"; return b }()), `{"error":"Sequencer model must contain at least 3 characters."}`},
}
