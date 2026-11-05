// Copyright 2026 Chainguard, Inc.
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v72/github"
)

func TestHasBypassAllowances(t *testing.T) {
	tests := []struct {
		name string
		in   *github.BypassPullRequestAllowances
		want bool
	}{
		{name: "nil", in: nil, want: false},
		{name: "empty lists", in: &github.BypassPullRequestAllowances{
			Users: []*github.User{}, Teams: []*github.Team{}, Apps: []*github.App{},
		}, want: false},
		{name: "zero value", in: &github.BypassPullRequestAllowances{}, want: false},
		{name: "user", in: &github.BypassPullRequestAllowances{
			Users: []*github.User{{Login: github.Ptr("octocat")}},
		}, want: true},
		{name: "team", in: &github.BypassPullRequestAllowances{
			Teams: []*github.Team{{Slug: github.Ptr("admins")}},
		}, want: true},
		{name: "app", in: &github.BypassPullRequestAllowances{
			Apps: []*github.App{{Slug: github.Ptr("dependabot")}},
		}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasBypassAllowances(tt.in); got != tt.want {
				t.Errorf("hasBypassAllowances() = %v, want %v", got, tt.want)
			}
		})
	}
}

// protectionJSON is a branch protection response that passes every check,
// with the bypass allowances supplied by the caller.
func protectionJSON(bypass string) string {
	return `{
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true,
    "required_approving_review_count": 1,
    "bypass_pull_request_allowances": ` + bypass + `
  },
  "enforce_admins": {"enabled": true},
  "required_status_checks": {"strict": true, "contexts": ["test"], "checks": [{"context": "test"}]}
}`
}

func TestBranchProtectionsBypass(t *testing.T) {
	tests := []struct {
		name       string
		bypass     string
		wantBypass bool
	}{
		{
			name:       "empty allowances are not a finding",
			bypass:     `{"users": [], "teams": [], "apps": []}`,
			wantBypass: false,
		},
		{
			name:       "app allowance is a finding",
			bypass:     `{"users": [], "teams": [], "apps": [{"id": 29110, "slug": "dependabot"}]}`,
			wantBypass: true,
		},
		{
			name:       "user allowance is a finding",
			bypass:     `{"users": [{"login": "octocat"}], "teams": [], "apps": []}`,
			wantBypass: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("GET /repos/org/repo", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, `{"default_branch": "main"}`)
			})
			mux.HandleFunc("GET /repos/org/repo/branches/main/protection", func(w http.ResponseWriter, _ *http.Request) {
				_, _ = io.WriteString(w, protectionJSON(tt.bypass))
			})
			srv := httptest.NewServer(mux)
			t.Cleanup(srv.Close)

			ghc := github.NewClient(nil)
			u, err := url.Parse(srv.URL + "/")
			if err != nil {
				t.Fatal(err)
			}
			ghc.BaseURL = u

			out := captureStdout(t, func() {
				if err := BranchProtections(context.Background(), ghc, "org", "repo"); err != nil {
					t.Fatalf("BranchProtections() = %v", err)
				}
			})

			if got := strings.Contains(out, "allows bypassing pull request reviews"); got != tt.wantBypass {
				t.Errorf("bypass finding = %v, want %v; output:\n%s", got, tt.wantBypass, out)
			}
			// The rest of the protection is compliant, so nothing else should fire.
			if !tt.wantBypass && out != "" {
				t.Errorf("unexpected findings:\n%s", out)
			}
		})
	}
}

// captureStdout runs fn and returns what it wrote to os.Stdout, which is
// where gherror emits findings.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()
	_ = w.Close()
	return <-done
}
