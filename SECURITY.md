# Security Policy

## Supported Versions

The `main` branch and latest semver tag receive security fixes.

## Reporting a Vulnerability

Please report vulnerabilities privately to florinbadita@users.noreply.github.com.

Do not open public issues with private videos, location details, license plates, API keys, or other sensitive evidence.

## Secrets

The frontend never requires secrets. Runtime backend secrets belong in `deploy/.env` or server environment variables and must not be committed.
