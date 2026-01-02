# Release Guide

This document explains how to create releases for S3 Bucket Mirror without manually building binaries.

## Automated Release Methods

### Method 1: GitHub Actions (Recommended)

The easiest way to create releases with automatic binary builds for multiple platforms.

#### Setup

1. **Install GoReleaser configuration**

   ```bash
   # Files are already created:
   # .goreleaser.yaml
   # .github/workflows/release.yml
   ```

2. **Commit and push to GitHub**

   ```bash
   git add .
   git commit -m "Add release automation"
   git push origin main
   ```

#### Creating a Release

1. **Create and push a git tag**

   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **Automatic build** - GitHub Actions will:
   - Run tests
   - Build binaries for all platforms:
     - Linux (amd64, arm64, arm)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Create archives (.tar.gz, .zip)
   - Generate checksums
   - Create GitHub Release
   - Upload all artifacts

3. **Find your release** at `https://github.com/michaelahli/s3-mirror/releases`

### Method 2: Local GoReleaser

For testing releases locally before pushing to GitHub.

#### Setup

1. **Install GoReleaser**

   ```bash
   # macOS
   brew install goreleaser
   
   # Linux
   curl -sfL https://goreleaser.com/static/run | bash
   
   # Or with Go
   go install github.com/goreleaser/goreleaser@latest
   ```

#### Creating a Snapshot Release (Testing)

```bash
# Build without publishing
goreleaser release --snapshot --clean

# Binaries will be in ./dist/ directory
ls -la dist/
```

#### Creating a Real Release

```bash
# Requires git tag
git tag v0.1.0

# Build and create GitHub release
export GITHUB_TOKEN="your_github_token"
goreleaser release --clean
```

### Method 3: Using Makefile

Simplified commands for common tasks.

```bash
# View all available commands
make help

# Run tests
make test

# Build locally
make build

# Create snapshot release (local testing)
make snapshot

# Create real release (requires git tag)
make release
```

## Release Workflow

### 1. Prepare Release

```bash
# Update version in relevant files if needed
# Update CHANGELOG.md or release notes

# Commit changes
git add .
git commit -m "Prepare v0.1.0 release"
git push origin main
```

### 2. Create Tag

```bash
# Create annotated tag with message
git tag -a v0.1.0 -m "Release v0.1.0: Initial release with S3 and MinIO support"

# Or simple tag
git tag v0.1.0

# Push tag to trigger GitHub Actions
git push origin v0.1.0
```

### 3. Monitor Build

```bash
# Watch GitHub Actions
# Go to: https://github.com/michaelahli/s3-mirror/actions

# Or check locally with goreleaser
goreleaser release --snapshot --clean
```

### 4. Verify Release

1. Check GitHub Releases page
2. Download binaries for different platforms
3. Test binaries work correctly
4. Update documentation if needed

## Versioning

Follow [Semantic Versioning](https://semver.org/):

- `v0.1.0` - Initial release
- `v0.1.1` - Bug fixes
- `v0.2.0` - New features (backward compatible)
- `v1.0.0` - First stable release
- `v2.0.0` - Breaking changes

## What Gets Built

GoReleaser automatically builds:

### Binaries

- **Linux**: amd64, arm64, arm7
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

### Archives

- **Linux/macOS**: `.tar.gz`
- **Windows**: `.zip`
- Each includes: binary, README.md, LICENSE, config.yaml.example

### Additional Outputs

- `checksums.txt` - SHA256 checksums for all files
- Source code archives (automatic by GitHub)
- Docker images (optional, if configured)

## Platform-Specific Downloads

Users can download from GitHub Releases:

```bash
# Linux amd64
wget https://github.com/michaelahli/s3-mirror/releases/download/v0.1.0/s3-mirror_0.1.0_Linux_x86_64.tar.gz

# macOS arm64 (Apple Silicon)
wget https://github.com/michaelahli/s3-mirror/releases/download/v0.1.0/s3-mirror_0.1.0_Darwin_arm64.tar.gz

# Windows amd64
wget https://github.com/michaelahli/s3-mirror/releases/download/v0.1.0/s3-mirror_0.1.0_Windows_x86_64.zip
```

## Docker Images (Optional)

If you enable Docker builds in `.goreleaser.yaml`:

```bash
# Pull image
docker pull michaelahli/s3-mirror:v0.1.0

# Run with config
docker run -v $(pwd)/config.yaml:/config/config.yaml \
  michaelahli/s3-mirror:v0.1.0
```

## Troubleshooting

### GitHub Actions fails

```bash
# Check GitHub Actions logs
# Common issues:
# - Missing GITHUB_TOKEN (should be automatic)
# - Test failures
# - Invalid goreleaser config
```

### Local release fails

```bash
# Check goreleaser config
goreleaser check

# Test without publishing
goreleaser release --snapshot --clean

# Common issues:
# - Missing git tag
# - Dirty git working directory
# - Missing GITHUB_TOKEN environment variable
```

### Binary doesn't work on target platform

```bash
# Verify GOOS/GOARCH in .goreleaser.yaml
# Test with snapshot build first
make snapshot
```

## Quick Reference

```bash
# Complete release workflow
git add .
git commit -m "Your changes"
git push origin main
git tag v0.1.0
git push origin v0.1.0
# Wait for GitHub Actions to complete

# Test locally before releasing
make snapshot
./dist/s3-mirror_linux_amd64_v1/s3-mirror -version

# Manual release
export GITHUB_TOKEN="your_token"
git tag v0.1.0
make release
```

## Advanced: Pre-release and Draft Releases

### Pre-release (Beta, RC)

```bash
# Tag with pre-release suffix
git tag v0.1.0-beta.1
git push origin v0.1.0-beta.1
# GitHub will mark as "pre-release" automatically
```

### Draft Release

Edit `.goreleaser.yaml`:

```yaml
release:
  draft: true  # Change to true
```

Then manually publish the draft from GitHub Releases page.

## Resources

- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
