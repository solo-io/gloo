# Renovate Configuration for Gloo

This document describes the Renovate configuration for automated dependency management in the Gloo project.

## Overview

Renovate is configured to automatically:
- Update Go module dependencies
- Update Docker images in Dockerfiles and CI
- Update GitHub Actions
- Update Helm charts and values
- Create security vulnerability PRs immediately
- Group related dependencies for easier review

## Configuration Files

- `renovate.json` - Main configuration file
- `.github/renovate.json5` - GitHub-specific settings with enhanced formatting
- This documentation - `docs/content/guides/dev/renovate.md`

## Schedule

Renovate runs on the following schedule:
- **Regular updates**: Before 6am on Monday (Pacific Time)
- **Security updates**: Immediately when detected
- **Rate limiting**: Maximum 3 PRs per hour, 5 concurrent PRs

## Dependency Groups

Dependencies are grouped for easier review:

### Go Dependencies
- All Go module dependencies
- Minimum release age: 3 days
- Weekly updates on Monday

### Solo.io Dependencies
- All `github.com/solo-io/*` packages
- Minimum release age: 1 day
- Weekly updates on Monday

### Kubernetes Dependencies
- All `k8s.io/*` and `sigs.k8s.io/*` packages
- Minimum release age: 7 days
- Conservative update schedule

### Envoy Dependencies
- All Envoy-related packages
- Minimum release age: 7 days
- Conservative update schedule

### GitHub Actions
- All GitHub Actions updates
- Auto-merge enabled for patch updates
- Minimum release age: 3 days

### Docker Images
- Dockerfile and docker-compose updates
- Minimum release age: 3 days

## Auto-merge Policy

The following updates are automatically merged:
- GitHub Actions patch updates (after 3 days)
- Non-critical patch updates (excluding K8s, Envoy, Solo.io packages)

Critical dependencies require manual review.

## Security Updates

- Vulnerability alerts are enabled and processed immediately
- OSV (Open Source Vulnerabilities) database integration
- Security PRs bypass normal scheduling

## Branch and PR Naming

- Branch prefix: `renovate/`
- PR title format: `deps: Update dependency {name} to v{version}`
- Semantic commit messages enabled
- Labeled with `dependencies` and `renovate`

## Team Assignment

- Assignees: `@solo-io/gloo-maintainers`
- Reviewers: `@solo-io/gloo-maintainers`
- CI images: `@solo-io/gloo-ci`

## Customization

To modify Renovate behavior:

1. Edit `renovate.json` for general settings
2. Edit `.github/renovate.json5` for GitHub-specific settings
3. Use package rules to customize behavior for specific dependencies
4. Test changes in a separate branch before merging

## Monitoring

- Dependency Dashboard is enabled in GitHub Issues
- Check the Renovate logs in GitHub Actions for troubleshooting
- Monitor the `#gloo-dev` Slack channel for notifications

## Manual Override

To manually trigger Renovate:
1. Add a comment to any Renovate PR: `@renovatebot rebase`
2. Close and reopen a Renovate PR to trigger a rebase
3. Use the Dependency Dashboard to check boxes for immediate updates

## Troubleshooting

Common issues and solutions:

### Renovate not creating PRs
- Check if the dependency is excluded in `packageRules`
- Verify the minimum release age has passed
- Check rate limiting (max 3 PRs/hour)

### Auto-merge not working
- Ensure required status checks are passing
- Verify branch protection rules allow auto-merge
- Check if the package is excluded from auto-merge

### Too many PRs
- Adjust `prHourlyLimit` and `prConcurrentLimit`
- Use more aggressive grouping
- Increase `minimumReleaseAge` for stability

For more information, see the [Renovate documentation](https://docs.renovatebot.com/). 