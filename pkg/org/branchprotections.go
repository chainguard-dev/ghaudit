// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/chainguard-dev/ghaudit/pkg/repo"
	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

func branchProtections(ghc *github.Client, org *string) *cobra.Command {
	rm := NewRepoMapper("branch-protections", ghc, org, repo.BranchProtections)

	return &cobra.Command{
		Use:           "branch-protections",
		Short:         "Audit the branch protections.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rm.Execute(cmd.Context())
		},
	}
}
