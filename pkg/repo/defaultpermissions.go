// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
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
		fmt.Fprintf(os.Stdout, `::error title="Elevated default actions permissions"::Elevated permissions in %s/%s%s`, org, repo, "\n")
	}

	// Check whether workflows can approve PRs.
	// TODO(mattmoor): We need to figure out how to disable checks for
	// repos, since the advisory repos approve PRs from actions.
	if dwp.GetCanApprovePullRequestReviews() {
		fmt.Fprintf(os.Stdout, `::error title="Actions can approve PRs"::Action approvers in %s/%s%s`, org, repo, "\n")
	}
	return nil
}
