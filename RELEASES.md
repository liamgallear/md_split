# Creating Releases

This repository uses GitHub Actions to automatically create releases when a new version tag is pushed.

## How to Create a Release

1. Make sure all changes are committed and pushed to the main branch
2. Create and push a version tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The GitHub Actions workflow will automatically:
   - Build binaries for multiple platforms (Linux, macOS, Windows on AMD64 and ARM64)
   - Create a GitHub release with auto-generated release notes
   - Attach all built binaries as downloadable assets

## Version Tagging Convention

- Use semantic versioning (e.g., `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`)
- Always prefix tags with `v` (e.g., `v1.0.0` not `1.0.0`)
- The workflow is triggered by any tag that starts with `v`

## Built Artifacts

The release workflow creates the following binary artifacts:
- `md_split-linux-amd64`
- `md_split-linux-arm64`
- `md_split-darwin-amd64` (macOS Intel)
- `md_split-darwin-arm64` (macOS Apple Silicon)
- `md_split-windows-amd64.exe`
- `md_split-windows-arm64.exe`