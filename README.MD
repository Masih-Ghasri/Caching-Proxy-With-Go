# In-Memory Caching Server with Go 🚀

A high-performance in-memory caching server built with Go, inspired by systems like Redis and Memcached. This server supports key features including TTL (Time To Live), LRU eviction policy, and persistence, optimized for handling multiple concurrent clients efficiently.

This project was developed as part of the [roadmap.sh Caching Server Project](https://roadmap.sh/projects/caching-server) learning path.

## Key Features ⚡️

- **Concurrent TCP Server:** Uses Goroutines to serve multiple clients simultaneously without blocking, ensuring high throughput and responsiveness.

- **Key-Value Storage:** Supports essential `GET`, `SET`, and `DELETE` operations for data management with a simple, Redis-like interface.

- **Time-To-Live (TTL):** Individual expiration time setting for each key. Expired keys are removed using two mechanisms:
    - **Passive Deletion:** When an expired key is accessed, it gets removed immediately.
    - **Active Deletion:** A background process periodically scans the entire cache to remove expired keys.

- **LRU Eviction Policy (Least Recently Used):** To prevent memory overflow, the server automatically removes the oldest and least-used items when reaching maximum capacity.

- **Persistence:** The server periodically saves a snapshot of cache data to a disk file (`cache.gob`). On restart, data is recovered from this file to prevent data loss.

- **Thread-Safe Operations:** All cache access operations are protected using Mutex to prevent race conditions in concurrent environments.

- **Modular Code Architecture:** Core cache logic is separated into a dedicated package (`/cache`) for clean, testable, and reusable code.

## Project Structure 📦

```
/
├── cache/
│   └── cache.go      # Core cache engine (LRU, TTL, Persistence logic)
├── go.mod            # Module dependencies management
├── LICENSE           # Project license (MIT)
├── main.go           # Network layer (TCP Server) & application entry point
└── README.md         # Project documentation
```

## Prerequisites 📋

- **Go:** Version 1.18 or higher
- **Port 8080:** Should be available (or modify in code)

## Installation & Usage ⚙️

### 1. Clone the Repository

```bash
git clone https://github.com/Masih-Ghasri/Caching-Proxy-With-Go.git
cd Caching-Proxy-With-Go
```

### 2. Run the Server

For development:
```bash
go run .
```

For production:
```bash
go build -o caching-server .
./caching-server
```

After running, the server will start listening on port `8080` and display:
```
Caching server started on :8080
```

## How to Use 🔧

### Connect to Server

You can connect to the server using tools like `telnet`, `netcat`, or `nc`:

```bash
telnet localhost 8080
```

or:

```bash
nc localhost 8080
```

### Supported Commands

#### SET - Store Key-Value Pair

Stores a key with a value and optional expiration time (in seconds).

**Format:** `SET key value [duration_in_seconds]`

**Examples:**
```
SET name John
> OK

SET session_token abc123 300
> OK

SET temp_data hello 10
> OK
```

- If duration is not specified, default TTL of 10 minutes is applied
- Duration is specified in seconds

#### GET - Retrieve Value

Retrieves the value associated with a key.

**Format:** `GET key`

**Examples:**
```
GET name
> John

GET nonexistent
> Key not found
```

#### DELETE - Remove Key

Removes a key from the cache.

**Format:** `DELETE key`

**Examples:**
```
DELETE name
> 1

DELETE nonexistent
> 0
```

- `1`: Key existed and was deleted
- `0`: Key did not exist

## Practical Usage Example 💡

```bash
# Connect to server
$ telnet localhost 8080

# Store a session token with 1-hour expiration
SET user:123:session abc123def456 3600
OK

# Store temporary data with 30-second expiration
SET temp:verification_code 456789 30
OK

# Retrieve session
GET user:123:session
abc123def456

# After 30 seconds, try to get verification code
GET temp:verification_code
Key not found

# Delete session
DELETE user:123:session
1
```

## Customizable Settings ⚙️

In `main.go`, you can modify these settings:

```go
// Change server port
listener, err := net.Listen("tcp", ":8080") // Change to ":9090"

// Modify cache settings
c := cache.NewCache(
    100,                // Maximum number of items
    5*time.Minute,      // Cleanup interval
)
```

In `cache/cache.go`:

```go
// Change default expiration time
var duration time.Duration = 10 * time.Minute // in main.go

// Change save interval
go c.saveLoop(5*time.Minute, SaveFilePath) // in NewCache()
```

## Advanced Features 🔥

### 1. Automatic Persistence

- Cache is automatically saved to `cache.gob` file every 5 minutes
- On startup, previous data is recovered from the file
- Expired keys are automatically removed during recovery

### 2. Smart Memory Management

- **LRU Eviction:** Maximum 100 items (configurable)
- **Active Cleanup:** Expired keys are removed every 5 minutes
- **Passive Cleanup:** Expired keys are removed immediately upon access

### 3. Comprehensive Logging

The server logs all important activities:
- Client connections and disconnections
- Expired key removals
- Save and load operations
- System errors and warnings

## Performance Testing 🚀

To test concurrency, you can connect multiple clients simultaneously:

```bash
# Terminal 1
telnet localhost 8080

# Terminal 2
telnet localhost 8080

# Terminal 3
nc localhost 8080
```

## Code Architecture 🏗️

### Cache Package (`cache/cache.go`)

- **`entry`:** Data structure for each cache item (key, value, expiration time)
- **`Cache`:** Main structure with Mutex, LRU List, and Map for fast access
- **`NewCache()`:** Constructor with cleanup and persistence settings
- **`Set()`, `Get()`, `Delete()`:** Core CRUD operations
- **`cleanupLoop()`:** Periodic removal of expired keys
- **`saveLoop()`:** Periodic file persistence

### Main Package (`main.go`)

- **TCP Server:** Network connection management
- **Protocol Handler:** Text-based command processing
- **Connection Management:** Concurrent client handling

## Contributing 🤝

This project is open source and welcomes contributions:

1. Fork the project
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## Performance Characteristics 📊

- **Concurrency:** Handles multiple clients simultaneously using Goroutines
- **Memory Usage:** O(1) access time with hash map + doubly linked list
- **Persistence:** Automatic background saves without blocking operations
- **Scalability:** Configurable cache size and cleanup intervals
- **Reliability:** Thread-safe operations with proper error handling

## License 📄

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Author 👨‍💻

**Masih Ghasri**
- GitHub: [@Masih-Ghasri](https://github.com/Masih-Ghasri)

---

