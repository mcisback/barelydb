# BarelyDB

A lightweight, JSON-based database with a RESTful API. Think of it as a simpler version of MongoDB that's incredibly easy to set up and run.

## What is BarelyDB?

BarelyDB is a minimalist database system that stores data as JSON files on disk and provides a simple HTTP API for accessing your data. Perfect for prototyping, small projects, or when you need a quick database solution without the overhead of traditional database systems.
Or when you need a quick and fast database for startups.

## Features

- **Zero configuration** - Just run the binary and start using it
- **JSON storage** - All data stored as human-readable JSON files
- **RESTful API** - Simple HTTP endpoints for data access
- **Lightweight** - Minimal dependencies, fast startup
- **File-based** - Easy to backup, version control, and inspect

## Installation

### Prerequisites

- Go 1.16 or higher

### Build from source

```bash
go build -o barelydb
```

## Quick Start

1. Run BarelyDB:
```bash
./barelydb
```

The server will start on port `3838` and create a root database directory automatically.

2. Access your data via HTTP:
```bash
# Get all records from a table
curl http://localhost:3838/mydb/users

# Get a specific record by ID
curl http://localhost:3838/mydb/users/user123
```

## API Reference

### Get All Records

Retrieve all records from a table.

**Endpoint:** `GET /:database/:table`

**Example:**
```bash
curl http://localhost:3838/myapp/users
```

**Response:**
```json
{
  "user1": {
    "name": "John Doe",
    "email": "john@example.com"
  },
  "user2": {
    "name": "Jane Smith",
    "email": "jane@example.com"
  }
}
```

### Get Single Record

Retrieve a specific record by its ID.

**Endpoint:** `GET /:database/:table/:id`

**Example:**
```bash
curl http://localhost:3838/myapp/users/user1
```

**Response:**
```json
{
  "id": "user1",
  "name": "John Doe",
  "email": "john@example.com"
}
```

### Error Responses

**Database Not Found (404):**
```json
{
  "error": "Database Not Found"
}
```

**Table Not Found (404):**
```json
{
  "error": "Table Not Found"
}
```

**Record Not Found (404):**
```json
{
  "error": "Record Not Found"
}
```

## Directory Structure

BarelyDB organizes data in a simple file structure:

```
barelydb_data/
├── database1/
│   ├── table1.json
│   └── table2.json
└── database2/
    └── table1.json
```

Each table is stored as a JSON file with records keyed by ID:

```json
{
  "record1": {
    "field1": "value1",
    "field2": "value2"
  },
  "record2": {
    "field1": "value3",
    "field2": "value4"
  }
}
```

## Configuration

### Port Configuration

By default, BarelyDB listens on port `3838`. To change this, modify the `LISTEN_PORT` constant in `main.go`:

```go
const LISTEN_PORT = ":3838"
```

### Data Directory

The root database directory is automatically determined by the `getRootDatabaseDirectory()` function. Check your implementation for the specific location.

## Use Cases

BarelyDB is perfect for:

- **Rapid prototyping** - Get a database up and running in seconds
- **Development environments** - No need to install heavy database systems
- **Small applications** - Personal projects, scripts, and tools
- **Learning projects** - Understand database concepts without complexity
- **Configuration storage** - Store app settings and configurations
- **Testing** - Easily create and teardown test databases

## Limitations

BarelyDB is designed for simplicity, not production scale. Be aware of:

- No built-in authentication or authorization
- Limited query capabilities (currently no filtering, sorting, or pagination)
- File-based storage may not scale to millions of records
- No transaction support
- No concurrent write safety

## Roadmap

Future enhancements may include:

- POST/PUT/DELETE operations for creating and modifying data
- Query parameters for filtering and sorting
- Basic authentication
- Indexing for faster lookups
- Backup and restore utilities

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## License

[Add your license here]

## Why "Barely"?

Because it's *barely* a database - and that's exactly the point. Sometimes you don't need the complexity of a full database system. Sometimes you just need something that *barely* gets the job done, and does it well.

---

**Getting Started in 30 Seconds:**

```bash
# Build
go build -o barelydb

# Create database
mkdir -p ./barelydb_data/mydb
echo "\{\}" > ./barelydb_data/mydb/users.json

# Run
./barelydb

# Use
curl http://localhost:3838/mydb/users
```

That's it. You're running a database.
