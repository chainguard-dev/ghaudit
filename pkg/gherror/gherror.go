// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package gherror

import (
	"fmt"
	"os"
)

type Error interface {
	Emit(msg string, args ...interface{})
}

func New(title string) Error {
	return &errorImpl{title: title}
}

type errorImpl struct {
	title string
}

func (e *errorImpl) Emit(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, `::error title=%s::%s%s`, e.title, fmt.Sprintf(msg, args...), "\n")
}
