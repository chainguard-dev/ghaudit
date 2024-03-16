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
	errDefaultPermissions  = gherror.New("Elevated default actions permissions")
	errApprovePullRequests = gherror.New("Actions can approve PRs")
)

func defaultPermissions(ghc *github.Client, org, repo *string) *cobra.Command {
	return &cobra.Command{
		Use:           "default-permissions",
		Short:         "Audit the default permissions.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DefaultPermissions(cmd.Context(), ghc, *org, *repo)
		},
	}
}

func DefaultPermissions(ctx context.Context, ghc *github.Client, org, repo string) error {
	dwp, _, err := ghc.Repositories.GetDefaultWorkflowPermissions(ctx, org, repo)
	if err != nil {
		return err
	}

	// Check whether the default workflow permissions are write.
	if dwp.GetDefaultWorkflowPermissions() == "write" {
		errDefaultPermissions.Emit("Elevated permissions in %s/%s", org, repo)
	}

	// Check whether workflows can approve PRs.
	// TODO(mattmoor): We need to figure out how to disable checks for
	// repos, since the advisory repos approve PRs from actions.
	if dwp.GetCanApprovePullRequestReviews() {
		errApprovePullRequests.Emit("Action approvers in %s/%s", org, repo)
	}
	return nil
}
