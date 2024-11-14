# 🐦 Chirpy

## Description
Chirpy is a robust HTTP server built in Go that powers a microblogging platform. It provides a RESTful API for creating and retrieving short messages called "chirps", along with secure user authentication.

## Why?
Modern social platforms are often bloated with features and complex infrastructure. Chirpy aims to:
- Demonstrate clean API design with Go
- Show how to build a maintainable microservice
- Implement secure user authentication with JWT
- Showcase PostgreSQL integration with type-safe queries
- Provide practical examples of password hashing and validation

## Quick Start
1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up PostgreSQL and create a database named `chirpy`
4. Create a `.env` file:
   ```
   DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"
   JWT_SECRET="your-secret-key"
   ```
5. Run migrations:
   ```bash
   goose -dir sql/schema postgres "${DB_URL}" up
   ```
6. Run the server:
   ```bash
   go run .
   ```

## Usage

### Authentication
```http
POST /api/users
Create a new user with email and password
{
    "email": "user@example.com",
    "password": "securepassword"
}

POST /api/login
Login with email and password (supports configurable token expiration)
{
    "email": "user@example.com",
    "password": "securepassword",
    "expires_in_seconds": 3600  // optional
}
```

### Chirps
```http
POST /api/chirps
Create a new chirp (requires JWT, 140 char limit)
{
    "body": "Hello world!"
}

GET /api/chirps
Get all chirps

GET /api/chirps/{chirpID}
Get a specific chirp
```

### Admin
```http
GET /admin/metrics
View request metrics

POST /admin/reset
Reset database (dev only)
```

### Health Check
```http
GET /api/healthz
Check API health
```

## Security Features
- Password hashing using bcrypt
- JWT-based authentication
- Configurable token expiration
- Content moderation (bad word filtering)
- Database-level email uniqueness
- Environment-based security controls

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Dependencies
- github.com/golang-jwt/jwt/v5
- github.com/google/uuid
- github.com/joho/godotenv
- github.com/lib/pq
- golang.org/x/crypto