// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"log"
	"os"

	"github.com/chainguard-dev/ghaudit/pkg/gherror"
	"github.com/chainguard-dev/ghaudit/pkg/org"
	"github.com/chainguard-dev/ghaudit/pkg/repo"
	"github.com/google/go-github/v60/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func New(ghc *github.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ghaudit",
		Short:         "GitHub Audit",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	// Add sub-commands.
	cmd.AddCommand(
		org.New(ghc),
		repo.New(ghc),
	)

	return cmd
}

func main() {
	ctx := context.Background()

	tok, ok := os.LookupEnv("GH_TOKEN")
	if !ok {
		log.Fatal("GH_TOKEN must be set")
	}

	ghc := github.NewClient(
		oauth2.NewClient(ctx,
			oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: tok,
			}),
		),
	)

	cmd := New(ghc)

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}

	// Exit with a non-zero status code if there were any errors.
	if gherror.HadErrors() {
		os.Exit(1)
	}
}
