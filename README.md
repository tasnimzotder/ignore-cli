# Ignore

`Ignore` is a command-line tool for managing `.gitignore` files and performing various operations related to ignoring files in a Git repository.

## Table of Contents

- [Ignore](#ignore)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [From Source](#from-source)
  - [Usage](#usage)
    - [Commands](#commands)
  - [Contributing](#contributing)
  - [License](#license)

## Installation

To install the `ignore` binary, you can download it from the [releases](https://github.com/tasnimzotder/ignore-cli/releases) page or build it from source.

### From Source

To build the binary from source, ensure you have Go installed and run the following commands:

```sh
git clone https://github.com/your-repo/ignore.git
cd ignore
make build
```

The binary will be available in the build/ directory.

## Usage

You can use the `ignore` tool to manage your `.gitignore` files and perform various operations. Below are some example commands:

```sh
./build/ignore add
./build/ignore list
./build/ignore search
```

### Commands

- `add`: Add a new entry to the `.gitignore` file.
- `list`: List all entries in the `.gitignore` file.
- `search`: Search for a specific entry in the `.gitignore` file.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License
This project is licensed under the MIT License. See the LICENSE file for details.