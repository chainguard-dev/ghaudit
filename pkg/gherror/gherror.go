// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package gherror

import (
	"fmt"
	"os"
)

// Error is an interface for emitting errors to GitHub Actions.
type Error interface {
	// Emit emits an error message to GitHub Actions.
	Emit(msg string, args ...interface{})
}

// New creates a new class of GitHub Action errors with the given title.
func New(title string) Error {
	return &errorImpl{title: title}
}

type errorImpl struct {
	title string
}

// Emit implements the Error interface.
func (e *errorImpl) Emit(msg string, args ...interface{}) {
	sawError()
	_, _ = fmt.Fprintf(os.Stdout, `::error title=%s::%s%s`, e.title, fmt.Sprintf(msg, args...), "\n")
}
