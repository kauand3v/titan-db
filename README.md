TitanDB
An experimental key-value storage prototype in Go, built to explore TCP networking, binary protocols, concurrency, and distributed consensus.

Warning

Experimental project — not production-ready. The current implementation provides a TCP-based key-value request path backed by an in-memory map. The Raft package is still a prototype: replicated log entries, majority-based commits, durable storage, and end-to-end multi-node consensus are not implemented. Do not use TitanDB to store production or valuable data.

Overview
TitanDB is a learning-oriented systems project focused on understanding how a small storage service can be built from lower-level primitives instead of relying on an HTTP framework or a managed database.

The current Go implementation includes:

A TCP server and a custom binary request protocol.
A command-line client for PUT, GET, DELETE, and an experimental 3SUM operation.
An in-memory key-value store protected by sync.RWMutex.
An early Raft implementation scaffold containing node roles, election timers, vote-request handling, and heartbeat-related code.
The project is being developed incrementally. Its current state should be understood as a prototype for learning and experimentation, not as a complete distributed database.

Project status
Component	Current status
TCP server	Implemented prototype
Binary request framing	Implemented prototype
CLI client	Present
PUT, GET, DELETE	Operate against the local in-memory store
Thread-safe in-memory map	Implemented with sync.RWMutex
3SUM command	Experimental demo operation
Raft node roles and election scaffold	Work in progress; not validated end to end
Raft AppendEntries log replication	Not implemented
Majority-based commit	Not implemented
Write-ahead log and snapshots	Not implemented
Crash recovery and durable storage	Not implemented
Production security and operational hardening	Not implemented
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
The repository also contains an experimental internal/raft package. It is not yet connected to the normal application request path as a consensus layer. In particular, incoming application writes are currently applied directly to the local store rather than being appended to a replicated Raft log and committed by a cluster majority.

Main packages
Package	Responsibility
cmd/server	Starts a TitanDB server node and its supporting components.
cmd/client	Sends key-value commands to a server over TCP.
internal/rpc	Encodes and decodes application requests and handles TCP connections.
internal/kv	Provides the in-memory key-value store.
internal/raft	Contains the experimental Raft node and election/heartbeat code.
pkg/logger	Provides project logging utilities.
Binary protocol
The application request header currently uses a fixed 9-byte layout:

Field	Size	Description
Payload length	4 bytes	Length of the request payload, big-endian.
Message ID	4 bytes	Request identifier.
Command	1 byte	Operation code.
Payload	Variable	Command-specific data.
Current command codes:

Code	Command	Payload
1	PUT	Key length, key, value length, value.
2	GET	Key length and key.
3	DELETE	Key length and key.
4	3SUM	No command-specific payload.
All multi-byte integer fields are encoded in big-endian byte order. The protocol is a work in progress; it should not yet be treated as a stable public wire-protocol specification.

Requirements
Go 1.22 or later.
A local terminal.
The current Go module is located in the titan-db/ subdirectory of this repository.

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
6. Run the experimental 3SUM operation
The 3SUM command reads values from the in-memory store, parses values that are integers, and searches for unique triplets whose sum is zero.

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
    │   └── rpc/
    ├── pkg/
    │   └── logger/
    └── go.mod
The layout is being simplified. The long-term goal is to move the Go module to the repository root and keep optional language clients under a clearly named clients/ directory.

Current limitations
The following limitations are known from the current source tree:

In-memory data only. The key-value store uses a Go map protected by a read/write mutex. Its Close method does not persist data.
No durable log. There is no write-ahead log or snapshot-based recovery path.
No replicated write path. The RPC server applies PUT and DELETE directly to its local store. The Raft SubmitCommand method is not used by that request handler.
Raft transport is not integrated. The experimental Raft code uses a 5-byte message header, while the application RPC server decodes a 9-byte header. The server's accepted-connection handler does not dispatch connections to HandleRaftConnection.
AppendEntries is incomplete. The current handler does not validate or append log entries; it updates node state and sends a response.
No production security. Authentication, authorization, transport encryption, resource limits, and operational hardening are not implemented.
No stable compatibility contract. The binary protocol may change while the project evolves.
Because of these limitations, the multi-node cluster and failover scenarios described in earlier documentation are not currently supported claims. They should only be documented as working after automated integration tests demonstrate them.

Development roadmap
The roadmap is deliberately ordered around correctness before performance.

 Move the Go module to the repository root and update the module path to github.com/kauand3v/titan-db.
 Add unit tests for protocol encoding/decoding and key-value operations.
 Add integration tests for the TCP server and CLI client.
 Define a single, versioned transport framing format for client and Raft traffic, or use separate listeners with explicit protocol handling.
 Integrate Raft message handling with the server's networking layer.
 Complete RequestVote handling, including term and log-freshness rules.
 Implement AppendEntries with previous-log validation and actual log replication.
 Implement majority-based commit and apply committed entries to the state machine.
 Add a write-ahead log and snapshots.
 Recover node state and committed data after process restarts.
 Add multi-node integration tests for leader election, replication, leader failure, and network partitions.
 Add benchmarks for throughput, latency, and resource usage.
 Add authentication, transport security, and configurable resource limits.
Testing and quality checks
Automated test coverage is a required next step. As test packages are added, use these checks from the Go module directory:

go test ./...
go test -race ./...
go vet ./...
Do not interpret a successful build or a test run with no test cases as proof of Raft correctness. The distributed-systems claims require multi-process integration tests that verify replicated state and failure behavior.

Non-goals for the current prototype
Replacing Redis, etcd, FoundationDB, or another production database.
Providing production-grade durability or availability.
Offering a stable client SDK or wire protocol.
Claiming linearizability or fault tolerance before those properties are tested.
Contributing
Issues and focused contributions are welcome. Before implementing a major feature, open an issue describing:

The problem being addressed.
The expected behavior.
The design and relevant trade-offs.
The tests that will demonstrate correctness.
Keep pull requests small and include tests for protocol, storage, concurrency, and failure-handling changes.

License
No LICENSE file is currently present in the repository. Until a license is added, the project should not be assumed to grant permission to reuse, modify, or redistribute its code.

Author
Kauan Oliveira — @kauand3v

TitanDB is an educational systems project under active development. Its documentation describes the current prototype and planned work separately; features are only promoted to “implemented” after they are present in code and covered by reproducible tests.