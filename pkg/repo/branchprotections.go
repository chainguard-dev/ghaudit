// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/chainguard-dev/ghaudit/pkg/gherror"
	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
)

var errBranchProtection = gherror.New("Default branch protection")

func branchProtections(ghc *github.Client, org, repo *string) *cobra.Command {
	return &cobra.Command{
		Use:           "branch-protections",
		Short:         "Audit the branch protections.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return BranchProtections(cmd.Context(), ghc, *org, *repo)
		},
	}
}

func BranchProtections(ctx context.Context, ghc *github.Client, org, repo string) error {
	r, _, err := ghc.Repositories.Get(ctx, org, repo)
	if err != nil {
		return err
	}

	if prot, _, err := ghc.Repositories.GetBranchProtection(ctx, org, repo, r.GetDefaultBranch()); err != nil {
		errBranchProtection.Emit("%s/%s branch %s returned %v", org, repo, r.GetDefaultBranch(), err)
	} else {
		// TODO(mattmoor): Check prot.GetRequiredPullRequestReviews().RequiredApprovingReviewCount
		// TODO(mattmoor): Check(?) prot.GetRequiredPullRequestReviews().RequireCodeOwnerReviews
		// TODO(mattmoor): Check prot.GetRequiredPullRequestReviews().GetBypassPullRequestAllowances()
		// TODO(mattmoor): Check prot.GetEnforceAdmins() and prot.GetEnforceAdmins().Enabled
		// TODO(mattmoor): Check that one of these is non-empty:
		//  - prot.GetRequiredStatusChecks().Contexts
		//  - prot.GetRequiredStatusChecks().Checks
		// TODO(mattmoor): Check prot.GetBlockCreations().Enabled (if the branch protection has a wildcard in it)

		// TODO(mattmoor): Look for others?

		_ = prot // Check the contents of the default branch protection.
	}
	return nil
}
