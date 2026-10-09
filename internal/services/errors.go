package services

import (
	"errors"
	"fmt"
	"strings"
)

// Generic errors
var ErrNotFound = errors.New("record not found")
var ErrInternal = errors.New("internal system error")
var ErrConflict = errors.New("record already exists")
var ErrUnauthorized = errors.New("unauthorized")

// Specific errors
var ErrInvalidCountryCode = errors.New("invalid country code")
var ErrConflictEmail = errors.New("email already exists")
var ErrConflictUsername = errors.New("username already exists")
var ErrEmailMismatch = errors.New("emails must match")
var ErrPasswordMismatch = errors.New("passwords must match")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDisabledUser = errors.New("disabled user")
var ErrOriginNotFound = errors.New("origin not found")
var ErrUserNotFound = errors.New("user not found")
var ErrSampleSourceNotFound = errors.New("sample source not found")
var ErrMicroorganismNotFound = errors.New("microorganism not found")
var ErrSequencerNotFound = errors.New("sequencer not found")
var ErrLaboratoryNotFound = errors.New("laboratory not found")
var ErrHealthServiceNotFound = errors.New("health service not found")
var ErrMissingFiles = errors.New("missing files")
var ErrMissingFastq1 = errors.New("missing fastq1 file")
var ErrMissingFastq2 = errors.New("missing fastq2 file")
var ErrCreateFolder = errors.New("cannot create folder")
var ErrDeleteRunningAnalysis = errors.New("cannot delete analysis")
var ErrSampleNotFound = errors.New("sample not found")
var ErrExceededDownloadLimit = errors.New("exceeded download limit")
var ErrFastQCDownload = errors.New("FASTQC analysis cannot be downloaded")
var ErrNoDashboardData = errors.New("no analyses for in-network samples")
var ErrZipNotFound = errors.New("zip file not available for this analysis")
var ErrTicketAlreadyResolvedStatus = errors.New("ticket already resolved")
var ErrTicketIsNotOpen = errors.New("ticket in progress or resolved")
var ErrDeleteActiveTicket = errors.New("cannot delete active ticket")
var ErrUserPartOfNetwork = errors.New("user is part of the network")
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token is expired")
var ErrCurrentPasswordMismatch = errors.New("current password is incorrect")
var ErrInvalidEmailUpdateToken = errors.New("invalid email update token")
var ErrExpiredEmailUpdateToken = errors.New("email update token is expired")
var ErrEmailSame = errors.New("new email is the same as current email")
var ErrDuplicateTask = errors.New("duplicate task already pending")
var ErrInvalidStatusTransition = errors.New("invalid status transition")
var ErrAnalysisSampleLimit = errors.New("sample already has 2 analyses")

// Template table errors
var ErrInvalidTable = errors.New("invalid template table")
var ErrEmptyTable = errors.New("template table has no rows")

type TableHeadersError struct {
	Invalid []string
}

func (e *TableHeadersError) Error() string {
	return "missing or invalid table headers: " +
		strings.Join(e.Invalid, ", ")
}

type TableValueError struct {
	Row int
	Col string
}

func (e *TableValueError) Error() string {
	return fmt.Sprintf("invalid value at row %d, column %s", e.Row, e.Col)
}
