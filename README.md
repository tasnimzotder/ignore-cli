# Ignore CLI

<div align="center">

![GitHub Release](https://img.shields.io/github/v/release/tasnimzotder/ignore-cli?style=flat-square)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/tasnimzotder/ignore-cli/ci.yml?branch=main&style=flat-square)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/tasnimzotder/ignore-cli/total?style=flat-square)

</div>

A TUI tool for managing `.gitignore` files using GitHub's gitignore templates. Features an interactive template picker with fuzzy search and live preview.

## Installation

```sh
# Homebrew (macOS/Linux)
brew install tasnimzotder/tap/ignore-cli

# From source
go install github.com/tasnimzotder/ignore-cli@latest
```

Or download binaries from [releases](https://github.com/tasnimzotder/ignore-cli/releases).

## Usage

```sh
# Interactive mode — launch TUI picker
ignore add

# Add templates directly
ignore add Go
ignore add Go Python Node

# Replace existing .gitignore instead of appending
ignore add Go --override

# Browse available templates
ignore list
```

### TUI Controls

- `/` — fuzzy search
- `space` — toggle selection
- `enter` — confirm and add selected templates
- `q` / `esc` — cancel

## License

MIT License - see [LICENSE](LICENSE) for details.
