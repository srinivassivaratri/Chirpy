# 🐦 Chirpy

## Description
Chirpy is a robust HTTP server built in Go that powers a microblogging platform. Think of it like Twitter, but simpler and more focused. Users can post short messages called "chirps" and we handle all the behind-the-scenes magic to make it work smoothly.

## Why?
Modern social platforms are often bloated with features and complex infrastructure. Chirpy aims to:
- Keep things simple and clean
- Show how to build a secure and reliable service
- Make user data safe with proper login systems
- Store data efficiently using PostgreSQL
- Handle passwords securely (because nobody wants their account hacked!)

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
   POLKA_KEY="your-polka-api-key"
   ```
5. Run migrations:
   ```bash
   goose -dir sql/schema postgres "${DB_URL}" up
   ```
6. Run the server:
   ```bash
   go run .
   ```

## Features

### 🌟 Chirpy Red
Our premium membership that gives users extra cool features:
- Edit chirps after posting
- More features coming soon!
- Automatic activation through Polka payment system

### 👤 User Management
```http
POST /api/users
Create a new account

POST /api/login
Sign in and get your access passes

PUT /api/users
Update your profile (need to be logged in)

POST /api/refresh
Get a new access pass using your refresh token

POST /api/revoke
Log out (invalidate your refresh token)
```

### 📝 Chirps
```http
POST /api/chirps
Post a new chirp (140 char limit)

GET /api/chirps
See all chirps
Optional: Filter by author with ?author_id=<uuid>

GET /api/chirps/{chirpID}
Look at a specific chirp

DELETE /api/chirps/{chirpID}
Delete your chirp (you can only delete your own!)
```

### 💳 Polka Integration
```http
POST /api/polka/webhooks
Handles automatic Chirpy Red membership activation
Requires Polka API key for security
```

### 🔧 Admin Tools
```http
GET /admin/metrics
Check how many people are using Chirpy

POST /admin/reset
Reset everything (only works in development)
```

### 🏥 Health Check
```http
GET /api/healthz
Make sure the server is alive and kicking
```

## Security Features
- Super secure password storage
- Login tokens that expire (so hackers can't use old ones)
- Refresh tokens for staying logged in safely
- Content filtering (keeps things family-friendly)
- Email uniqueness (no duplicate accounts)
- Environment-based security
- Polka webhook authentication
- Only authors can delete their chirps

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/cool-new-thing`)
3. Make your changes
4. Test everything (`go test ./...`)
5. Commit your changes (`git commit -m 'feat: add cool new thing'`)
6. Push to your branch (`git push origin feature/cool-new-thing`)
7. Open a Pull Request

## Dependencies
- github.com/golang-jwt/jwt/v5
- github.com/google/uuid
- github.com/joho/godotenv
- github.com/lib/pq
- golang.org/x/crypto