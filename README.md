# 🐦 Chirpy

## Description
Chirpy is a robust HTTP server built in Go that powers a microblogging platform. It provides a RESTful API for creating and retrieving short messages called "chirps", along with user management capabilities. The server includes features like content moderation, request metrics tracking, and PostgreSQL integration.

## Why?
Modern social platforms are often bloated with features and complex infrastructure. Chirpy aims to:
- Demonstrate clean API design with Go
- Show how to build a maintainable microservice
- Implement practical features like content moderation and metrics
- Serve as a reference for building production-ready web services
- Showcase integration with PostgreSQL using type-safe queries

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
   ```
5. Run the server:
   ```bash
   go run .
   ```

## Usage
The API provides the following endpoints:

### Users
```
POST /api/users
Create a new user with email
```

### Chirps
```
POST /api/chirps
Create a new chirp (140 char limit)

GET /api/chirps
Get all chirps

GET /api/chirps/{chirpID}
Get a specific chirp
```

### Admin
```
GET /admin/metrics
View request metrics

POST /admin/reset
Reset database (dev only)
```

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request