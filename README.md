# GoPass - Secure Password Manager

GoPass is a production-ready, highly secure Password Manager built with Golang, Fiber, and Clean Architecture. It features AES-GCM encryption, Google OIDC authentication, and enterprise-grade security practices.

![Dashboard Screenshot](https://via.placeholder.com/800x400?text=GoPass+Dashboard+Preview)

## 🚀 Features

-   **Zero-Knowledge Architecture**: Secrets are encrypted using AES-GCM before storage.
-   **Authentication**: Secure Google OIDC login with Redis-backed session management.
-   **Secrets Management**: Create, Read, Update, and Delete secrets securely.
-   **Encrypted Backups**: Export and Import secrets as encrypted JSON files.
-   **Modern UI**: Server-side rendered UI using Fiber Templates and TailwindCSS.
-   **Documentation**: Interactive Swagger API documentation.

## 🛠 Tech Stack

-   **Language**: Go 1.23+
-   **Web Framework**: [Fiber v2](https://github.com/gofiber/fiber)
-   **Database**: PostgreSQL 16
-   **Cache/Session**: Redis 7
-   **Configuration**: Viper
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

Create a `.env` file or export the following variables (Docker Compose uses environment variables):

```bash
# App
SERVER_PORT=:8080
ENCRYPTION_KEY=your_32_byte_master_key_here_1234 # MUST be 32 bytes for AES-256

# Database
DB_SOURCE=postgresql://postgres:postgres@postgres:5432/password_manager?sslmode=disable

# Redis
REDIS_ADDR=redis:6379

# Google OAuth2
GOOGLE_CLIENT_ID=your_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/callback
SESSION_SECRET=super-secret-session-key
```

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

### Integration Tests
Run integration tests using Testcontainers (requires Docker):

```bash
go test -v ./tests/integration/...
```

## 🛡 Security Notes

-   **Encryption**: This MVP uses a server-side Master Key (`ENCRYPTION_KEY`). For an enterprise deployment, consider using a Key Management Service (KMS) or implementing client-side encryption.
-   **Session**: Sessions are stored in Redis with secure cookie attributes (HttpOnly).

## 📄 License

MIT License.
