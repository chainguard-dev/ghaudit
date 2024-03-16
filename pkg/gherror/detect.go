// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package gherror

import "sync/atomic"

var errors = atomic.Bool{}

// HadErrors returns true if any errors have been emitted.
func HadErrors() bool {
	return errors.Load()
}

func sawError() {
	errors.Store(true)
}
