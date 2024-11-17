# 🐦 Chirpy

## Description
Chirpy is a robust HTTP server built in Go that powers a microblogging platform. Think of it like Twitter, but simpler and more focused. Users can post short messages called "chirps" and we handle all the behind-the-scenes magic to make it work smoothly.

## Architecture
Our system is built like a well-organized restaurant:

### 🏠 Front Door (HTTP Server)
- Like a restaurant's entrance, this is where all requests come in
- Handles things like "I want to post a chirp" or "Show me all chirps"
- Makes sure people are who they say they are (authentication)

### 🗄️ Kitchen (Database)
- PostgreSQL database, like a restaurant's kitchen where all the ingredients are stored
- Keeps track of:
  - Users and their information
  - All the chirps people have posted
  - Who has premium membership (Chirpy Red)
  - Login tokens (like membership cards)

### 🔒 Security System
- JWT tokens (like special VIP passes that prove who you are)
- Password hashing (turns passwords into secret codes)
- API keys for external services (like special passwords for delivery partners)

### 💳 Payment System
- Integration with Polka (our payment processor)
- Handles premium membership upgrades
- Like a cash register that talks to the bank

### 🚦 Traffic Control (Middleware)
- Checks if users are allowed to do certain things
- Counts how many people visit
- Like having security guards and people counters at the door

## Architecture

```mermaid
graph TD
    subgraph Client
        UI[HTTP Client/Browser]
    end

    subgraph HTTP_Server[HTTP Server - Port 8080]
        MUX[Router/Mux]
        Auth[Auth Middleware]
        Handlers[Request Handlers]
        Metrics[Metrics Counter]
    end

    subgraph Core_Services[Core Services]
        UserService[User Service]
        ChirpService[Chirp Service]
        AuthService[Auth Service]
        PaymentService[Payment Service]
    end

    subgraph Database[PostgreSQL Database]
        Users[(Users Table)]
        Chirps[(Chirps Table)]
        RefreshTokens[(Refresh Tokens Table)]
    end

    subgraph External[External Services]
        Polka[Polka Payment System]
    end

    %% Client interactions
    UI -->|HTTP Requests| MUX
    
    %% HTTP Server flow
    MUX --> Auth
    Auth --> Handlers
    Handlers --> Metrics
    
    %% Core Services interactions
    Handlers -->|User Operations| UserService
    Handlers -->|Chirp Operations| ChirpService
    Handlers -->|Auth Operations| AuthService
    Handlers -->|Payment Webhooks| PaymentService
    
    %% Database interactions
    UserService -->|CRUD| Users
    ChirpService -->|CRUD| Chirps
    AuthService -->|Manage Tokens| RefreshTokens
    
    %% External service interactions
    Polka -->|Webhooks| PaymentService
    PaymentService -->|Update Status| Users

    %% Styling
    classDef primary fill:#f9f,stroke:#333,stroke-width:2px
    classDef secondary fill:#bbf,stroke:#333,stroke-width:2px
    classDef database fill:#dfd,stroke:#333,stroke-width:2px
    classDef external fill:#fdd,stroke:#333,stroke-width:2px
    
    class UI,MUX,Auth,Handlers,Metrics primary
    class UserService,ChirpService,AuthService,PaymentService secondary
    class Users,Chirps,RefreshTokens database
    class Polka external
```

### Component Details

1. **Client Layer**
   - HTTP clients (browsers, apps) that send requests to our server

2. **HTTP Server Layer**
   - Router: Directs traffic to correct handlers
   - Auth: Checks user permissions
   - Handlers: Process requests
   - Metrics: Tracks usage

3. **Core Services**
   - User Service: Manages accounts
   - Chirp Service: Handles posts
   - Auth Service: Handles login/security
   - Payment Service: Manages subscriptions

4. **Database Layer**
   - Users: Stores account info
   - Chirps: Stores posts
   - Refresh Tokens: Manages sessions

5. **External Services**
   - Polka: Handles payments

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
Optional parameters:
- sort=asc|desc (default: asc) - Sort by creation time
- author_id=<uuid> - Filter by author

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