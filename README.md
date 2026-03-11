# baton-oracle-scm

`baton-oracle-scm` is a connector for Oracle Fusion Cloud Supply Chain Management (SCM) built using the [Baton SDK](https://github.com/ConductorOne/baton-sdk). It communicates with the Oracle Fusion Cloud REST APIs to sync user accounts and role assignments.

## Capabilities

| Resource Type | Sync | Provisioning |
|--------------|------|-------------|
| User         | Yes  | No          |
| Role         | Yes  | No          |

## Prerequisites

- An Oracle Fusion Cloud instance with SCM enabled
- A user account with API access and the following minimum roles:
  - `IT Security Manager` or equivalent for reading user accounts and roles
- The instance URL (e.g., `https://servername.fa.us2.oraclecloud.com`)

## Configuration

| Field          | Description                                     | Required |
|---------------|-------------------------------------------------|----------|
| `instance-url` | Oracle Fusion Cloud instance URL                | Yes      |
| `username`     | Username for API access                         | Yes      |
| `password`     | Password for API access (Basic Authentication)  | Yes      |

## Usage

```bash
baton-oracle-scm \
  --instance-url "https://your-instance.fa.us2.oraclecloud.com" \
  --username "your-api-user" \
  --password "your-api-password"
```

## API Endpoints Used

This connector uses the following Oracle Fusion Cloud REST API endpoints:

- **User Accounts**: `/hcmRestApi/resources/latest/userAccounts` - Lists all user accounts
- **Roles**: `/hcmRestApi/resources/latest/roles` - Lists all security roles
- **User Role Assignments**: `/hcmRestApi/resources/latest/userAccounts/{UserId}/child/userAccountRoles` - Lists roles assigned to a user

## Authentication

The connector uses **Basic Authentication** over HTTPS, which is the standard authentication method for Oracle Fusion Cloud REST APIs. Ensure your Oracle Fusion Cloud instance enforces SSL/TLS.
