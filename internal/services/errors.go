package services

import "errors"

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
var ErrInvalidTicketStatus = errors.New("invalid ticket status")
var ErrDeleteActiveTicket = errors.New("cannot delete active ticket")
