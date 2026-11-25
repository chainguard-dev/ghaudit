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
		Use:   "default-permissions",
		Short: "Audit the default permissions.",
		Long: `Audit the default permissions.

This checks every repository in the organization that is not archived and
does not have the "no-ghaudit" or "no-ghaudit-default-permissions" topic.
Listing the repositories requires metadata access, and reading their default
GitHub Actions workflow permissions requires read access to their
administration settings.

Use an organization-level Octo STS trust policy, checked into the
organization's .github repository, and leave out "repositories" so that the
token can access every repository. A sample Octo STS permissions block:

  permissions:
    metadata: read
    administration: read`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rm.Execute(cmd.Context())
		},
	}
}
