# Contributing to Duet

First off, thank you for considering contributing to Duet! It's people like you that make Duet such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the issue list as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible.

### Suggesting Enhancements

If you have a suggestion for a new feature or enhancement, first check the issue list to see if it's already been proposed. If it hasn't, feel free to create a new issue detailing your suggestion.

### Pull Requests

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run the tests (`task test`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages. This helps us automatically determine version numbers and generate changelogs.

Examples:

- `feat: add new AWS provider`
- `fix: correct SSH connection handling`
- `docs: update installation instructions`
- `chore: update dependencies`

## Development Setup

1. Install Go 1.21 or later
2. Install Task (`go install github.com/go-task/task/v3/cmd/task@latest`)
3. Clone the repository
4. Install dependencies (`task install-tools`)
5. Run tests (`task test`)

## License

By contributing, you agree that your contributions will be licensed under its MIT License.
