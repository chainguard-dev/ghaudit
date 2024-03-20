// Copyright 2024 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"

	"github.com/google/go-github/v60/github"
	"k8s.io/apimachinery/pkg/util/sets"
)

type RepoMapper interface {
	Execute(context.Context) error
}

type RepoFunc func(ctx context.Context, ghc *github.Client, org, repo string) error

func NewRepoMapper(name string, ghc *github.Client, org *string, rf RepoFunc) RepoMapper {
	return &repoMapper{
		name: name,
		ghc:  ghc,
		org:  org,
		rf:   rf,
	}
}

type repoMapper struct {
	name string
	ghc  *github.Client
	org  *string
	rf   RepoFunc
}

func (rm *repoMapper) Execute(ctx context.Context) error {
	page := 0
	for {
		repos, resp, err := rm.ghc.Repositories.ListByOrg(ctx, *rm.org, &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		})
		if err != nil {
			return err
		}

		for _, r := range repos {
			// Skip archived repositories.
			if r.GetArchived() {
				continue
			}
			// Skip repositories with the "no-ghaudit" topic.
			topics := sets.New[string](r.Topics...)
			if topics.Has("no-ghaudit") || topics.Has("no-ghaudit:"+rm.name) {
				continue
			}

			if err := rm.rf(ctx, rm.ghc, *rm.org, r.GetName()); err != nil {
				return err
			}
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return nil
}
