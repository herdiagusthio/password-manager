# GoPass - Secure Password Manager

GoPass is a production-ready, highly secure Password Manager built with Golang, Fiber, and Clean Architecture. It features AES-GCM encryption, Google OIDC authentication, and enterprise-grade security practices.

## 🚀 Features

-   **Zero-Knowledge Architecture**: Secrets are encrypted using AES-GCM before storage.
-   **Authentication**: Secure Google OIDC login with Redis-backed session management.
-   **Secrets Management**: Create, Read, Update, and Delete secrets securely.
-   **Encrypted Backups**: Export and Import secrets as encrypted JSON files.
-   **Modern UI**: Server-side rendered UI (Fiber Templates + TailwindCSS) with:
    -   Secure Login Page
    -   Dashboard with Copy-to-Clipboard & Reveal functionality
    -   Add/Edit/Delete Modals
-   **Documentation**: Interactive Swagger API documentation.

## 🛠 Tech Stack

-   **Language**: Go 1.23+
-   **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
-   **Database**: PostgreSQL 16
-   **Cache/Session**: Redis 7
-   **Configuration**: Viper
-   **Testing**: Testcontainers (Integration), Testify & MockGen (Unit)
-   **CI/CD**: GitHub Actions
-   **Containerization**: Docker & Docker Compose

## 📦 Installation & Setup

### Prerequisites

-   Docker & Docker Compose
-   Go 1.23+ (optional, for local dev)
-   Google Cloud Console Project (for OAuth2)

### 1. Clone the Repository

```bash
git clone https://github.com/herdiagusthio/password-manager.git
cd password-manager
```

### 2. Configure Environment

Copy the example configuration file and fill in your credentials:

```bash
cp .env.example .env
```

**Required Variables**:
-   `GOOGLE_CLIENT_ID` & `GOOGLE_CLIENT_SECRET`: From Google Cloud Console.
-   `ENCRYPTION_KEY`: A **32-byte** hex string for AES-256 encryption.
-   `SESSION_SECRET`: Random string for signing session cookies.

### 3. Run with Docker Compose

```bash
docker-compose up -d --build
```

### 4. Run Migrations

Connect to the PostgreSQL container and run the migrations:

```bash
# Using a tool like Migrate or manually executing SQL in psql
cat migrations/000001_init.up.sql | docker-compose exec -T postgres psql -U postgres -d password_manager
```

## 📖 Usage

### User Interface
-   **Login**: Visit [http://localhost:8080/login](http://localhost:8080/login) and sign in with Google.
-   **Dashboard**: Manage your secrets at [http://localhost:8080/dashboard](http://localhost:8080/dashboard).

### API Documentation (Swagger)
Interactive API documentation is available at:
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 🧪 Testing

### Unit Tests
Run table-driven unit tests for the core business logic (UseCases):
```bash
go test -v ./internal/usecase/...
```

### Integration Tests
Run integration tests using Testcontainers (requires Docker):
```bash
go test -v ./tests/integration/...
```

## ⚙️ CI/CD

This project uses **GitHub Actions** for Continuous Integration. On every push to `main`, the workflow:
1.  Builds the application.
2.  Runs Unit Tests.
3.  Runs Integration Tests (spinning up ephemeral Postgres/Redis containers).

## 🛡 Security Notes

-   **Encryption**: This MVP uses a server-side Master Key (`ENCRYPTION_KEY`). For an enterprise deployment, consider using a Key Management Service (KMS) or implementing client-side encryption.
-   **Session**: Sessions are stored in Redis with secure cookie attributes (HttpOnly).

## 📄 License

MIT License.
