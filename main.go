package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

var (
	org = flag.String("org", "", "organization to list repositories for")
)

func main() {
	flag.Parse()
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

	page := 0
	for {
		repos, cur, err := ghc.Repositories.ListByOrg(ctx, *org, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page: page,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			// Skip archived repositories.
			if repo.GetArchived() {
				continue
			}

			// Check for any Deploy Keys.
			{
				keys, _, err := ghc.Repositories.ListKeys(ctx, repo.Owner.GetLogin(), repo.GetName(), &github.ListOptions{})
				if err != nil {
					log.Fatal(err)
				}

				// Check whether there are any deploy keys.
				// TODO(mattmoor): bump the severity if there are any non-readonly ones?
				if len(keys) > 0 {
					fmt.Fprintf(os.Stdout, `::error title="Found deploy keys"::Deploy keys used in %s%s`, *repo.FullName, "\n")
				}
			}

			// Check the default workflow permissions.
			{
				dwp, _, err := ghc.Repositories.GetDefaultWorkflowPermissions(ctx, repo.Owner.GetLogin(), repo.GetName())
				if err != nil {
					log.Fatal(err)
				}

				// Check whether the default workflow permissions are write.
				if dwp.GetDefaultWorkflowPermissions() == "write" {
					fmt.Fprintf(os.Stdout, `::error title="Elevated default actions permissions"::Elevated permissions in %s%s`, *repo.FullName, "\n")
				}

				// Check whether workflows can approve PRs.
				// TODO(mattmoor): We need to figure out how to disable checks for
				// repos, since the advisory repos approve PRs from actions.
				if dwp.GetCanApprovePullRequestReviews() {
					fmt.Fprintf(os.Stdout, `::error title="Actions can approve PRs"::Action approvers in %s%s`, *repo.FullName, "\n")
				}
			}

			// Check the default branch protection.
			{
				if prot, _, err := ghc.Repositories.GetBranchProtection(ctx, repo.Owner.GetLogin(), repo.GetName(), repo.GetDefaultBranch()); err != nil {
					fmt.Fprintf(os.Stdout, `::error title="Default branch protection"::%s branch %s returned %v%s`, *repo.FullName, repo.GetDefaultBranch(), err, "\n")
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
			}

			// TODO(mattmoor): branch protections on other branches (requires "contents: read" permission)

			// TODO(mattmoor): enumerate secrets and check when they were last updated (org, repo, environment)
			// org secrets: organization_secrets:read
			// repo secrets: secrets:read
			// env secrets: environments:read and secrets:read

			// TODO(mattmoor): Check runner group visibility restrictions using:
			// ghc.Actions.ListOrganizationRunnerGroups()

			// TODO(mattmoor): Is there a way to check for the age of particular
			// runners, to assess ephemerality?
		}

		if page = cur.NextPage; page == 0 {
			break
		}
	}
}
