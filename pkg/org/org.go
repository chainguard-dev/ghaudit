// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
)

func New(ghc *github.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "org",
		Short:         "Commands to audit github organizations.",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	var org string
	cmd.PersistentFlags().StringVarP(&org, "organization", "o", "", "organization to perform audits on.")

	// Add sub-commands.
	cmd.AddCommand(
		deployKeys(ghc, &org),
		defaultPermissions(ghc, &org),
		branchProtections(ghc, &org),
	)

	return cmd
}
