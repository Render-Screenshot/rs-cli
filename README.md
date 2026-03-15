# rs — RenderScreenshot CLI

Command-line interface for the [RenderScreenshot](https://renderscreenshot.com) API. Capture web screenshots, generate signed URLs, manage cache, and preview page metadata from your terminal.

## Installation

### Homebrew

```bash
brew install render-screenshot/tap/rs
```

### Go Install

```bash
go install github.com/Render-Screenshot/rs-cli/cmd/rs@latest
```

### Binary Download

Download from [GitHub Releases](https://github.com/Render-Screenshot/rs-cli/releases).

## Quick Start

```bash
# Authenticate
rs login

# Take a screenshot
rs take https://example.com -o screenshot.png

# Use a preset
rs take https://example.com --preset og_card

# JSON output
rs take https://example.com --json

# Preview page metadata (no API credits)
rs preview https://github.com

# Generate signed URL
rs signed-url https://example.com --expires 7d

# Batch screenshots
rs batch https://example.com https://github.com -d ./screenshots

# List presets and devices
rs presets
rs devices
```

## Authentication

```bash
# Interactive login
rs login

# With signed URL support
rs login --signed-urls

# Environment variable
export RS_API_KEY=rs_live_xxxxx

# Per-command flag
rs take https://example.com --api-key rs_live_xxxxx
```

Precedence: `--api-key` flag > `RS_API_KEY` env > config file.

## Commands

| Command | Description |
|---------|-------------|
| `rs take` | Capture a screenshot |
| `rs batch` | Screenshot multiple URLs |
| `rs signed-url` | Generate signed URLs |
| `rs preview` | Fetch page metadata |
| `rs cache` | Manage cached screenshots |
| `rs presets` | List screenshot presets |
| `rs devices` | List device presets |
| `rs whoami` | Show account info |
| `rs config` | Manage configuration |
| `rs login` | Authenticate |
| `rs logout` | Remove credentials |

## Global Flags

| Flag | Description |
|------|-------------|
| `--api-key` | API key (overrides env and config) |
| `--json` | Output as JSON |
| `--quiet` | Suppress progress output |
| `--verbose` | Show detailed request/response info |
| `--version` | Print version |

## Documentation

Full documentation at [renderscreenshot.com/docs/cli](https://renderscreenshot.com/docs/cli).

## License

MIT
