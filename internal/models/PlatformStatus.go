package models

import "slices"

type PlatformStatus string

const (
	PlatformStatusOk                PlatformStatus = "ok"
	PlatformStatusPlatformInTesting PlatformStatus = "platformInTesting"
	PlatformStatusMaintenance       PlatformStatus = "maintenance"
	PlatformStatusAdminOnly         PlatformStatus = "adminOnly"
	PlatformStatusAnalysisDisabled  PlatformStatus = "analysisDisabled"
)

var platformStatusValues = []PlatformStatus{
	PlatformStatusOk,
	PlatformStatusPlatformInTesting,
	PlatformStatusMaintenance,
	PlatformStatusAdminOnly,
	PlatformStatusAnalysisDisabled,
}

func ParsePlatformStatus(s string) PlatformStatus {
	ps := PlatformStatus(s)
	if slices.Contains(platformStatusValues, ps) {
		return ps
	}
	return PlatformStatusOk
}

func (ps PlatformStatus) String() string {
	return string(ps)
}
