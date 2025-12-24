# Publishing a New Version of the Go Module

This guide explains how to publish a new version of the Ujeebu Go SDK.

## Prerequisites

- Git repository with a remote origin configured
- All changes committed and pushed to the main branch
- Tests passing (`make test`)

## Quick Release

Use the Makefile for easy version bumping:

```bash
# Bump patch version (v0.1.5 → v0.1.6)
make release-patch

# Bump minor version (v0.1.5 → v0.2.0)
make release-minor

# Bump major version (v0.1.5 → v1.0.0)
make release-major
```

## Manual Release

### 1. Check Current Version

```bash
make version
# or
git describe --tags --abbrev=0
```

### 2. Ensure Code Quality

```bash
# Run tests
make test

# Format code
make fmt

# Build to verify compilation
make build
```

### 3. Commit All Changes

```bash
git add -A
git commit -m "Your commit message"
git push origin main
```

### 4. Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag
git push origin v1.0.0
```

## Semantic Versioning

Follow [semver](https://semver.org/) conventions:

| Version | When to Use |
|---------|-------------|
| **Patch** (v1.0.X) | Bug fixes, minor changes that don't affect API |
| **Minor** (v1.X.0) | New features, backward-compatible changes |
| **Major** (vX.0.0) | Breaking changes, API modifications |

## Major Version 2+

For v2 and above, Go modules require updating the module path:

1. Update `go.mod`:
   ```go
   module github.com/ujeebu/ujeebu-go-sdk/v2
   ```

2. Update all internal imports to include `/v2`

3. Tag with the new major version:
   ```bash
   git tag -a v2.0.0 -m "Release v2.0.0"
   git push origin v2.0.0
   ```

## After Publishing

### Verify on pkg.go.dev

After pushing a tag, the module will be available at:
- https://pkg.go.dev/github.com/ujeebu/ujeebu-go-sdk

The Go module proxy typically indexes new versions within a few minutes.

### Force Proxy Update (if needed)

```bash
GOPROXY=proxy.golang.org go list -m github.com/ujeebu/ujeebu-go-sdk@v1.0.0
```

## Users Can Install

```bash
# Latest version
go get github.com/ujeebu/ujeebu-go-sdk@latest

# Specific version
go get github.com/ujeebu/ujeebu-go-sdk@v1.0.0
```

## Troubleshooting

### Tag Not Showing Up

1. Ensure the tag is pushed: `git push origin --tags`
2. Wait a few minutes for the proxy to index

### Wrong Version Published

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Create new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Clear Module Cache

```bash
go clean -modcache
```

