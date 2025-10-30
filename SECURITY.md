# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

1. **DO NOT** create a public GitHub issue
2. Email security details to the maintainers
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if available)

## Security Best Practices

When using DBML:

1. **Schema Validation**: Always validate generated DBML before use in production environments.

2. **Input Sanitization**: If generating DBML from user input, sanitize strings properly to prevent injection attacks.

3. **Generated SQL**: If using DBML to generate SQL schemas, review the output for any security implications before executing in production databases.

## Security Features

DBML is designed with security in mind:

- Zero external dependencies (reduces supply chain risks)
- No network operations
- No file system operations beyond normal Go imports
- Pure Go implementation
- Thread-safe operations

## Acknowledgments

We appreciate responsible disclosure of security vulnerabilities.
