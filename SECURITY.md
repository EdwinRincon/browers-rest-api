# Security Policy

## Supported Versions

We actively maintain and provide security updates for the following versions of Browers REST API:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| older   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue in this project, please follow responsible disclosure practices.

### How to Report

**Please do NOT report security vulnerabilities via public GitHub issues.**

Instead, report vulnerabilities through one of the following channels:

1. **GitHub Private Security Advisory**: Use the [Security tab](https://github.com/EdwinRincon/browers-rest-api/security/advisories/new) to create a private advisory.
2. **Direct contact**: Reach out to the maintainer directly via GitHub profile.

### What to Include

When reporting a vulnerability, please include:

- A clear description of the vulnerability
- Steps to reproduce the issue
- Potential impact and severity assessment
- Any suggested fixes or mitigations (optional)
- Your contact information for follow-up

## Response Timeline

- **Acknowledgement**: Within 48 hours of receiving your report
- **Initial Assessment**: Within 5 business days
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days

## Security Measures

This repository employs the following security practices:

- **Dependabot**: Automated dependency vulnerability scanning and updates
- **gosec**: Static security analysis for Go code (runs on every PR)
- **golangci-lint**: Code quality and security linting
- **Branch Protection**: Required checks on the `main` branch
- **CODEOWNERS**: All changes require review from designated code owners

## Scope

### In Scope
- Authentication and authorization vulnerabilities
- SQL injection and database security
- API endpoint security issues
- Dependency vulnerabilities
- Container image security
- Secrets or credentials accidentally committed

### Out of Scope
- Issues in third-party dependencies (report directly to the vendor)
- Social engineering attacks
- Physical security

## Disclosure Policy

We follow a **coordinated disclosure** policy. We ask that:

1. You give us reasonable time to investigate and patch the vulnerability before public disclosure
2. You do not exploit the vulnerability beyond what is necessary to demonstrate it
3. You do not access or modify user data without permission

We will credit security researchers in our release notes (unless you prefer anonymity).

---

*This security policy is based on best practices for open-source projects.*
