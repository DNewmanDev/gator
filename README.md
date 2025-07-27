# Gator - RSS Feed Aggregator #

This interactive command-line RSS feed aggregator built with Go features user management, feed subscription, and automated content fetching with PostgreSQL persistence.

## 🚀 Key Features

- **Multi-user Support**: User registration and authentication system
- **Feed Management**: Subscribe to RSS feeds with duplicate detection
- **Background Processing**: Concurrent feed fetching with configurable intervals
- **Content Persistence**: PostgreSQL database with comprehensive schema
- **CLI Interface**: Intuitive command-line interface with context-aware operations

## 🏗️ Technical Architecture

### Database Design
- **5 Migration System**: Evolutionary database schema with proper versioning using Goose
- **Complex Relationships**: Users, feeds, feed_follows, and posts with foreign key constraints
- **SQLC Integration**: Type-safe, compiled SQL queries with zero runtime reflection

### Concurrency & Performance
- **Context-Aware Processing**: Proper Go context usage for request lifecycle management
- **Goroutine-Based Fetching**: Concurrent RSS processing with rate limiting
- **HTTP Client Optimization**: Custom timeouts and connection management

### Architecture Patterns
- **Repository Pattern**: Clean data access layer abstraction
- **Command Pattern**: Extensible CLI command registration system
- **Middleware Architecture**: Authentication and validation layers

## 💻 Technologies Used

- **Go 1.24+**: Modern Go with modules and generics
- **PostgreSQL**: Relational database with complex schemas
- **SQLC**: Compile-time SQL query generation
- **RSS/XML Processing**: Custom RSS parser with HTML entity decoding
- **JSON Configuration**: User configuration persistence

## 🔧 Technical Highlights

### Advanced Go Patterns
```go
// Context-aware HTTP requests with timeouts
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Type-safe database queries via SQLC
user, err := db.GetUser(ctx, username)
```

### Database Schema Management
- Incremental migrations with rollback support
- Foreign key relationships ensuring data integrity
- Optimized queries with proper indexing

### RSS Processing Pipeline
- XML parsing with robust error handling
- HTML entity decoding for proper content display
- Duplicate post detection and filtering

## 🎯 Skills Demonstrated

- **Database Architecture**: Complex relational schema design and migration management
- **Concurrent Programming**: Goroutines, context management, and thread safety
- **HTTP Client Development**: Custom clients with proper timeout and error handling
- **CLI Framework Design**: Command registration and execution patterns
- **Error Handling**: Comprehensive error propagation and user feedback
- **Code Generation**: Integration with SQLC for type-safe database operations

## 🚀 Getting Started

### Prerequisites
- Go 1.24+
- PostgreSQL 12+

### Installation
```bash
# Clone and build
git clone https://github.com/DonJNewman/gator.git
cd gator
go build -o gator
```

### Database Setup
```bash
# Create database and run migrations
createdb gator
# Migrations are handled automatically on first run
```

### Usage
```bash
# Register a new user
./gator register <username>

# Add RSS feeds
./gator addfeed <name> <url>

# Follow feeds
./gator follow <feed_name>

# Start background aggregation (30 second intervals)
./gator agg 30s

# Browse collected posts
./gator browse [limit]
```

## 📋 Available Commands

- `register <username>` - Create new user account
- `login <username>` - Switch active user
- `reset` - Clear all data
- `users` - List all registered users
- `addfeed <name> <url>` - Add new RSS feed
- `feeds` - List all available feeds
- `follow <feed_name>` - Subscribe to a feed
- `following` - Show your feed subscriptions
- `unfollow <feed_name>` - Unsubscribe from feed
- `agg <duration>` - Start background feed fetching
- `browse [limit]` - View collected posts

This project demonstrates production-ready Go development practices, including proper database design, concurrent programming, and robust error handling suitable for enterprise applications.
