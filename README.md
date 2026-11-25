TitanDB

An experimental distributed key-value storage system in Go, exploring TCP networking, binary protocols, concurrency, storage engineering, observability, and distributed consensus.

"Go" (https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
"TCP" (https://img.shields.io/badge/Networking-TCP-2563EB)
"Protocol" (https://img.shields.io/badge/Protocol-Custom%20Binary-7C3AED)
"Status" (https://img.shields.io/badge/Status-Experimental-orange)

«[!WARNING]
Experimental project — not production-ready. TitanDB currently provides a TCP-based key-value request path backed by an in-memory map. Its Raft implementation is still a prototype: replicated log entries, majority-based commits, durable storage, and end-to-end multi-node consensus are not implemented.»

---

The Problem

Building a storage system is much more than implementing "PUT" and "GET" operations:

- Network communication is complex. Applications need reliable request framing, message identification, and a protocol for exchanging data over TCP.
- Concurrent access introduces risks. Multiple clients accessing shared data require synchronization to prevent race conditions and inconsistent state.
- Data durability is essential. In-memory storage is fast, but data disappears when the process terminates.
- Distributed consistency is difficult. Multiple nodes need a coordination mechanism to agree on which operations are committed.
- Operational visibility matters. Logging, observability, and performance measurements are necessary to understand system behavior.

The Solution

TitanDB is an experimental project designed to explore the internal components of a networked key-value storage system:

- Custom TCP communication: a binary request protocol for communication between clients and the server.
- Key-value operations: local "PUT", "GET", and "DELETE" operations backed by an in-memory storage engine.
- Concurrent access: synchronized access to shared data using Go's "sync.RWMutex".
- Command-line client: a CLI for sending requests to the TitanDB server.
- Experimental query operation: a "3SUM" operation that searches stored integer values for triplets whose sum is zero.
- Raft prototype: an early implementation exploring node roles, election timers, and vote requests as a foundation for distributed consensus.

Architecture

TitanDB currently uses a simple client-server architecture. The Raft package is separate from the application request path and is not yet responsible for committing key-value operations.

flowchart TB
    Client["CLI Client"]

    subgraph Server["TitanDB Server"]
        TCP["TCP Connection"]
        RPC["Binary RPC Protocol"]
        Handler["Request Handler"]
        KV["In-Memory Key-Value Store"]
        Mutex["sync.RWMutex"]
    end

    subgraph Experimental["Experimental Components"]
        Raft["Raft Prototype"]
        Roles["Node Roles"]
        Election["Election Timers"]
        Vote["RequestVote Handling"]
    end

    Client -->|"TCP Requests"| TCP
    TCP --> RPC
    RPC --> Handler
    Handler --> KV
    KV --> Mutex

    Raft --> Roles
    Raft --> Election
    Raft --> Vote

Request flow: the CLI client establishes a TCP connection with the server and sends a binary-framed request. The RPC layer decodes the message, and the request handler executes the corresponding operation against the local in-memory key-value store.

Important: key-value writes are currently applied directly to the local store. They are not replicated through Raft, and they are not committed through a majority of nodes.

Current Features

«Keep this list limited to functionality documented in the current implementation. Experimental components are explicitly identified.»

- TCP-based client-server communication
- Custom binary request framing
- Message identification through request IDs
- CLI client for interacting with the server
- Local "PUT", "GET", and "DELETE" operations
- In-memory key-value storage
- Concurrent access protection using "sync.RWMutex"
- Experimental "3SUM" operation over stored integer values
- Initial Raft node structure and role definitions
- Election timer and "RequestVote" prototype
- Logging package for application observability

Tech Stack

Layer| Technology| Purpose
Core language| Go 1.22+| System implementation and concurrency
Networking| TCP| Client-server communication
Application protocol| Custom binary protocol| Structured request framing
Storage engine| Go map| In-memory key-value storage
Synchronization| "sync.RWMutex"| Concurrent access protection
Distributed consensus| Raft prototype| Exploration of leader election and consensus
CLI| Go| Client-side request execution
Logging| Custom logger package| Application logging

Quick Start

Prerequisites: Go 1.22+ and Git.

1. Clone the repository

git clone https://github.com/kauand3v/titan-db.git
cd titan-db/titan-db

2. Start the TitanDB server

go run ./cmd/server/main.go \
  --id=1 \
  --address=127.0.0.1:8001 \
  --cluster=1@127.0.0.1:8001

The server will start listening on the configured TCP address.

3. Send a PUT request

Open another terminal from the module directory:

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=PUT \
  --key=greeting \
  --val=hello

4. Retrieve a value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=GET \
  --key=greeting

5. Delete a value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=DELETE \
  --key=greeting

6. Run the experimental 3SUM operation

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=3SUM

«[!NOTE]
The "--cluster" argument does not mean that distributed consensus is operational. The current key-value request path does not use Raft replication.»

Binary Protocol

TitanDB uses a custom binary protocol over TCP.

Each application request begins with a fixed-size header of 9 bytes:

Field| Size| Description
Payload length| 4 bytes| Request payload length field
Message ID| 4 bytes| Request identifier
Command| 1 byte| Operation identifier
Payload| Variable| Command-specific data

All multi-byte integer fields use big-endian byte order.

Supported commands

Command| ID| Description
"PUT"| "1"| Stores a key-value pair
"GET"| "2"| Retrieves a value by key
"DELETE"| "3"| Removes a key-value pair
"3SUM"| "4"| Searches stored integer values for zero-sum triplets

The protocol is experimental and is not intended to be a stable or externally supported wire-protocol specification.

Storage and Consensus

Key-value storage

The current storage implementation uses an in-memory Go map protected by "sync.RWMutex".

This provides synchronized access to the local store, but it does not provide persistence.

- Data is held in process memory.
- Writes are local to the server instance.
- Restarting the process loses stored data.
- There is no write-ahead log or snapshot-based recovery.

Raft implementation

The "internal/raft" package is an experimental foundation for future distributed consensus functionality.

The current prototype includes:

- Node role definitions.
- Election timer-related logic.
- Initial "RequestVote" handling.
- Heartbeat-related code.

The following capabilities are not implemented as a complete, operational consensus system:

- Replicated log entries.
- Functional "AppendEntries" log replication.
- Majority-based commit decisions.
- Raft-backed key-value writes.
- Durable Raft state and log storage.
- Snapshotting and crash recovery.
- End-to-end multi-node consensus.

Project Structure

titan-db/
├── .vscode/
├── ruby/
├── titan-db/
│   ├── cmd/
│   │   ├── client/       # CLI client
│   │   └── server/       # TCP server
│   ├── internal/
│   │   ├── kv/           # Key-value storage
│   │   ├── raft/         # Experimental Raft implementation
│   │   └── rpc/          # Binary RPC protocol
│   ├── pkg/
│   │   └── logger/       # Logging utilities
│   └── go.mod
├── .gitignore
└── README.md

Testing and Development

From the Go module directory:

Run tests

go test ./...

Run tests with the race detector

go test -race ./...

Run static analysis

go vet ./...

These commands are development checks, not evidence of complete distributed-system correctness. In particular, Raft behavior requires dedicated multi-node tests before it can be considered reliable.

Known Limitations

TitanDB is an engineering and learning project. Its current implementation has important limitations:

- No persistent key-value storage.
- No write-ahead log or crash recovery.
- No functional Raft log replication.
- No majority-based commit mechanism.
- Key-value requests bypass the Raft package.
- No complete multi-node consensus workflow.
- The application RPC and experimental Raft transport use different message framing.
- No production-grade authentication, authorization, or TLS.
- No comprehensive resource limits or production hardening.

Roadmap

The roadmap focuses on evolving TitanDB from an experimental TCP key-value prototype into a more complete storage system.

- [ ] Testing foundation: unit tests for storage, RPC framing, and server behavior.
- [ ] Protocol engineering: document and standardize message framing, validation, and error handling.
- [ ] Storage persistence: implement a write-ahead log and durable storage.
- [ ] Raft integration: connect the application write path to the consensus layer.
- [ ] Log replication: implement functional "AppendEntries" and replicated log management.
- [ ] Consensus and commits: implement majority-based commit decisions and state machine application.
- [ ] Crash recovery: add snapshots, log replay, and recovery mechanisms.
- [ ] Multi-node testing: validate elections, replication, failure handling, and recovery across multiple nodes.
- [ ] Observability: improve structured logging, metrics, and operational diagnostics.
- [ ] Performance: add benchmarks for storage operations, protocol overhead, and concurrent workloads.
- [ ] Security and hardening: introduce connection limits, request validation, and transport security.

Contributing

Issues, suggestions, and pull requests are welcome.

For substantial changes, open an issue first to discuss the proposed implementation and its impact on the architecture.

License

No license has currently been declared for this project. Until a license is added, reuse, modification, and redistribution permissions should not be assumed.

---

Developed by "Kauan Oliveira" (https://github.com/kauand3v)

"GitHub Repository" (https://github.com/kauand3v/titan-db)