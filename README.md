# ghaudit

This repository contains an experimental CLI tool that checks for a variety of
github org and repository settings.

## Usage

`ghaudit` reads a GitHub token from `GITHUB_TOKEN` (or `GH_TOKEN`) and reports
each finding as a GitHub Actions error annotation. It exits with a non-zero
status if it found anything.

```shell
# Audit a single repository.
ghaudit repo --organization chainguard-dev --repository ghaudit branch-protections

# Audit every repository in an organization.
ghaudit org --organization chainguard-dev deploy-keys
```

The `org` commands skip archived repositories. Add the `no-ghaudit` topic to a
repository to skip it entirely, or a `no-ghaudit-<command>` topic (for example
`no-ghaudit-deploy-keys`) to skip a single check.

## Required permissions

Every command's `--help` output lists the permissions it needs. They are:

| Command               | What it reads                                   | Permissions                              |
|-----------------------|-------------------------------------------------|------------------------------------------|
| `deploy-keys`         | The repository's deploy keys                    | `metadata: read`, `administration: read` |
| `default-permissions` | The default GitHub Actions workflow permissions | `metadata: read`, `administration: read` |
| `branch-protections`  | The protection rules of the default branch      | `metadata: read`, `administration: read` |

The `org` variants run the same checks on every repository in the organization,
so the token must have access to all of them. Listing the repositories only
needs `metadata: read`.

To run `ghaudit` from GitHub Actions with [Octo STS](https://github.com/octo-sts/app),
check a trust policy like this into `.github/chainguard/ghaudit.sts.yaml` of the
repository to audit:

```yaml
issuer: https://token.actions.githubusercontent.com
subject: repo:my-org/my-repo:ref:refs/heads/main

permissions:
  metadata: read
  administration: read
```

To audit a whole organization, check the policy into the organization's
`.github` repository instead, leave out `repositories` so that the token can
access every repository in the organization, and federate the token with the
organization as its `scope`. See
[`presubmit-testing.yaml`](.github/workflows/presubmit-testing.yaml) for a
workflow that federates a token and runs every `repo` command.
