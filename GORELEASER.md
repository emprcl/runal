# GoReleaser Setup for Runal

This document explains how to use GoReleaser with the Runal project for building cross-platform binaries and distribution packages.

## Overview

GoReleaser replaces our previous custom GitHub Actions workflow and provides:

- Cross-platform binary builds (Linux, macOS, Windows)
- Multiple architectures (amd64, arm64)
- Linux distribution packages (deb, rpm, apk, archlinux)
- Homebrew formula generation

- Automated release management
- Changelog generation

## Local Development

### Prerequisites

1. Install GoReleaser:
   ```bash
   # macOS
   brew install goreleaser
   
   # Linux
   curl -sfL https://goreleaser.com/static/run | bash
   
   # Or download from: https://github.com/goreleaser/goreleaser/releases
   ```



### Testing Locally

1. **Dry run** (check configuration without building):
   ```bash
   goreleaser check
   ```

2. **Build snapshot** (local build without releasing):
   ```bash
   goreleaser build --snapshot --clean
   ```

3. **Full snapshot release** (build everything locally):
   ```bash
   goreleaser release --snapshot --clean
   ```

4. **Test specific targets**:
   ```bash
   # Build only for Linux
   goreleaser build --single-target --snapshot --clean
   
   # Build packages only
   goreleaser release --snapshot --clean --skip=publish
   ```

### Output Structure

After running GoReleaser, you'll find outputs in the `dist/` directory:

```
dist/
├── checksums.txt
├── config.yaml
├── metadata.json
├── runal_linux_amd64_v1/
│   └── runal
├── runal_darwin_amd64_v1/
│   └── runal
├── runal_windows_amd64_v1/
│   └── runal.exe
├── runal_*.tar.gz
├── runal_*.zip
├── runal_*.deb
├── runal_*.rpm
├── runal_*.apk
└── runal_*.pkg.tar.zst
```

## CI/CD Workflow

### GitHub Actions Integration

The workflow automatically:

1. **On Push to Main**: Creates snapshot builds for testing
2. **On Pull Request**: Creates snapshot builds for validation
3. **On Tag Push** (`v*`): Creates and publishes full release

### Release Process

1. **Create a new tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically**:
   - Build cross-platform binaries
   - Create distribution packages
   - Generate changelog
   - Create GitHub release with assets

   - Update Homebrew formula (if tap exists)

## Configuration Details

### Build Configuration

- **Target Platforms**: Linux, macOS, Windows
- **Architectures**: amd64, arm64
- **CGO**: Disabled for static binaries
- **LDFLAGS**: Strips debug info and injects version/commit/date

### Package Formats

- **deb**: Debian/Ubuntu packages
- **rpm**: RHEL/Fedora/CentOS packages  
- **apk**: Alpine Linux packages
- **archlinux**: Arch Linux packages



### Homebrew

- Tap: `emprcl/homebrew-tap` (needs to be created)
- Formula location: `Formula/runal.rb`
- Auto-updates on release

## Customization

### Version Information

GoReleaser injects build information into the binary:

```go
var (
    version = "dev"     // Git tag or "dev"
    commit  = "none"    // Git commit hash
    date    = "unknown" // Build timestamp
)
```

### Adding New Platforms

Edit `.goreleaser.yaml` to add new targets:

```yaml
builds:
  - goos:
      - linux
      - darwin 
      - windows
      - freebsd  # Add new OS
    goarch:
      - amd64
      - arm64
      - 386      # Add new architecture
```

### Package Metadata

Update package information in the `nfpms` section:

```yaml
nfpms:
  - package_name: runal
    vendor: emprcl
    homepage: https://github.com/emprcl/runal
    description: "Your package description"
    license: MIT
```

## Troubleshooting

### Common Issues

1. **Build fails on specific platform**:
   - Check Go version compatibility
   - Verify CGO requirements
   - Test locally with: `GOOS=target_os GOARCH=target_arch go build`

2. **Package validation fails**:
   - Install package locally and test
   - Check file permissions and paths
   - Verify package metadata



### Debug Commands

```bash
# Verbose output
goreleaser release --snapshot --clean -v

# Skip specific steps
goreleaser release --snapshot --clean --skip=publish

# Only run specific steps
goreleaser build --snapshot --clean --single-target
```

## Migration Notes

### Changes from Previous Workflow

1. **Removed**: Custom build matrix in GitHub Actions
2. **Removed**: Manual archive creation
3. **Added**: Native package generation
3. **Added**: Homebrew formula generation

### Version Handling

- **Before**: Embedded VERSION file
- **After**: Build-time injection via LDFLAGS

### Artifact Names

- **Before**: `runal_VERSION_OS.tar.gz`
- **After**: `runal_VERSION_OS_ARCH.tar.gz`

This provides better consistency and follows common naming conventions.