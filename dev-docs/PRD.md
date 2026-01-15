Product Requirements Document (PRD)
1. Executive Summary
Project Name: GoGuard (Placeholder Name) Objective: Develop a self-hosted, web-based password management solution that emphasizes data sovereignty. The system allows users to generate, store, and manage credentials securely, with a disaster recovery mechanism (backup/restore) strictly bound to the application instance to prevent external data leaks.

2. User Personas
The Self-Hoster: A tech-savvy user who prefers hosting their own tools via Docker to avoid vendor lock-in and ensure data privacy.

The Privacy Advocate: Prioritizes strong encryption and wants assurance that backups cannot be opened if stolen by third parties without the application context.

3. Functional Requirements
3.1 Authentication & Session
REQ-AUTH-01: Users must log in using Google OAuth2 (OIDC).

REQ-AUTH-02: No local email/password registration shall be supported (to reduce attack surface).

REQ-AUTH-03: Session management must use HTTP-Only, Secure Cookies backed by Redis.

3.2 Password Generator
REQ-GEN-01: Users can generate passwords with customizable length (8-64 chars).

REQ-GEN-02: Options to include/exclude: Uppercase, Lowercase, Numbers, Symbols.

REQ-GEN-03: The default setting must be "Strong" (16 chars, mixed types).

3.3 Credential Management (The Vault)
REQ-VAULT-01 (Create): Save a new entry with fields: Title, Username, Password, URL, Notes.

REQ-VAULT-02 (Read): View list of saved secrets (passwords masked by default).

REQ-VAULT-03 (Update): Edit existing credentials.

REQ-VAULT-04 (Delete): Remove credentials permanently.

REQ-VAULT-05 (Copy): "Click to Copy" functionality for passwords.

3.4 Security & Backups
REQ-BACK-01 (Export): Authenticated users can download a full backup of their vault.

REQ-BACK-02 (Encryption): The backup file must be encrypted using a server-side Master Key. It must be indecipherable if opened outside this application instance.

REQ-BACK-03 (Restore): Users can upload a backup file. The system will decrypt it (using the internal Master Key) and restore the entries.

Conflict Resolution: Strategy shall be "Upsert" (Update if exists, Insert if new) based on Entry Title/ID.

4. Non-Functional Requirements
Security: All secrets in the database should ideally be encrypted at rest (optional scope, but recommended). The Backup must be encrypted.

Performance: API response time < 200ms for vault retrieval.

Deployment: Fully containerized (Docker) with no external host dependencies.

Code Quality: 80%+ Test Coverage via Table-Driven Tests