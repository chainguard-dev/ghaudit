// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/chainguard-dev/ghaudit/pkg/repo"
	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

func defaultPermissions(ghc *github.Client, org *string) *cobra.Command {
	rm := NewRepoMapper("default-permissions", ghc, org, repo.DefaultPermissions)

	return &cobra.Command{
		Use:           "default-permissions",
		Short:         "Audit the default permissions.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rm.Execute(cmd.Context())
		},
	}
}
