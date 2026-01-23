package services

import "errors"

// Generic errors
var ErrNotFound = errors.New("record not found")
var ErrInternal = errors.New("internal system error")
var ErrConflict = errors.New("record already exists")

// Specific errors
var ErrInvalidCountryCode = errors.New("invalid country code")
var ErrConflictEmail = errors.New("email already exists")
var ErrConflictUsername = errors.New("username already exists")
var ErrEmailMismatch = errors.New("emails must match")
var ErrPasswordMismatch = errors.New("passwords must match")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrDisabledUser = errors.New("disabled user")
var ErrUnauthorized = errors.New("unauthorized")
