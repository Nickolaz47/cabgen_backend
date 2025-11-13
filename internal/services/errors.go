package services

import "errors"

// Generic errors
var ErrNotFound = errors.New("record not found")
var ErrInternal = errors.New("internal system error")
var ErrConflict = errors.New("record already exists")