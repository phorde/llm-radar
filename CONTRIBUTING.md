# Contributing to OpenCode Check

Thank you for your interest in contributing to OpenCode Check! This document provides guidelines and instructions for contributing.

## 🌟 Ways to Contribute

- 🐛 **Report bugs** - Open an issue with details about the problem
- 💡 **Suggest features** - Share ideas for improvements
- 📝 **Improve documentation** - Fix typos, add examples, clarify instructions
- 🔧 **Submit code** - Fix bugs or implement new features
- 🧪 **Test** - Try the tool with different configurations and report results

## 🚀 Getting Started

### Prerequisites
16: 
17: - Go 1.24 or higher
18: - OpenCode CLI installed and configured
19: - Git
20: - A GitHub account
21: 

### Setting Up Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/opencode-check.git
   cd opencode-check
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/original-owner/opencode-check.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Build and test**:
   ```bash
   go build -o opencode-check
   ./opencode-check --help
   ```

### 🤖 AI Agent Guidelines
51: 
52: All AI agents working in this environment **MUST** follow the specific instructions in [docs/AGENT_GUIDELINES.md](docs/AGENT_GUIDELINES.md).
53: 

## 🔨 Development Workflow

### Making Changes

1. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/bug-description
   ```

2. **Make your changes** following our coding standards (see below)

3. **Test your changes**:
   ```bash
   # Build
   go build -o opencode-check
   
   # Run basic tests
   ./opencode-check --version
   ./opencode-check --help
   
   # Test with different configurations
   ./opencode-check -c 3
   ./opencode-check --cache
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   # or
   git commit -m "fix: resolve issue with..."
   ```

   Use conventional commits:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes (formatting, etc.)
   - `refactor:` - Code refactoring
   - `test:` - Adding or updating tests
   - `chore:` - Maintenance tasks

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request** on GitHub

### Pull Request Guidelines

- **Title**: Use clear, descriptive titles
- **Description**: Explain what changes you made and why
- **Link issues**: Reference related issues with `Fixes #123` or `Relates to #456`
- **Keep it focused**: One PR should address one feature or fix
- **Update documentation**: If you change behavior, update README.md

## 📋 Coding Standards

### Go Code Style

- Follow standard Go formatting: use `gofmt`
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Keep functions focused and small
- Add comments for exported functions and complex logic
- Use meaningful variable and function names

### Code Organization

```
opencode-check/
├── main.go              # Main application logic
├── docs/                # Documentation
├── config/              # Configuration files
├── models/              # Data models (if needed)
└── kb-custom.json       # Knowledge base example
```

### Error Handling

- Always check and handle errors
- Provide meaningful error messages
- Use `fmt.Errorf` for error wrapping

Example:
```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

### Comments

- Add doc comments for exported types and functions:
  ```go
  // ModelResult represents the test result for a single model
  type ModelResult struct {
      // ...
  }
  ```

## 🧪 Testing

Before submitting a PR, ensure:

1. **Code compiles** without errors:
   ```bash
   go build
   ```

2. **No formatting issues**:
   ```bash
   gofmt -w .
   ```

3. **Manual testing** with real OpenCode CLI:
   ```bash
   ./opencode-check -c 2
   ```

## 📝 Documentation

When adding features:

1. Update `README.md` with new flags or behavior
2. Update `README.pt-BR.md` with Portuguese translation
3. Add examples if applicable
4. Update relevant documentation in `docs/`

## 🐛 Reporting Bugs

When reporting bugs, include:

- **Description**: Clear description of the problem
- **Steps to reproduce**: Numbered list of steps
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Environment**:
  - OS (Linux/macOS/Windows)
  - Go version (`go version`)
  - OpenCode version (`opencode --version`)
  - Terminal emulator
- **Logs/Screenshots**: Any relevant output or screenshots

## 💡 Feature Requests

For feature requests:

- Check if it's already been requested in Issues
- Describe the use case and problem it solves
- Explain how you envision it working
- Consider if it fits the project's scope

## 🎯 Areas for Contribution

Some areas where we especially welcome contributions:

- 🧪 **Testing**: Automated tests for core functionality
- 📊 **Knowledge Base**: Updates to model classifications
- 🌍 **Internationalization**: Translations to other languages
- 🎨 **UI/UX**: Improvements to TUI interface
- 📖 **Documentation**: Examples, tutorials, use cases
- 🔧 **Performance**: Optimizations for speed or resource usage

## ❓ Questions

If you have questions:

- Check existing [Issues](https://github.com/your-username/opencode-check/issues)
- Check the [OpenCode Documentation](https://opencode.ai/docs)
- Open a new issue with the `question` label

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Thank You!

We appreciate your time and effort in contributing to OpenCode Check. Every contribution, no matter how small, helps make this tool better for the community!

---

**Note**: This is a community project not affiliated with Anomaly (OpenCode creators). For OpenCode CLI issues, please refer to the [official repository](https://github.com/anomalyco/opencode).
