# MQTT Motor Backend - Device Management Complete

A Go backend server for MQTT motor control with incremental development. This project is built step-by-step, with each phase adding new features while maintaining clean, well-documented code.

## 🎯 Current Phase: Device Management ✅

### What We've Built

#### ✅ **Foundation (Phase 1)**
- **Configuration Management**: Environment variables with sensible defaults
- **Database Connection**: SQLite with GORM ORM
- **HTTP Server**: Gin framework with health endpoint
- **Project Structure**: Clean, modular architecture

#### ✅ **User Authentication (Phase 2)**
- **User Model**: Database schema with password hashing
- **Registration**: `POST /register` endpoint with validation
- **Login**: `POST /login` endpoint with JWT token generation
- **Authentication Middleware**: JWT token validation for protected routes
- **Environment Configuration**: `.env` file support with comprehensive settings

#### ✅ **Device Management (Phase 3)**
- **Device Models**: Database schema for devices and activation logs
- **Device Activation**: `POST /api/activate-device` endpoint with queue system
- **Asynchronous Processing**: Background goroutine for device control
- **Quota Management**: Daily usage limits with thread-safe implementation
- **Device State Management**: ON/OFF state tracking with database persistence

## 🗄️ Database Schema (ERD)

### Current Schema

```
┌─────────────────┐    ┌─────────────────────┐    ┌─────────────────┐
│      users      │    │      devices        │    │ deviceActivation│
├─────────────────┤    ├─────────────────────┤    ├─────────────────┤
│ id (PK)         │    │ id (PK)             │    │ id (PK)         │
│ email (UNIQUE)  │    │ name                │    │ user_id (FK)    │
│ password        │    │ state               │    │ device_id       │
│ created_at      │    │ created_at          │    │ duration        │
│ updated_at      │    │ updated_at          │    │ created_at      │
│ deleted_at      │    │ deleted_at          │    └─────────────────┘
└─────────────────┘    └─────────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                                 ▼
                        ┌─────────────────┐
                        │   deviceLogs    │
                        ├─────────────────┤
                        │ id (PK)         │
                        │ user_id (FK)    │
                        │ device_id       │
                        │ changed_at      │
                        │ state           │
                        │ duration        │
                        └─────────────────┘
```

### Schema Details

| Table | Field | Type | Constraints | Description |
|-------|-------|------|-------------|-------------|
| `users` | `id` | `uint` | `PRIMARY KEY, AUTO_INCREMENT` | Unique identifier for each user |
| `users` | `email` | `varchar(255)` | `UNIQUE, NOT NULL` | User's email address (unique) |
| `users` | `password` | `varchar(255)` | `NOT NULL` | Hashed password using bcrypt |
| `users` | `created_at` | `timestamp` | `NOT NULL` | When the user account was created |
| `users` | `updated_at` | `timestamp` | `NOT NULL` | When the user account was last updated |
| `users` | `deleted_at` | `timestamp` | `NULL` | Soft delete timestamp (NULL = active) |
| `devices` | `id` | `uint` | `PRIMARY KEY, AUTO_INCREMENT` | Unique identifier for each device |
| `devices` | `name` | `varchar(255)` | `NOT NULL` | Device name (e.g., "Motor") |
| `devices` | `state` | `enum` | `NOT NULL, DEFAULT 'UNKNOWN'` | Current state (ON/OFF/UNKNOWN) |
| `devices` | `created_at` | `timestamp` | `NOT NULL` | When the device was created |
| `devices` | `updated_at` | `timestamp` | `NOT NULL` | When the device was last updated |
| `devices` | `deleted_at` | `timestamp` | `NULL` | Soft delete timestamp (NULL = active) |
| `deviceActivation` | `id` | `uint` | `PRIMARY KEY, AUTO_INCREMENT` | Unique identifier for each activation |
| `deviceActivation` | `user_id` | `uint` | `FOREIGN KEY` | User who requested activation |
| `deviceActivation` | `device_id` | `uint` | `FOREIGN KEY` | Device that was activated |
| `deviceActivation` | `duration` | `time.Duration` | `NOT NULL` | How long device was active |
| `deviceActivation` | `created_at` | `timestamp` | `NOT NULL` | When activation was logged |
| `deviceLogs` | `id` | `uint` | `PRIMARY KEY, AUTO_INCREMENT` | Unique identifier for each log |
| `deviceLogs` | `user_id` | `uint` | `FOREIGN KEY` | User who triggered the change |
| `deviceLogs` | `device_id` | `uint` | `FOREIGN KEY` | Device that changed state |
| `deviceLogs` | `changed_at` | `timestamp` | `NOT NULL` | When the change occurred |
| `deviceLogs` | `state` | `varchar(50)` | `NOT NULL` | New state (ON/OFF) |
| `deviceLogs` | `duration` | `time.Duration` | `NULL` | How long in that state (optional) |

### Database Features

- **Soft Deletes**: Users and devices are not permanently deleted, just marked as deleted
- **Timestamps**: Automatic creation and update timestamps
- **Password Security**: Passwords are hashed using bcrypt
- **Email Uniqueness**: Prevents duplicate user accounts
- **Device State Tracking**: Real-time device state management
- **Activation Logging**: Comprehensive logging of device activations
- **GORM Integration**: Automatic schema management and migrations

### Project Structure
```
mqtt-motor-backend-v2/
├── main.go              # 🚀 Application entry point with routes
├── go.mod               # 📦 Go module dependencies
├── .env                 # ⚙️  Environment variables (configurable)
├── .env.example         # 📋 Example environment variables
├── .gitignore           # 🚫 Git ignore rules
├── config/
│   └── config.go        # ⚙️  Configuration management
├── database/
│   └── database.go      # 🗄️  Database connection and setup
├── models/
│   ├── user.go          # 👤 User model with password hashing
│   ├── device.go        # 🔧 Device model for motor control
│   ├── deviceActivation.go # 📊 Device activation logging
│   └── deviceLog.go     # 📝 Device state change logging
├── handlers/
│   ├── user.go          # 🔐 User registration and login handlers
│   └── EnqueueDeviceActivation.go # ⚡ Device activation with queue system
└── middleware/
    └── auth.go          # 🛡️  JWT authentication middleware
```

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

#### 3. Configure Environment (Optional)
```bash
# Copy the example .env file and modify as needed
cp .env.example .env

# Edit the .env file with your specific values
nano .env
```

#### 4. Run the Server
```bash
# Start the development server
go run main.go
```

You should see output like:
```
2025/08/04 11:57:58 Starting MQTT Motor Backend on port 8080
2025/08/04 11:57:58 Database connected successfully
2025/08/04 11:57:58 Running in debug mode
[GIN-debug] Listening and serving HTTP on :8080
```

## 🔐 API Endpoints

### Public Endpoints (No Authentication Required)

#### Health Check
```bash
GET /health
```
Response:
```json
{
  "status": "ok",
  "message": "MQTT Motor Backend is running"
}
```

#### User Registration
```bash
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```
Response:
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-08-04T11:57:58.418603+05:00",
    "updated_at": "2025-08-04T11:57:58.418603+05:00"
  }
}
```

#### User Login
```bash
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```
Response:
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-08-04T11:57:58.418603+05:00",
    "updated_at": "2025-08-04T11:57:58.418603+05:00"
  }
}
```

### Protected Endpoints (Authentication Required)

#### User Profile
```bash
GET /api/profile
Authorization: Bearer <JWT_TOKEN>
```
Response:
```json
{
  "message": "Protected endpoint accessed successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2025-08-04T11:57:58.418603+05:00",
    "updated_at": "2025-08-04T11:57:58.418603+05:00"
  }
}
```

#### Device Activation
```bash
POST /api/activate-device
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "device_id": 1,
  "duration": 30
}
```
Response:
```json
{
  "status": "Request added to queue"
}
```

**Notes:**
- `device_id`: Integer ID of the device to activate
- `duration`: Integer representing minutes (will be converted to `duration * time.Minute`)
- **Asynchronous**: Request is queued and processed in background
- **Quota Check**: Subject to daily usage limits (1 hour by default)
- **Queue Protection**: Returns 429 if queue is full (max 100 pending requests)
- **Database Only**: Currently updates database state (MQTT integration coming in Phase 4)

## ⚙️ Configuration

Our application uses environment variables for configuration. All variables are optional and have sensible defaults.

### Environment Variables

| Variable | Default | Description | Example |
|----------|---------|-------------|---------|
| `DB_PATH` | `data.db` | SQLite database file path | `./myapp.db` |
| `MQTT_BROKER` | `tcp://localhost:1883` | MQTT broker URL (for Phase 4) | `tcp://broker.example.com:1883` |
| `JWT_SECRET` | `supersecret` | Secret for JWT token signing | `my-super-secret-key-123` |
| `PORT` | `8080` | HTTP server port | `3000` |
| `DEBUG_MODE` | `true` | Enable debug logging | `false` |
| `DAILY_QUOTA` | `1h` | Daily motor usage limit | `2h30m` |
| `MAX_RETRIES` | `3` | Maximum retry attempts | `5` |

### Setting Environment Variables

#### Using .env File (Recommended)
```bash
# Edit the .env file
nano .env

# Set your values
DB_PATH=./myapp.db
PORT=3000
JWT_SECRET=my-secret-key
DEBUG_MODE=false
```

#### Using System Environment Variables
```bash
# macOS/Linux
export DB_PATH="./myapp.db"
export PORT="3000"
export JWT_SECRET="my-secret-key"
export DEBUG_MODE="false"

# Run the application
go run main.go
```

## 🏗️ Architecture Overview

### Current Architecture
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
│   Middleware    │  ← JWT authentication
│   (Auth)        │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   Handlers      │  ← Business logic
│   (User/Device) │
└─────────────────┘
         │
         ▼
┌─────────────────┐    ┌─────────────────┐
│   SQLite        │    │  Background     │  ← Asynchronous
│   Database      │    │   Processor     │    device control
└─────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│   Device        │    │   MQTT Broker   │  ← Real-time
│   State         │    │   (Phase 4)     │    communication
└─────────────────┘    └─────────────────┘
```

### How It Works

1. **Client Request**: A client sends an HTTP request with JWT token
2. **Gin Router**: Routes the request to appropriate handler
3. **Middleware**: JWT authentication for protected routes
4. **Handler Processing**: Business logic (device activation, etc.)
5. **Queue System**: Device requests are queued for background processing
6. **Background Processing**: Asynchronous device control with quota management
7. **Database Operations**: Device state and activation logging
8. **Response**: Immediate JSON response with queue status

### Key Technologies

- **Gin**: High-performance HTTP web framework for Go
- **GORM**: Object-Relational Mapping for database operations
- **SQLite**: Lightweight, file-based database
- **JWT**: JSON Web Tokens for authentication
- **bcrypt**: Secure password hashing
- **godotenv**: Environment variable management
- **Goroutines**: Concurrent background processing
- **Channels**: Thread-safe communication between components

## 🔄 Next Phases

### Phase 4: MQTT Integration (Coming Next)
- **MQTT Client**: Connection to MQTT broker for device communication
- **Real-time Control**: Direct MQTT commands to ESP32 devices
- **Device Communication**: Publish/subscribe for device state updates
- **Live State Updates**: Real-time device state synchronization

### Phase 5: Advanced Features
- **Device Discovery**: Automatic device registration
- **Real-time Monitoring**: Live device state updates
- **Advanced Quota**: Per-user and per-device quotas
- **Device Scheduling**: Time-based device activation

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

### Testing the API
```bash
# Test health endpoint
curl http://localhost:8080/health

# Test user registration
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'

# Test user login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}'

# Test device activation (replace TOKEN with actual JWT token)
curl -X POST http://localhost:8080/api/activate-device \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"device_id": 1, "duration": 30}'

# Test protected endpoint (replace TOKEN with actual JWT token)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer TOKEN"
```

## 📝 Code Quality

This project emphasizes:
- **Comprehensive Comments**: Every function and important line is documented
- **Clean Architecture**: Separation of concerns with clear module boundaries
- **Incremental Development**: Building features step by step
- **Error Handling**: Proper error handling throughout the application
- **Security**: Password hashing, JWT authentication, input validation
- **Configuration Management**: Environment-based configuration
- **Concurrency**: Thread-safe operations with mutexes and channels
- **Asynchronous Processing**: Non-blocking API responses with background processing

## 🔒 Security Features

- **Password Hashing**: bcrypt for secure password storage
- **JWT Authentication**: Stateless authentication with tokens
- **Input Validation**: Email format and password strength validation
- **Error Messages**: Generic error messages to prevent information leakage
- **Protected Routes**: Middleware-based route protection
- **Quota Enforcement**: Daily usage limits to prevent abuse
- **Queue Protection**: Prevents system overload with capacity limits

## 🤝 Contributing

When adding new features:
1. Follow the incremental phase approach
2. Add comprehensive comments explaining what and why
3. Update this README with new features
4. Test thoroughly before moving to next phase
5. Ensure thread-safety for concurrent operations

## 📄 License

MIT License - feel free to use this code for your own projects!
