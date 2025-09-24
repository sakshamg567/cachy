# Cachy - Distributed Cache System

A distributed caching system built with Go and gRPC, featuring consistent hashing for scalable data distribution and LRU eviction policies.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the System](#running-the-system)
- [API Usage](#api-usage)
- [Components Breakdown](#components-breakdown)
- [Configuration](#configuration)
- [Development](#development)
- [Contributing](#contributing)

## Overview

Cachy is a distributed caching system that automatically distributes data across multiple cache nodes using consistent hashing. The system provides high availability, horizontal scalability, and efficient data retrieval with LRU (Least Recently Used) eviction policies.

## Features

- **Distributed Architecture**: Multiple cache nodes with automatic load balancing
- **Consistent Hashing**: Efficient data distribution and minimal data movement during scaling
- **LRU Eviction**: Intelligent cache management with Least Recently Used eviction
- **Dynamic Node Addition**: Add new cache nodes without service interruption
- **gRPC Communication**: High-performance inter-service communication
- **RESTful API**: Easy-to-use HTTP endpoints for cache operations
- **Data Migration**: Automatic data redistribution when nodes are added
- **Concurrent Safe**: Thread-safe operations across all components

## Architecture

```
┌─────────────────┐    HTTP      ┌─────────────────┐
│   Client Apps   │ ────────────▶│  Server (8080)  │
└─────────────────┘              │  (Coordinator)  │
                                 └─────────────────┘
                                          │ gRPC
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ▼                     ▼                     ▼
        ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
        │ Cache Node :50051│  │ Cache Node :50052│  │ Cache Node :50053│
        │   (LRU Cache)    │  │   (LRU Cache)    │  │   (LRU Cache)    │
        └──────────────────┘  └──────────────────┘  └──────────────────┘
```

## Project Structure

```
cachy/
├── cmd/
│   ├── cache-node/          # Cache node executable
│   │   └── main.go         # Cache node entry point
│   └── server/             # Coordinator server executable
│       └── main.go         # Server entry point
├── internal/
│   ├── cache/              # Cache implementation
│   │   ├── lru.go         # LRU cache with doubly-linked list
│   │   └── node.go        # gRPC cache node implementation
│   └── coordinator/        # Coordination logic
│       ├── coordinator.go  # Main coordinator logic
│       └── hashRing.go    # Consistent hashing implementation
├── shared/
│   └── proto/
│       ├── cacheNodepb/   # Generated protobuf code
│       └── cache-node.proto # Protocol buffer definitions
├── util/
│   └── hash.go           # SHA256-based consistent hashing utility
├── Makefile             # Build and deployment automation
├── go.mod              # Go module dependencies
└── .gitignore         # Git ignore rules
```

## Prerequisites

- **Go**: Version 1.24.1 or higher
- **Protocol Buffers**: For gRPC code generation
- **Make**: For build automation (optional)
- **lsof**: For port management in Makefile

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sakshamg567/cachy.git
   cd cachy
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Build all binaries:**
   ```bash
   make build-all
   ```
   Or manually:
   ```bash
   go build -o bin/cache-node ./cmd/cache-node
   go build -o bin/server ./cmd/server
   ```

## Running the System

### Option 1: Using Makefile (Recommended)

**Start all services:**
```bash
make run-all
```

**View logs:**
```bash
make tail-logs
```

**Stop all services:**
```bash
make stop-all
```

**Clean up:**
```bash
make clean
```

### Option 2: Manual Start

**Start cache nodes:**
```bash
# Terminal 1
./bin/cache-node --port 50051

# Terminal 2  
./bin/cache-node --port 50052

# Terminal 3
./bin/cache-node --port 50053
```

**Start coordinator server:**
```bash
# Terminal 4
./bin/server --port 8080
```

## API Usage

### Set a Value
```bash
curl -X POST http://localhost:8080/set \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "value": "john_doe"}'
```

### Get a Value
```bash
curl http://localhost:8080/get?key=user:123
```

### Add a New Cache Node
```bash
curl -X POST http://localhost:8080/add-node \
  -H "Content-Type: application/json" \
  -d '{"address": "localhost:50054"}'

# Don't forget to start the new node:
./bin/cache-node --port 50054
```

## Components Breakdown

### 1. **Cache Node** (`internal/cache/`)
- **LRU Cache** (`lru.go`): Thread-safe LRU implementation using doubly-linked list
- **gRPC Server** (`node.go`): Handles cache operations (Get, Set, Delete, GetAllKeys)
- **Capacity Management**: Configurable cache size with automatic eviction

### 2. **Coordinator** (`internal/coordinator/`)
- **Request Routing** (`coordinator.go`): Routes cache requests to appropriate nodes
- **Consistent Hashing** (`hashRing.go`): SHA256-based hash ring for data distribution
- **Dynamic Scaling**: Handles node addition with automatic data migration

### 3. **Hash Ring** (`internal/coordinator/hashRing.go`)
- **Node Management**: Add/remove nodes from the hash ring
- **Data Distribution**: Uses consistent hashing to minimize data movement
- **Migration Logic**: Automatically redistributes data when topology changes

### 4. **Server** (`cmd/server/main.go`)
- **HTTP API**: RESTful endpoints for client interaction
- **Request Coordination**: Delegates operations to appropriate cache nodes
- **JSON Serialization**: Handles request/response formatting

### 5. **Utilities** (`util/hash.go`)
- **SHA256 Hashing**: Generates 32-bit hash values for consistent distribution

## Configuration

### Cache Node Configuration
- **Default Port**: 50051 (configurable via `--port` flag)
- **Default Capacity**: 100 items per node
- **Eviction Policy**: LRU (Least Recently Used)

### Server Configuration
- **Default Port**: 8080 (configurable via `--port` flag)
- **Default Cache Nodes**: localhost:50051, localhost:50052, localhost:50053

### Makefile Configuration
```makefile
CACHE_PORTS = 50051 50052 50053  # Cache node ports
SERVER_PORT = 8080               # Coordinator server port
```

