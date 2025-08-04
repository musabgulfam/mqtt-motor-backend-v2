# MQTT Motor Backend - Phase 1: Foundation

A Go backend server for MQTT motor control with incremental development. This project is built step-by-step, with each phase adding new features while maintaining clean, well-documented code.

## 🎯 Current Phase: Foundation ✅

### What We've Built

#### ✅ **Configuration Management**
- **Purpose**: Centralizes all application settings in one place
- **How it works**: Reads environment variables with sensible defaults
- **Benefits**: Easy to configure for different environments (dev, staging, production)
- **Files**: `config/config.go`

#### ✅ **Database Connection**
- **Purpose**: Provides persistent storage for users, device data, and logs
- **Technology**: SQLite with GORM ORM for easy database operations
- **Benefits**: Lightweight, file-based, no separate server needed
- **Files**: `database/database.go`

#### ✅ **HTTP Server**
- **Purpose**: Provides REST API endpoints for client applications
- **Framework**: Gin (high-performance Go web framework)
- **Features**: Built-in logging, error recovery, JSON handling
- **Files**: `main.go`

#### ✅ **Health Check Endpoint**
- **Purpose**: Allows monitoring systems to verify the server is running
- **Endpoint**: `GET /health`
- **Response**: JSON with status and message
- **Use cases**: Load balancers, monitoring tools, client health checks

### Project Structure Explained

```
mqtt-motor-backend-v2/
├── main.go              # 🚀 Application entry point - orchestrates everything
├── go.mod               # 📦 Go module dependencies and versions
├── config/
│   └── config.go        # ⚙️  Configuration management - environment variables
└── database/
    └── database.go      # 🗄️  Database connection and setup
```

#### **File Descriptions:**

- **`main.go`**: The heart of our application. It:
  - Loads configuration from environment variables
  - Establishes database connection
  - Sets up HTTP server with Gin framework
  - Defines API endpoints
  - Starts listening for requests

- **`config/config.go`**: Manages all application settings:
  - Database path (where SQLite file is stored)
  - MQTT broker URL (for device communication)
  - JWT secret (for user authentication)
  - Server port (HTTP listening port)
  - Daily quotas (usage limits)

- **`database/database.go`**: Handles database operations:
  - Connects to SQLite database
  - Sets up GORM ORM for easy data operations
  - Provides global database access
  - Will handle schema migrations in future phases

## 🚀 Installation & Setup

### Prerequisites
- **Go** (1.18 or newer) - Download from [golang.org](https://golang.org/dl/)

### Step-by-Step Setup

#### 1. Clone and Navigate
```bash
# Navigate to your project directory
cd mqtt-motor-backend-v2
```

#### 2. Install Dependencies
```bash
# Download and install all required Go packages
go mod tidy
```

#### 3. Run the Server
```bash
# Start the development server
go run main.go
```

You should see output like:
```
2025/08/04 11:30:33 Starting MQTT Motor Backend on port 8080
2025/08/04 11:30:33 Database connected successfully
[GIN-debug] Listening and serving HTTP on :8080
```

#### 4. Test the Health Endpoint
```bash
# Test that the server is running
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "MQTT Motor Backend is running"
}
```

## ⚙️ Configuration

Our application uses environment variables for configuration. All variables are optional and have sensible defaults.

### Environment Variables

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `DB_PATH` | `data.db` | SQLite database file path | `./myapp.db` |
| `MQTT_BROKER` | `tcp://localhost:1883` | MQTT broker URL | `tcp://broker.example.com:1883` |
| `JWT_SECRET` | `supersecret` | Secret for JWT token signing | `my-super-secret-key-123` |
| `PORT` | `8080` | HTTP server port | `3000` |
| `DAILY_QUOTA` | `1h` | Daily motor usage limit | `2h30m` |

### Setting Environment Variables

#### macOS/Linux
```bash
# Set variables for current session
export DB_PATH="./myapp.db"
export PORT="3000"
export JWT_SECRET="my-secret-key"

# Run the application
go run main.go
```

#### Windows
```cmd
# Set variables for current session
set DB_PATH=./myapp.db
set PORT=3000
set JWT_SECRET=my-secret-key

# Run the application
go run main.go
```

## 🏗️ Architecture Overview

### Current Architecture (Phase 1)
```
┌─────────────────┐
│   Client App    │  ← HTTP requests (REST API)
│  (Web/Mobile)   │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│  Gin HTTP       │  ← Web server with routing
│    Server       │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   SQLite        │  ← File-based database
│   Database      │
└─────────────────┘
```

### How It Works

1. **Client Request**: A client (web app, mobile app, etc.) sends an HTTP request
2. **Gin Router**: Our Gin server receives the request and routes it to the appropriate handler
3. **Handler Processing**: The handler processes the request (currently just health check)
4. **Database Operations**: When needed, handlers interact with the SQLite database
5. **Response**: The server sends back a JSON response to the client

### Key Technologies

- **Gin**: High-performance HTTP web framework for Go
- **GORM**: Object-Relational Mapping for database operations
- **SQLite**: Lightweight, file-based database
- **Go Modules**: Dependency management

## 🔄 Next Phases

### Phase 2: User Management (Coming Next)
- **User Model**: Database schema for storing user information
- **Registration**: Endpoint for creating new user accounts
- **Login**: Authentication with JWT tokens
- **Middleware**: Authentication middleware for protected endpoints

### Phase 3: MQTT Integration
- **MQTT Client**: Connection to MQTT broker for device communication
- **Motor Control**: Endpoints for controlling the motor
- **Device Communication**: Real-time communication with ESP32 devices

### Phase 4: Advanced Features
- **Motor Queue**: Queue system for motor activation requests
- **Daily Quota**: Usage limits and quota enforcement
- **Device Logging**: Comprehensive logging of device activations

## 🧪 Development & Testing

### Running Tests
```bash
# Run all tests (when we add them in future phases)
go test ./...
```

### Code Formatting
```bash
# Format all Go code
go fmt ./...
```

### Code Linting
```bash
# Check for common issues
go vet ./...
```

## 📝 Code Quality

This project emphasizes:
- **Comprehensive Comments**: Every function and important line is documented
- **Clean Architecture**: Separation of concerns with clear module boundaries
- **Incremental Development**: Building features step by step
- **Error Handling**: Proper error handling throughout the application
- **Configuration Management**: Environment-based configuration

## 🤝 Contributing

When adding new features:
1. Follow the incremental phase approach
2. Add comprehensive comments explaining what and why
3. Update this README with new features
4. Test thoroughly before moving to next phase

## 📄 License

MIT License - feel free to use this code for your own projects!
