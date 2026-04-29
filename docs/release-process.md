# Release Process

This document describes the process for releasing new versions of `daraja-client-go`.

## Versioning

This project follows [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** — breaking changes to the public API
- **MINOR** — backward-compatible new functionality
- **PATCH** — backward-compatible bug fixes

## Prerequisites

- Write access to the repository
- [GitHub CLI](https://cli.github.com/) (`gh`) installed
- Go toolchain matching `go.mod`

## Release Steps

1. **Ensure `main` is green**
   - All CI checks pass on `main`
   - No unresolved high/critical security findings

2. **Update `CHANGELOG.md`**
   - Review merged PRs since last release
   - Add entries under the new version header
   - Follow [Keep a Changelog](https://keepachangelog.com/) format

3. **Tag the release**
   ```bash
   git checkout main
   git pull origin main
   git tag -a v1.0.0 -m "v1.0.0"
   git push origin v1.0.0
   ```

4. **Create GitHub Release**
   ```bash
   gh release create v1.0.0 \
     --title "v1.0.0" \
     --notes "See CHANGELOG.md for details" \
     --verify-tag
   ```

5. **Verify module availability**
   ```bash
   go list -m github.com/Kisotu/daraja-client-go@v1.0.0
   ```

## Hotfix Release

For critical bug fixes, create a branch from the release tag:

```bash
git checkout -b hotfix/v1.0.1 v1.0.0
# Apply fix
git tag -a v1.0.1 -m "v1.0.1"
git push origin v1.0.1
```