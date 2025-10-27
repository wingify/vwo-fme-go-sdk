# Git Hooks Setup

This project includes automated git hooks to ensure code quality and consistency.

## Overview

The git hooks system includes:

- **pre-push**: Runs comprehensive validations before pushing code
- **commit-msg**: Validates commit message format
- **Copyright Management**: Scripts to add and verify copyright headers

## Quick Setup

Run the setup script to install all git hooks:

```bash
./git-hooks/scripts/setup-git-hooks.sh
```

This will:
- Install pre-push and commit-msg hooks
- Set proper permissions
- Test all dependencies
- Verify scripts are working

## Git Hooks Details

### Pre-Push Hook

Runs the following validations before allowing a push:

1. **Copyright Header Check**: Ensures all Go files have proper copyright headers
2. **Go Format Check**: Runs `go fmt ./...` to format code
3. **Go Vet Check**: Runs `go vet ./...` to check for potential issues
4. **Go Build Check**: Runs `go build ./...` to ensure code compiles
5. **Go Test Suite**: Runs `go test ./test/... -v` to execute all tests

### Commit Message Hook

Validates commit messages follow the conventional commit format:

```
<type>(<scope>): <subject>
```

**Valid types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `cleanup`: Code cleanup
- `revert`: Reverting changes
- `tracking`: Tracking changes

**Examples:**
- ✅ `feat: add new API endpoint`
- ✅ `fix(api): resolve authentication issue`
- ✅ `docs: update README with setup instructions`
- ❌ `invalid commit message`
- ❌ `WIP: work in progress` (ignored)

## Copyright Management

### Adding Copyright Headers

To add copyright headers to all Go files:

```bash
# Using Node.js script
node git-hooks/scripts/add-copyright.js

# Using shell script (if available)
./git-hooks/scripts/add-copyright.sh

# Add to specific directory
node git-hooks/scripts/add-copyright.js pkg/
```

### Checking Copyright Headers

To verify all files have proper copyright headers:

```bash
node git-hooks/scripts/check-license.js
```

## Manual Hook Management

### Disable Hooks Temporarily

```bash
git config core.hooksPath /dev/null
```

### Re-enable Hooks

```bash
git config --unset core.hooksPath
```

### Run Hooks Manually

```bash
# Test pre-push hook
node .git/hooks/pre-push

# Test commit-msg hook
echo "feat: test commit" | node .git/hooks/commit-msg
```

## Requirements

- **Node.js**: Required for running the git hooks
- **Go**: Required for Go-specific validations
- **Git**: Required for git hook functionality

## Troubleshooting

### Common Issues

1. **"Cannot find module" errors**: Make sure you're running from the project root
2. **Permission denied**: Run `chmod +x` on script files
3. **Go vet errors**: Fix formatting issues in Go code
4. **Commit message rejected**: Use proper conventional commit format

### Debug Mode

To see detailed output from hooks, you can modify the scripts to include more verbose logging.

## File Structure

```
git-hooks/
├── pre-push                 # Pre-push validation hook
├── commit-msg               # Commit message validation hook
├── scripts/
│   ├── add-copyright.js     # Add copyright headers script
│   ├── check-license.js     # Verify copyright headers script
│   ├── setup-git-hooks.sh   # Installation script
│   ├── run-tests.sh         # Test runner script
│   ├── start.sh             # Quick start script
│   ├── utils/
│   │   ├── run.js          # Task runner utility
│   │   └── CheckLicenseUtil.js  # License checking utility
│   └── enums/
│       └── AnsiColorEnum.js # Color output utilities
```

## Test Management

### Running Tests

To run the complete test suite:

```bash
# Using the test runner script
./git-hooks/scripts/run-tests.sh

# Using Go directly
go test ./test/... -v

# Run specific test categories
go test ./test/e2e/... -v    # E2E tests only
go test ./test/unit/... -v    # Unit tests only

# Run with additional options
go test ./test/... -cover     # With coverage
go test ./test/... -race      # With race detection
go test ./test/... -short     # Skip long-running tests
```

### Test Categories

The test suite includes:

- **E2E Tests** (`test/e2e/`): End-to-end feature flag functionality tests
- **Unit Tests** (`test/unit/`): Individual component tests including:
  - Segmentation evaluators
  - Operators (AND, OR, NOT, etc.)
  - Data type validations
  - Decision maker logic
  - Client validation

### Test Data

Test data is organized in `test/data/`:
- **Settings**: Various campaign configurations
- **Test Cases**: Structured test scenarios
- **Storage**: Storage-related test utilities

## Contributing

When contributing to this project:

1. Ensure your commit messages follow the conventional commit format
2. All Go files should have proper copyright headers
3. Code should pass all pre-push validations including tests
4. Test your changes with the git hooks before pushing
5. Run the test suite to ensure no regressions: `./git-hooks/scripts/run-tests.sh`
