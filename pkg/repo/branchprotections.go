// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/chainguard-dev/ghaudit/pkg/gherror"
	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
)

var (
	errBranchProtection = gherror.New("Default branch protection")
	errRequiredReviews  = gherror.New("Required reviews")
	errAdminBypass      = gherror.New("Admin bypass")
	errRequiredChecks   = gherror.New("Required status checks")
)

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
		if rev := prot.GetRequiredPullRequestReviews(); rev == nil {
			errRequiredReviews.Emit("%s/%s branch %s returned does not require pull request reviews", org, repo, r.GetDefaultBranch())
		} else {
			if rev.RequiredApprovingReviewCount == 0 {
				errRequiredReviews.Emit("%s/%s branch %s does not require pull request approval", org, repo, r.GetDefaultBranch())
			}
			if bypasses := rev.GetBypassPullRequestAllowances(); bypasses != nil {
				errRequiredReviews.Emit("%s/%s branch %s allows bypassing pull request reviews", org, repo, r.GetDefaultBranch())
			}
			if !rev.DismissStaleReviews {
				errRequiredReviews.Emit("%s/%s branch %s does not dismiss stale reviews when new commits are pushed", org, repo, r.GetDefaultBranch())
			}
			if !rev.RequireCodeOwnerReviews {
				errRequiredReviews.Emit("%s/%s branch %s does not require reviews from Code Owners", org, repo, r.GetDefaultBranch())
			}
			// TODO(mattmoor): How to check whether approval of the most recent reviewable push is enabled?
		}

		if admin := prot.GetEnforceAdmins(); admin == nil || !admin.Enabled {
			errAdminBypass.Emit("%s/%s branch %s allows admins to bypass branch protections", org, repo, r.GetDefaultBranch())
		}

		if checks := prot.GetRequiredStatusChecks(); checks == nil {
			errRequiredChecks.Emit("%s/%s branch %s does not require status checks", org, repo, r.GetDefaultBranch())
		} else if checks.Contexts != nil && len(*checks.Contexts) == 0 {
			// Either contexts of checks must be specified.
			errRequiredChecks.Emit("%s/%s branch %s does not have any required contexts", org, repo, r.GetDefaultBranch())
		} else if checks.Checks != nil && len(*checks.Checks) == 0 {
			// Either contexts of checks must be specified.
			errRequiredChecks.Emit("%s/%s branch %s does not have any required checks", org, repo, r.GetDefaultBranch())
		}

		// TODO(mattmoor): Check prot.GetBlockCreations().Enabled (if the branch protection has a wildcard in it)
	}
	return nil
}
