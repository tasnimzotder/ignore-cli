# Ignore CLI

<div align="center">

![GitHub Release](https://img.shields.io/github/v/release/tasnimzotder/ignore-cli?style=flat-square)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/tasnimzotder/ignore-cli/build-test-release.yml?branch=main&style=flat-square)
![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/tasnimzotder/ignore-cli/total?style=flat-square)

</div>

A CLI tool for managing `.gitignore` files using GitHub's gitignore templates.

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
# Add a template to .gitignore
ignore add Go
ignore add Python --override  # Replace existing .gitignore

# List available templates
ignore list

# Search templates
ignore search node
```

## License

MIT License - see [LICENSE](LICENSE) for details.
