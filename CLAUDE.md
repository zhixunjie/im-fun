# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

im-fun is an IM (Instant Messaging) system inspired by [goim](https://github.com/Terry-Mao/goim). It supports single/group chat, unicast/multicast/broadcast, and custom WebSocket protocol parsing.

## Build & Run Commands

```bash
# Build all three services (outputs to target/)
make build

# Run all tests
make test
# or with Go directly:
go test -v ./...

# Run a single test file/package
go test -v ./internal/logic/tests/...

# Run services (after build)
make run

# Stop services
make stop

# Clean build artifacts
make clean
```

**Run individual services directly (development):**
```bash
go run cmd/single/comet/main.go -conf=cmd/single/comet/comet.yaml
go run cmd/single/logic/main.go -conf=cmd/single/logic/logic.yaml
go run cmd/single/job/main.go  -conf=cmd/single/job/job.yaml
```

**Regenerate Wire dependency injection code:**
```bash
cd cmd/single/logic/wire && wire
```

**Regenerate protobuf files (run from `api/` directory):**
```bash
cd api && bash cmd.sh
```

## Architecture

The system has three services that work together:

```
Client
  |
  v (TCP :12571 or WebSocket :12572)
Comet  <----gRPC---- Job  <----Kafka---- Logic
  |                                        |
  +-----gRPC----------------------------> Logic
```

### Comet (`internal/comet/`, `cmd/single/comet/`)
Long-connection server. Manages persistent TCP/WebSocket connections from clients.
- **Bucket**: Sharded map of Channels/Rooms (hash-distributed by connection key for concurrency). Number of buckets is configurable.
- **Channel** (`internal/comet/channel/`): Represents one client TCP connection; has a ring buffer for pending messages.
- **Room**: Groups channels for multicast (e.g., group chat rooms).
- **Round**: Pool of timers and read/write buffers reused across connections to reduce allocations.
- Exposes gRPC server (`:12570`) for Job to push messages to connected clients.
- Registers itself in etcd for service discovery by Job.

### Logic (`internal/logic/`, `cmd/single/logic/`)
Business logic server. Stateless; handles auth, message storage, and routing decisions.
- **HTTP API** (`internal/logic/api/http/`): Gin-based REST endpoints for send message, fetch messages, contacts, users, groups.
- **gRPC API** (`internal/logic/api/grpc/`): Receives events from Comet (connect, disconnect, heartbeat, message receive).
- **Biz layer** (`internal/logic/biz/`): Use cases (MessageUseCase, GroupMessageUseCase, ContactUseCase, UserUseCase, UserGroupUseCase). Assembled via Google Wire.
- **Data layer** (`internal/logic/data/`): MySQL (GORM + gorm/gen), Redis, Kafka producer. Generated models/queries live in `internal/logic/data/ent/generate/`.
- After processing a message, Logic publishes to Kafka for Job to deliver.

### Job (`internal/job/`, `cmd/single/job/`)
Kafka consumer that bridges Logic and Comet.
- Consumes Kafka messages and routes them based on type: `ToUsers`, `ToRoom`, or `ToAll`.
- Watches etcd for live Comet instances and maintains gRPC connections to each (`invoker.CometInvoker`).
- For `ToRoom`, maintains `RoomJob` with batching to reduce gRPC calls.

## Key Design Patterns

- **Dependency Injection**: Google Wire is used in the Logic service. The generated file is `cmd/single/logic/wire/wire_gen.go`; edit `wire.go` and re-run `wire` to regenerate.
- **Protocol**: Custom binary protocol defined in `api/protocol.proto` (ver, op, seq, body). Connection I/O is in `api/protocol/connection_io.go`.
- **Service Discovery**: etcd via `go-kratos` registry (`pkg/registry/`). Comet registers on startup; Job watches for changes.
- **Cluster mode**: `cmd/cluster/` contains configs for running multiple instances of each service.
- **Ports (single mode)**:
  - Comet: gRPC `:12570`, TCP `:12571`, WebSocket `:12572`, pprof `:6060`, Prometheus `:7060`
  - Logic: pprof `:6062`, Prometheus `:7062`
  - Job: pprof `:6061`, Prometheus `:7061`

## Configuration

Each service reads a YAML config file via `-conf` flag. Config files are at:
- `cmd/single/comet/comet.yaml`
- `cmd/single/logic/logic.yaml`
- `cmd/single/job/job.yaml`

The `Makefile` currently targets `darwin/arm64`. Change the `ARCH` variable at the top for cross-compilation.

## Package Highlights

- `pkg/websocket/`: Custom WebSocket implementation (frame read/write, upgrade).
- `pkg/buffer/`: Custom buffer pools for memory reuse.
- `pkg/kafka/`: Sarama-based sync producer and consumer group wrappers.
- `pkg/goredis/`: Redis pool and distributed spin lock.
- `pkg/logging/`: Logrus-based logging.
- `pkg/registry/`: etcd-backed kratos registry helpers.