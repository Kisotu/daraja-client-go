# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability, please do NOT open a public issue.

Instead, report it via email to the project maintainers. We aim to acknowledge reports within 48 hours and provide a fix timeline.

## Supported Versions

| Version | Supported |
|---------|-----------|
| 1.x     | Yes       |

## Security Best Practices for Users

- Never hardcode credentials. Use environment variables or a secrets manager.
- Always use HTTPS for callback URLs.
- Validate callback signatures and correlation IDs in your callback handler.
- Rotate credentials regularly.