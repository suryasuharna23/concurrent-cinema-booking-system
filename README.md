# Concurrent Cinema Booking System

This is a Go backend project that simulates high-concurrency seat contention in a cinema booking flow.
Its main objective is to guarantee that only one user can book the same seat, even when a very large number of requests arrive at nearly the same time.

## Goals

- Prevent double booking for the same seat.
- Implement seat hold with TTL using Redis SET NX EX.
- Test extreme contention scenarios with many goroutines.

## Tech Stack

- Go 1.25.3
- Redis 7 for distributed lock and seat hold
- Docker Compose (optional) to run Redis locally

## Project Structure

```text
.
|- cmd/
|  \- main.go
|- internal/
|  |- adapters/
|  |  \- redis/
|  |     \- redis.go
|  \- booking/
|     |- domain.go
|     |- service.go
|     |- memory_store.go
|     |- concurrent_memory_store.go
|     |- redis_store.go
|     \- service_test.go
|- static/
|  \- index.html
|- docker.compose.yaml
|- go.mod
\- README.md
```

## Core Concepts

### 1) Booking Domain
The Booking model contains:

- ID
- MovieID
- SeatID
- UserID
- Status
- ExpiresAt

Important domain error:

- ErrSeatAlreadyBooked, used when a seat is already taken.

### 2) In-Memory Stores

- MemoryStore: not thread-safe, suitable for simple demos.
- ConcurrentMemoryStore: thread-safe using sync.RWMutex for in-process concurrency.

### 3) Redis Store

RedisStore uses this strategy:

- Seat key format: seat:<movieID>:<seatID>
- SET NX so only the first request wins.
- Default hold TTL is 2 minutes through defaultHoldTTL = 2 * time.Minute.
- Session key format: session:<sessionID> for hold-session mapping.

This approach keeps seat contention behavior consistent under heavy concurrent traffic.

## Running Redis

### Option A: Docker Compose (recommended)

```bash
docker compose -f docker.compose.yaml up -d
```

Services:

- Redis on localhost:6379
- Redis Commander on http://localhost:8081

### Option B: Local Redis

Make sure Redis is running on localhost:6379.

## Running Tests

```bash
go test ./...
```

Important note:

- TestConcurrentBooking_ExactlyOneWins runs 100_000 goroutines competing for the same seat.
- Redis must be up because the test initializes RedisStore and performs PING on startup.

Expected result:

- Exactly 1 successful booking.
- The remaining requests fail because the seat is already taken.

## Current Implementation Status

- Core booking logic and Redis locking are implemented.
- Concurrency stress test is implemented.
- Static UI in static/index.html is implemented.
- HTTP API backend endpoints such as /movies and /sessions are not implemented in the current Go code.
- cmd/main.go is still a placeholder with package declaration only.

In short, the static frontend already defines an API contract, but the backend HTTP server is still pending.

## Planned Endpoints Based on Frontend Contract

- GET /movies
- GET /movies/{movieID}/seats
- POST /movies/{movieID}/seats/{seatID}/hold
- PUT /sessions/{sessionID}/confirm
- DELETE /sessions/{sessionID}

Once these endpoints are implemented, static/index.html can be used as a realtime seat-map demo.

## Suggested Next Improvements

1. Implement an HTTP server in cmd/main.go (net/http or your preferred framework).
2. Add Confirm and Release use cases in service and store layers.
3. Store confirmed seat state separately from temporary hold state.
4. Add endpoint integration tests.
5. Add CI checks for linting and test automation.