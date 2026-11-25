TitanDB

<p align="center">
  <strong>An experimental key-value storage system built with Go.</strong>
  <br>
  Exploring networking, binary protocols, concurrency, storage engineering, observability, and distributed consensus.
</p><p align="center">
  <img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22+">
  <img src="https://img.shields.io/badge/Status-Experimental-orange?style=for-the-badge" alt="Experimental">
  <img src="https://img.shields.io/badge/Storage-Key--Value-blue?style=for-the-badge" alt="Key-Value Storage">
</p><p align="center">
  <img src="https://img.shields.io/badge/Networking-TCP-informational?style=flat-square" alt="TCP">
  <img src="https://img.shields.io/badge/Protocol-Binary-blueviolet?style=flat-square" alt="Binary Protocol">
  <img src="https://img.shields.io/badge/Concurrency-RWMutex-success?style=flat-square" alt="Concurrency">
  <img src="https://img.shields.io/badge/Consensus-Raft-orange?style=flat-square" alt="Raft">
  <img src="https://img.shields.io/badge/Observability-MVP%20Planned-yellow?style=flat-square" alt="Observability planned">
</p>«[!WARNING]
TitanDB is experimental and not production-ready. The current implementation uses a TCP-based request path backed by an in-memory key-value store. Raft replication, majority-based commits, durable storage, and end-to-end multi-node consensus are not implemented. Observability is planned for the single-node MVP. Do not use TitanDB to store production or valuable data.»

---

Contents

- "Overview" (#overview)
- "Features" (#features)
- "Project status" (#project-status)
- "Architecture" (#architecture)
- "Quick start" (#quick-start)
- "Binary protocol" (#binary-protocol)
- "Observability" (#observability)
- "Repository structure" (#repository-structure)
- "Development roadmap" (#development-roadmap)
- "Testing" (#testing)
- "Limitations" (#limitations)
- "Non-goals" (#non-goals)
- "Contributing" (#contributing)
- "License" (#license)
- "Author" (#author)

---

Overview

TitanDB is a learning-oriented systems engineering project focused on building a key-value storage service from lower-level primitives.

Instead of relying on an HTTP framework or a managed database, the project explores the internal components behind a storage system:

- TCP networking and connection handling.
- Custom binary protocol design.
- Concurrent data access.
- Storage management.
- Operational observability.
- Distributed consensus with Raft.

The project is being developed incrementally, with correctness, testing, and architectural clarity taking priority over premature optimization.

Features

Currently implemented

- TCP server prototype.
- Custom binary request framing.
- CLI client for interacting with the server.
- "PUT", "GET", and "DELETE" operations.
- Thread-safe in-memory key-value storage using "sync.RWMutex".
- Experimental "3SUM" operation.
- Initial Raft node and election-related code.

Planned for the single-node MVP

- Structured operational logging.
- Server health checks.
- Request and operation metrics.
- Error tracking and request latency.
- Active connection monitoring.
- Storage statistics.
- Automated unit and integration tests.

Planned for future development

- Persistent storage and write-ahead logging.
- Raft log replication.
- Majority-based commit.
- Multi-node consensus.
- Crash recovery.
- Cluster-level observability.

Project status

The following table distinguishes the current implementation from planned capabilities.

Component| Status
TCP server| <img src="https://img.shields.io/badge/Prototype-Implemented-success?style=flat-square">
Binary request framing| <img src="https://img.shields.io/badge/Prototype-Implemented-success?style=flat-square">
CLI client| Present
"PUT", "GET", "DELETE"| Local in-memory operations
Thread-safe storage| Implemented
"3SUM"| Experimental
Raft election scaffold| In progress
Structured logging| Existing utilities; coverage to be validated
Health checks| Planned
Request metrics| Planned
Storage statistics| Planned
Monitoring endpoints| Planned
Raft log replication| Not implemented
Majority-based commit| Not implemented
WAL and snapshots| Not implemented
Durable storage and recovery| Not implemented

Current data behavior: successful writes exist only in the local process memory. They are lost when the server exits.

---

Architecture

Current implementation

The current request path is intentionally simple:

       CLI CLIENT
           │
           │ TCP
           ▼
     TITAN RPC SERVER
      Binary Framing
           │
           ▼
    IN-MEMORY STORAGE
      map + RWMutex

The repository also contains an experimental "internal/raft" package. It is not currently integrated into the application write path.

Application writes are applied directly to the local store instead of being replicated and committed through a Raft cluster.

Planned single-node MVP

The next architectural milestone introduces an observability layer alongside the existing request path.

       CLI CLIENT
           │
           │ TCP
           ▼
     TITAN RPC SERVER
           │
           ▼
    IN-MEMORY STORAGE
           │
           ▼
  OBSERVABILITY COLLECTOR
           │
           ▼
   MONITORING INTERFACE
      ├── Health Checks
      ├── Request Metrics
      ├── Storage Statistics
      └── Operational Logs

The monitoring interface is planned to use a separate HTTP listener. This keeps monitoring traffic independent from the application's binary TCP protocol.

Package responsibilities

Package| Responsibility
"cmd/server"| Starts the TitanDB server and supporting components.
"cmd/client"| Sends key-value commands over TCP.
"internal/rpc"| Handles application request encoding, decoding, and TCP connections.
"internal/kv"| Implements the in-memory key-value store.
"internal/raft"| Contains experimental Raft election and heartbeat logic.
"internal/observability"| Planned health checks, metrics, and operational statistics.
"pkg/logger"| Provides project logging utilities.

---

Quick start

The following instructions run the current single-node prototype.

They demonstrate the TCP request path and local in-memory operations, not persistence or distributed consensus.

Requirements

- Go 1.22 or later.
- Git.
- A local terminal.

The current Go module is located inside the repository's "titan-db/" directory.

1. Clone the repository

git clone https://github.com/kauand3v/titan-db.git
cd titan-db/titan-db

2. Start the server

go run ./cmd/server/main.go \
  --id=1 \
  --address=127.0.0.1:8001 \
  --cluster=1@127.0.0.1:8001

Keep the server running.

3. Write a value

Open another terminal in the Go module directory:

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

6. Run the experimental 3SUM operation

The "3SUM" command reads integer values from the in-memory store and searches for unique triplets whose sum is zero.

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=3SUM

This is an experimental operation and is not part of the core key-value interface.

---

Binary protocol

TitanDB currently uses a custom binary application protocol.

Request header

The current request header has a fixed 9-byte layout:

Field| Size| Description
Payload length| 4 bytes| Payload length in big-endian order.
Message ID| 4 bytes| Request identifier.
Command| 1 byte| Operation code.
Payload| Variable| Command-specific data.

Command codes

Code| Command| Payload
"1"| "PUT"| Key length, key, value length, value.
"2"| "GET"| Key length and key.
"3"| "DELETE"| Key length and key.
"4"| "3SUM"| No command-specific payload.

All multi-byte integer fields use big-endian byte order.

The protocol is still evolving and is not a stable public wire-protocol specification.

Planned protocol improvements

- Explicit protocol versioning.
- Maximum payload size validation.
- Defined request and response formats.
- Error codes and malformed-message handling.
- Clear message ID correlation semantics.
- Automated encoding and decoding tests.

---

Observability

Observability is a planned part of the single-node MVP.

The goal is to make TitanDB easier to inspect, debug, and understand while running, providing visibility into server health, request processing, and storage behavior.

MVP capabilities

Capability| Description
Health checks| Report whether the local server is operational.
Request counters| Track the total number of processed requests.
Operation metrics| Track "PUT", "GET", and "DELETE" operations.
Error counters| Track failed and invalid requests.
Request latency| Measure request processing time.
Active connections| Track active TCP connections.
Storage statistics| Report key count and estimated memory usage.
Structured logging| Produce consistent operational logs.

Monitoring endpoints

The planned monitoring interface will expose a separate HTTP listener:

Endpoint| Purpose
"GET /health"| Local server health information.
"GET /metrics"| Operational metrics and storage statistics.

These endpoints are not currently implemented.

Local server health must remain separate from distributed consensus health. A running TCP server does not mean that a healthy Raft cluster exists.

Future observability

After the MVP:

- Prometheus-compatible metrics.
- Raft leader and follower status.
- Current Raft term.
- Replication lag.
- Pending log entries.
- Cluster health information.
- A dashboard for node and cluster inspection.

---

Repository structure

Current repository layout:

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

The long-term structural goal is to move the Go module to the repository root and organize optional language clients under a dedicated "clients/" directory.

---

Development roadmap

The roadmap prioritizes correctness and a demonstrable single-node MVP before distributed consensus and performance optimization.

Phase 1 — Core foundation

- [ ] Move the Go module to the repository root.
- [ ] Update the module path to "github.com/kauand3v/titan-db".
- [ ] Add protocol encoding and decoding unit tests.
- [ ] Add key-value storage unit tests.
- [ ] Add TCP server and CLI integration tests.
- [ ] Improve connection handling and malformed-request validation.
- [ ] Validate concurrent access with "go test -race ./...".

Phase 2 — Single-node MVP and observability

- [ ] Introduce the observability package.
- [ ] Implement basic server health checks.
- [ ] Add request counters and operation metrics.
- [ ] Track errors and request latency.
- [ ] Track active connections.
- [ ] Expose storage statistics.
- [ ] Improve structured logging.
- [ ] Add a separate monitoring listener.
- [ ] Implement "/health" and "/metrics".
- [ ] Add tests for health reporting and metrics collection.
- [ ] Document reproducible MVP demonstrations.

MVP milestone: a functional single-node TCP key-value service with automated tests, structured logging, health checks, and operational metrics.

This milestone does not imply persistence, replication, or production readiness.

Phase 3 — Transport and protocol reliability

- [ ] Define a versioned transport framing format or separate application and Raft listeners.
- [ ] Define maximum message and payload sizes.
- [ ] Define response formats and error codes.
- [ ] Integrate Raft message handling with the networking layer.
- [ ] Test incomplete messages, malformed frames, and connection failures.

Phase 4 — Persistent storage

- [ ] Design and implement a write-ahead log.
- [ ] Persist required Raft metadata and log entries.
- [ ] Implement snapshots.
- [ ] Implement log compaction.
- [ ] Recover node state and committed data after restarts.
- [ ] Test interrupted writes and recovery scenarios.

Phase 5 — Raft consensus

- [ ] Complete "RequestVote" term and log-freshness rules.
- [ ] Implement "AppendEntries" with previous-log validation and actual replication.
- [ ] Implement majority-based commit.
- [ ] Apply committed entries to the state machine.
- [ ] Integrate application writes with Raft.
- [ ] Validate leader and follower behavior.

Phase 6 — Distributed reliability

- [ ] Add multi-node integration tests.
- [ ] Verify leader election and log replication.
- [ ] Test leader failure and reelection.
- [ ] Test follower restarts and recovery.
- [ ] Simulate network partitions.
- [ ] Verify committed writes survive supported failure scenarios.
- [ ] Add cluster-level health and replication metrics.

Phase 7 — Performance and security

- [ ] Add throughput, latency, and resource benchmarks.
- [ ] Analyze storage and networking bottlenecks.
- [ ] Add authentication and authorization.
- [ ] Add transport security.
- [ ] Implement configurable resource limits.
- [ ] Improve operational hardening.

Phase 8 — Developer experience

- [ ] Add Docker Compose examples.
- [ ] Improve CLI diagnostics.
- [ ] Publish architecture documentation.
- [ ] Provide reproducible demonstrations.
- [ ] Add CI checks and release workflows.
- [ ] Evaluate a monitoring dashboard.

---

Testing

Automated testing is a required part of TitanDB's development.

Run the following commands from the Go module directory:

go test ./...
go test -race ./...
go vet ./...

CI checks should eventually run these commands on relevant changes.

A successful build or a test run with no test cases is not proof of Raft correctness. Distributed-systems claims require multi-process integration tests covering replicated state, committed writes, recovery, and failure behavior.

---

Limitations

The current implementation has several important limitations:

1. In-memory storage: data is held in a Go map protected by a read/write mutex. The store's "Close" method does not persist data.
2. No durable log: there is no write-ahead log or snapshot-based recovery.
3. No replicated write path: the RPC server applies "PUT" and "DELETE" directly to the local store. The Raft "SubmitCommand" method is not used by the request handler.
4. Raft transport is not integrated: the experimental Raft code uses a 5-byte message header, while application RPC uses a 9-byte header. Accepted application connections are not dispatched to "HandleRaftConnection".
5. Incomplete AppendEntries: the current handler does not validate or append log entries; it updates node state and sends a response.
6. No production security: authentication, authorization, transport encryption, resource limits, and operational hardening are not implemented.
7. No stable compatibility contract: the binary protocol may change as the project evolves.
8. No observability implementation: health checks, metrics endpoints, and storage statistics remain planned MVP capabilities.

Multi-node cluster and failover behavior should not be presented as supported until automated integration tests demonstrate it.

---

Non-goals

The current prototype is not intended to:

- Replace Redis, etcd, FoundationDB, or another production database.
- Provide production-grade durability or availability.
- Offer a stable client SDK or wire protocol.
- Claim linearizability or fault tolerance before those properties are tested.
- Compete with mature database systems without a specific technical advantage.

Contributing

Issues and focused contributions are welcome.

Before implementing a major feature, open an issue describing:

- The problem being addressed.
- The expected behavior.
- The design and relevant trade-offs.
- The tests that will demonstrate correctness.

Keep pull requests small and include tests for protocol, storage, concurrency, observability, and failure handling.

License

No "LICENSE" file is currently present in the repository. Until a license is added, the project should not be assumed to grant permission to reuse, modify, or redistribute its code.

Author

Kauan Oliveira — "@kauand3v" (https://github.com/kauand3v)

---

<p align="center">
  <sub>TitanDB is an educational systems engineering project under active development. Its documentation distinguishes implemented functionality from planned capabilities. Features are considered implemented only when they exist in the source code and are covered by reproducible tests.</sub>
</p>