TitanDB

An experimental key-value storage system in Go, exploring TCP networking, binary protocols, concurrency, storage engineering, observability, and distributed consensus.

«[!WARNING] Experimental project — not production-ready. TitanDB currently provides a TCP-based key-value request path backed by an in-memory map. The Raft package remains a prototype: replicated log entries, majority-based commits, durable storage, and end-to-end multi-node consensus are not implemented. Basic observability is a target of the single-node MVP, not an implemented feature. Do not use TitanDB to store production or valuable data.»

Overview

TitanDB is a learning-oriented systems engineering project focused on understanding how a storage service can be built from lower-level primitives instead of relying on an HTTP framework or a managed database.

The project explores the engineering foundations of a key-value storage system, including network communication, binary protocol design, concurrent access, storage management, observability, and distributed consensus.

The current Go implementation includes:

- A TCP server and a custom binary request protocol.
- A command-line client supporting "PUT", "GET", "DELETE", and an experimental "3SUM" operation.
- An in-memory key-value store protected by "sync.RWMutex".
- An early Raft implementation scaffold containing node roles, election timers, vote-request handling, and heartbeat-related code.

The planned single-node MVP will extend the current implementation with basic observability, structured logging, health checks, request metrics, storage statistics, and automated testing.

The project is being developed incrementally. Its current state should be understood as an experimental prototype, not as a complete distributed database.

Project status

Component| Current status
TCP server| Implemented prototype
Binary request framing| Implemented prototype
CLI client| Present
"PUT", "GET", "DELETE"| Operate against the local in-memory store
Thread-safe in-memory map| Implemented with "sync.RWMutex"
"3SUM" command| Experimental demo operation
Raft node roles and election scaffold| Work in progress; not validated end to end
Basic structured logging| Existing logging utilities; coverage to be validated
Health checks| Planned for the single-node MVP
Request metrics| Planned for the single-node MVP
Storage statistics| Planned for the single-node MVP
Observability endpoints| Planned for the single-node MVP
Raft "AppendEntries" log replication| Not implemented
Majority-based commit| Not implemented
Write-ahead log and snapshots| Not implemented
Crash recovery and durable storage| Not implemented
Production security and operational hardening| Not implemented

The project does not currently provide replicated or durable writes. A successful write is stored only in the process-local map and is lost when that process exits.

Architecture

Current request path

┌──────────────────┐
│   CLI Client     │
└────────┬─────────┘
         │ TCP
         ▼
┌──────────────────┐
│ TitanRPC Server  │
│ binary framing   │
└────────┬─────────┘
         │
         ▼
┌──────────────────────────┐
│ In-memory Key-Value Store│
│ map + sync.RWMutex       │
└──────────────────────────┘

The repository also contains an experimental "internal/raft" package. It is not yet connected to the normal application request path as a consensus layer. Incoming application writes are currently applied directly to the local store rather than being appended to a replicated Raft log and committed by a cluster majority.

Planned MVP architecture

The single-node MVP will introduce an observability layer alongside the existing request path. A separate monitoring interface will expose local server health and operational metrics without changing the application's binary TCP protocol.

┌──────────────────┐
│   CLI Client     │
└────────┬─────────┘
         │ TCP
         ▼
┌──────────────────┐
│ TitanRPC Server  │
└────────┬─────────┘
         │
         ▼
┌──────────────────────────┐
│ In-memory Key-Value Store│
└──────────────────────────┘
         │
         ▼
┌──────────────────────────┐
│  Observability Collector │
└────────┬─────────────────┘
         │
         ▼
┌──────────────────────────┐
│   Monitoring Interface   │
│                          │
│  Health Checks           │
│  Request Metrics         │
│  Storage Statistics      │
│  Operational Logs        │
└──────────────────────────┘

The observability layer is intended to collect information from the server and storage components. Its monitoring interface will use a separate HTTP listener, keeping monitoring traffic independent from the application's binary TCP protocol.

Main packages

Package| Responsibility
"cmd/server"| Starts a TitanDB server node and its supporting components.
"cmd/client"| Sends key-value commands to a server over TCP.
"internal/rpc"| Encodes and decodes application requests and handles TCP connections.
"internal/kv"| Provides the in-memory key-value store.
"internal/raft"| Contains the experimental Raft node and election/heartbeat code.
"internal/observability"| Planned health checks, metrics collection, and operational statistics.
"pkg/logger"| Provides project logging utilities.

Binary protocol

The application request header currently uses a fixed 9-byte layout:

Field| Size| Description
Payload length| 4 bytes| Length of the request payload, big-endian.
Message ID| 4 bytes| Request identifier.
Command| 1 byte| Operation code.
Payload| Variable| Command-specific data.

Current command codes:

Code| Command| Payload
"1"| "PUT"| Key length, key, value length, value.
"2"| "GET"| Key length and key.
"3"| "DELETE"| Key length and key.
"4"| "3SUM"| No command-specific payload.

All multi-byte integer fields are encoded in big-endian byte order.

The protocol is a work in progress and should not yet be treated as a stable public wire-protocol specification.

Future protocol improvements include:

- Explicit protocol versioning.
- Maximum payload size validation.
- Request and response definitions.
- Error codes and malformed-message handling.
- Clear message ID correlation semantics.
- Automated encoding and decoding tests.

Observability

Observability is a planned component of the single-node MVP.

The objective is to make TitanDB easier to inspect, debug, and understand while it is running, providing visibility into request handling, storage operations, and server health.

The initial implementation will prioritize a small set of operational signals rather than a complex monitoring stack.

Planned MVP capabilities

Capability| Description
Health checks| Report whether the local server is operational.
Request counters| Track the total number of processed requests.
Operation metrics| Track "PUT", "GET", and "DELETE" operations.
Error counters| Track failed and invalid requests.
Request latency| Measure request processing time.
Active connections| Track currently active TCP connections.
Storage statistics| Report key count and estimated memory usage.
Structured logging| Produce consistent operational logs.
Metrics endpoint| Expose metrics through a separate monitoring interface.

Planned monitoring endpoints

The MVP is expected to expose a separate HTTP listener for monitoring:

Endpoint| Purpose
"GET /health"| Local server health information.
"GET /metrics"| Operational metrics and storage statistics.

These endpoints are planned and are not currently available.

The health endpoint must distinguish between local server availability and distributed consensus health. Until Raft is fully integrated and validated, TitanDB must not report a functioning distributed cluster merely because its local TCP server is running.

Future observability improvements

After the MVP:

- Prometheus-compatible metrics.
- Raft leader and follower status.
- Current Raft term.
- Replication lag.
- Pending log entries.
- Cluster health information.
- Dashboard for node and cluster inspection.

Requirements

- Go 1.22 or later.
- A local terminal.

The current Go module is located in the "titan-db/" subdirectory of this repository.

Quick start: single-node prototype

The following commands exercise the current single-node TCP/RPC and in-memory key-value path. They do not demonstrate replicated writes or multi-node Raft consensus.

1. Clone the repository

git clone https://github.com/kauand3v/titan-db.git
cd titan-db/titan-db

2. Start one server node

go run ./cmd/server/main.go \
  --id=1 \
  --address=127.0.0.1:8001 \
  --cluster=1@127.0.0.1:8001

Keep this terminal running.

3. Write a value

Open a second terminal, enter the Go module directory, and run:

cd titan-db/titan-db

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=PUT \
  --key=greeting \
  --val=hello

4. Read the value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=GET \
  --key=greeting

5. Delete the value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=DELETE \
  --key=greeting

6. Run the experimental "3SUM" operation

The "3SUM" command reads values from the in-memory store, parses values that are integers, and searches for unique triplets whose sum is zero.

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=3SUM

This command is an experiment and is not part of the core key-value storage interface.

Repository layout

titan-db/
├── README.md
├── ruby/                      # Experimental Ruby client material
└── titan-db/                  # Current Go module
    ├── cmd/
    │   ├── client/
    │   └── server/
    ├── internal/
    │   ├── kv/
    │   ├── raft/
    │   ├── rpc/
    │   └── observability/     # Planned
    ├── pkg/
    │   └── logger/
    └── go.mod

The layout is being simplified. The long-term goal is to move the Go module to the repository root and keep optional language clients under a clearly named "clients/" directory.

Current limitations

The following limitations are known from the current source tree:

1. In-memory data only. The key-value store uses a Go map protected by a read/write mutex. Its "Close" method does not persist data.
2. No durable log. There is no write-ahead log or snapshot-based recovery path.
3. No replicated write path. The RPC server applies "PUT" and "DELETE" directly to its local store. The Raft "SubmitCommand" method is not used by that request handler.
4. Raft transport is not integrated. The experimental Raft code uses a 5-byte message header, while the application RPC server decodes a 9-byte header. The server's accepted-connection handler does not dispatch connections to "HandleRaftConnection".
5. AppendEntries is incomplete. The current handler does not validate or append log entries; it updates node state and sends a response.
6. No production security. Authentication, authorization, transport encryption, resource limits, and operational hardening are not implemented.
7. No stable compatibility contract. The binary protocol may change while the project evolves.
8. Observability is not implemented yet. Health checks, metrics endpoints, and storage statistics remain planned MVP capabilities.

Because of these limitations, the multi-node cluster and failover scenarios described in earlier documentation are not currently supported claims. They should only be documented as working after automated integration tests demonstrate them.

Development roadmap

The roadmap prioritizes correctness, observability, and reproducibility before distributed features and performance optimization.

Phase 1 — Core Foundation

- [ ] Move the Go module to the repository root and update the module path to "github.com/kauand3v/titan-db".
- [ ] Add unit tests for protocol encoding and decoding.
- [ ] Add unit tests for key-value operations.
- [ ] Add integration tests for the TCP server and CLI client.
- [ ] Improve connection handling and malformed-request validation.
- [ ] Validate concurrent storage access with "go test -race ./...".

Phase 2 — Single-Node MVP and Observability

- [ ] Introduce an observability package.
- [ ] Implement basic server health checks.
- [ ] Add request counters and operation metrics.
- [ ] Track request errors and processing latency.
- [ ] Track active connections.
- [ ] Expose basic storage statistics.
- [ ] Improve structured logging.
- [ ] Add a separate monitoring listener for health and metrics endpoints.
- [ ] Add tests for metrics collection and health reporting.
- [ ] Document the MVP architecture and reproducible demonstrations.

MVP milestone: a functional single-node TCP key-value service with automated tests, structured logging, basic health checks, and operational metrics.

This milestone does not imply persistence, replication, or production readiness.

Phase 3 — Transport and Protocol Reliability

- [ ] Define a single, versioned transport framing format for client and Raft traffic, or use separate listeners with explicit protocol handling.
- [ ] Define maximum message and payload sizes.
- [ ] Define response formats and error codes.
- [ ] Integrate Raft message handling with the server's networking layer.
- [ ] Add tests for incomplete messages, malformed frames, and connection failures.

Phase 4 — Persistent Storage

- [ ] Design and implement a write-ahead log.
- [ ] Persist the required Raft metadata and log entries.
- [ ] Implement snapshots.
- [ ] Implement log compaction.
- [ ] Recover node state and committed data after process restarts.
- [ ] Test interrupted writes and recovery scenarios.

Phase 5 — Raft Consensus

- [ ] Complete "RequestVote" handling, including term and log-freshness rules.
- [ ] Implement "AppendEntries" with previous-log validation and actual log replication.
- [ ] Implement majority-based commit.
- [ ] Apply committed entries to the state machine.
- [ ] Integrate the application write path with Raft.
- [ ] Validate leader and follower behavior.

Phase 6 — Distributed Reliability

- [ ] Add multi-node integration tests.
- [ ] Verify leader election.
- [ ] Verify log replication.
- [ ] Test leader failure and reelection.
- [ ] Test follower restarts and recovery.
- [ ] Simulate network partitions.
- [ ] Verify that acknowledged committed writes survive supported failure scenarios.
- [ ] Add cluster-level health and replication metrics.

Phase 7 — Performance and Security

- [ ] Add benchmarks for throughput, latency, and resource usage.
- [ ] Analyze storage and networking bottlenecks.
- [ ] Add authentication and authorization.
- [ ] Add transport security.
- [ ] Implement configurable resource limits.
- [ ] Improve operational hardening.

Phase 8 — Developer Experience

- [ ] Add Docker Compose examples.
- [ ] Improve CLI diagnostics.
- [ ] Publish architecture documentation.
- [ ] Provide reproducible demonstrations.
- [ ] Add CI checks and release workflows.
- [ ] Evaluate a monitoring dashboard for cluster inspection.

Testing and quality checks

Automated test coverage is a required part of the project's development.

Run these checks from the Go module directory:

go test ./...
go test -race ./...
go vet ./...

The project should also introduce automated CI checks to run these commands on relevant changes.

Do not interpret a successful build or a test run with no test cases as proof of Raft correctness. Distributed-systems claims require multi-process integration tests that verify replicated state, committed writes, recovery, and failure behavior.

Non-goals for the current prototype

- Replacing Redis, etcd, FoundationDB, or another production database.
- Providing production-grade durability or availability.
- Offering a stable client SDK or wire protocol.
- Claiming linearizability or fault tolerance before those properties are tested.
- Competing with mature database systems before establishing a specific technical advantage.

Contributing

Issues and focused contributions are welcome. Before implementing a major feature, open an issue describing:

- The problem being addressed.
- The expected behavior.
- The design and relevant trade-offs.
- The tests that will demonstrate correctness.

Keep pull requests small and include tests for protocol, storage, concurrency, observability, and failure-handling changes.

License

No "LICENSE" file is currently present in the repository. Until a license is added, the project should not be assumed to grant permission to reuse, modify, or redistribute its code.

Author

Kauan Oliveira — "@kauand3v" (https://github.com/kauand3v)

---

TitanDB is an educational systems engineering project under active development. Its documentation distinguishes the current implementation from planned capabilities. Features are only promoted to implemented status after they exist in the source code and are covered by reproducible tests.