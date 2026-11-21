// Copyright 2026 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"strings"
	"testing"

	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

// TestCommandsDocumentPermissions checks that every audit command tells users
// which token permissions it needs, with a sample Octo STS permissions block.
func TestCommandsDocumentPermissions(t *testing.T) {
	var leaves []*cobra.Command
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		if !cmd.HasAvailableSubCommands() {
			leaves = append(leaves, cmd)
			return
		}
		for _, sub := range cmd.Commands() {
			walk(sub)
		}
	}
	walk(New(github.NewClient(nil)))

	if len(leaves) == 0 {
		t.Fatal("found no audit commands")
	}
	for _, cmd := range leaves {
		t.Run(cmd.CommandPath(), func(t *testing.T) {
			if !strings.Contains(cmd.Long, "\n  permissions:\n") {
				t.Errorf("%q help does not include a sample Octo STS permissions block:\n%s", cmd.CommandPath(), cmd.Long)
			}
		})
	}
}
