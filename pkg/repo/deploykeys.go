// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/chainguard-dev/ghaudit/pkg/gherror"
	"github.com/google/go-github/v72/github"
	"github.com/spf13/cobra"
)

var errDeployKeys = gherror.New("Found deploy keys")

func deployKeys(ghc *github.Client, org, repo *string) *cobra.Command {
	return &cobra.Command{
		Use:           "deploy-keys",
		Short:         "Audit for usage of deploy keys.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DeployKeys(cmd.Context(), ghc, *org, *repo)
		},
	}
}

func DeployKeys(ctx context.Context, ghc *github.Client, org, repo string) error {
	keys, _, err := ghc.Repositories.ListKeys(ctx, org, repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	// Check whether there are any deploy keys.
	// TODO(mattmoor): bump the severity if there are any non-readonly ones?
	if len(keys) > 0 {
		errDeployKeys.Emit("Deploy keys used in %s/%s", org, repo)
	}
	return nil
}
