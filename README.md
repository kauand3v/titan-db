TitanDB

An experimental key-value storage prototype in Go, exploring TCP networking, binary protocols, concurrency, storage engineering, and distributed consensus.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white" alt="Go 1.22+"/>
  <img src="https://img.shields.io/badge/Status-Experimental-orange" alt="Experimental"/>
  <img src="https://img.shields.io/badge/Architecture-Client--Server-blue" alt="Client Server"/>
  <img src="https://img.shields.io/badge/Storage-In--Memory-lightgrey" alt="In Memory"/>
</p>«[!WARNING]
Experimental project — not production-ready.

TitanDB currently provides a TCP-based application request path backed by an in-memory key-value store. Its Raft implementation remains an early prototype. Log replication, majority-based commits, durable storage, crash recovery, and end-to-end distributed consensus are not implemented.

Do not use TitanDB to store production or valuable data.»

---

Table of Contents

- "Overview" (#overview)
- "Current Features" (#current-features)
- "Project Status" (#project-status)
- "Architecture" (#architecture)
- "Core Components" (#core-components)
- "Binary Protocol" (#binary-protocol)
- "Requirements" (#requirements)
- "Installation" (#installation)
- "Quick Start" (#quick-start)
- "Repository Structure" (#repository-structure)
- "Storage Model and Data Lifecycle" (#storage-model-and-data-lifecycle)
- "Raft and Distributed Systems" (#raft-and-distributed-systems)
- "Testing and Quality" (#testing-and-quality)
- "Limitations and Security" (#limitations-and-security)
- "Development Roadmap" (#development-roadmap)
- "Non-Goals" (#non-goals)
- "Contributing" (#contributing)
- "License" (#license)
- "Author" (#author)

---

Overview

TitanDB is an experimental storage-engine project written in Go.

Its primary goal is to explore the engineering concepts behind storage systems and distributed infrastructure by implementing fundamental components directly.

Rather than relying on an existing database engine or high-level networking framework, TitanDB explores how a small storage service can be assembled from lower-level primitives.

The project focuses on:

- TCP networking and connection handling.
- Custom binary protocol design.
- Client-server communication.
- Concurrent programming in Go.
- In-memory key-value storage.
- Storage persistence and recovery concepts.
- Raft consensus and leader election.
- Replication and distributed state management.

TitanDB is primarily a systems-engineering and learning project. Its architecture is being developed incrementally, with correctness and implementation maturity taking priority over production-oriented claims.

Current Features

Implemented application functionality

- TCP-based server.
- Custom binary request framing.
- Command-line client.
- "PUT", "GET", and "DELETE" operations.
- In-memory key-value storage.
- Concurrent access protection using "sync.RWMutex".
- Basic request processing through the RPC package.

Experimental functionality

- "3SUM" command for finding unique integer triplets whose sum is zero.
- Early Raft node implementation.
- Node roles and election-related state.
- Election timers.
- RequestVote-related handling.
- Heartbeat-related functionality.

These experimental components should not be interpreted as a complete distributed consensus implementation.

Project Status

The following table distinguishes working application functionality from incomplete infrastructure components.

Component| Status| Description
TCP server| Implemented| Accepts client connections.
Binary request framing| Implemented| Custom application request format.
CLI client| Implemented| Sends commands to the server.
PUT / GET / DELETE| Implemented| Operate against local memory.
In-memory storage| Implemented| Go map protected by "sync.RWMutex".
3SUM| Experimental| Demonstration operation.
Raft node structure| Experimental| Early node and role implementation.
Leader election| In Progress| Election-related timers and vote handling.
Heartbeat handling| Experimental| Early heartbeat-related functionality.
Raft transport integration| Not Implemented| Normal server connections do not dispatch Raft messages.
AppendEntries replication| Not Implemented| No actual replicated log implementation.
Majority-based commit| Not Implemented| Writes are not committed through a cluster majority.
Persistent storage| Not Implemented| Data exists only in memory.
Write-ahead log| Not Implemented| No durable log.
Snapshots and recovery| Not Implemented| No persistent recovery mechanism.
Production security| Not Implemented| No authentication or transport encryption.
Multi-node fault tolerance| Not Verified| No demonstrated end-to-end fault-tolerance guarantees.

Current maturity: single-node TCP storage prototype with an experimental, disconnected Raft implementation.

---

Architecture

Current application request path

The current application architecture is a straightforward client-server storage service.

flowchart TD
    A["CLI Client"] -->|"TCP connection"| B["TCP Server"]
    B --> C["TitanRPC"]
    C --> D["Request Decoding"]
    D --> E["Command Processing"]
    E --> F["In-Memory Key-Value Store"]
    F --> G["Go map + sync.RWMutex"]

The application request path currently handles client operations independently of Raft.

A typical write follows this path:

1. The client constructs a binary request.
2. The request is transmitted through TCP.
3. The server accepts the connection.
4. The RPC layer decodes the request.
5. The corresponding operation is executed against the local key-value store.

Experimental Raft subsystem

The repository also contains an experimental Raft package.

It is separate from the normal application request path.

flowchart TD
    A["Experimental Raft Package"] --> B["Node State"]
    B --> C["Node Roles"]
    B --> D["Election Timers"]
    B --> E["RequestVote Handling"]
    B --> F["Heartbeat-Related Code"]
    F --> G["AppendEntries Handler"]
    G --> H["Incomplete Replication Logic"]

    I["Application RPC Server"] -.->|"Not integrated"| A

The Raft implementation is not currently responsible for committing application writes.

In particular:

- Application writes are applied directly to local memory.
- The RPC request handler does not use Raft's "SubmitCommand" method.
- The TCP server does not dispatch accepted application connections to "HandleRaftConnection".
- Application RPC and experimental Raft traffic currently use different message-header layouts.
- "AppendEntries" does not implement actual log replication.

Consequently, TitanDB does not currently provide distributed write consistency, replicated storage, or automatic failover guarantees.

Core Components

The Go implementation is organized around the following packages.

Package| Responsibility
"cmd/server"| Starts the TitanDB server and supporting components.
"cmd/client"| Provides the command-line client for communicating with the server.
"internal/rpc"| Handles TCP connections and application request encoding and decoding.
"internal/kv"| Implements the in-memory key-value store.
"internal/raft"| Contains experimental Raft node, election, and heartbeat-related code.
"pkg/logger"| Provides project logging utilities.

Application components

Server ("cmd/server")

Responsible for starting the server node and configuring its network-related parameters.

Client ("cmd/client")

Provides a command-line interface for sending storage commands to a running server.

RPC ("internal/rpc")

Handles application-level communication through TCP and the custom binary request format.

Key-value storage ("internal/kv")

Provides the current in-memory storage implementation, using a Go map and read/write synchronization.

Experimental components

Raft ("internal/raft")

Contains the early implementation of distributed consensus concepts, including node roles and election-related functionality.

It is not yet integrated into the application write path.

Logger ("pkg/logger")

Provides logging utilities used by the project.

---

Binary Protocol

TitanDB uses a custom binary application protocol over TCP.

The current application request header has a fixed size of 9 bytes.

Request header

Field| Size| Description
Payload Length| 4 bytes| Request payload length.
Message ID| 4 bytes| Request identifier.
Command| 1 byte| Operation identifier.
Payload| Variable| Command-specific data.

All multi-byte integer fields use big-endian byte order.

Command identifiers

Code| Command| Payload
"1"| PUT| Key length, key, value length, value.
"2"| GET| Key length and key.
"3"| DELETE| Key length and key.
"4"| 3SUM| No command-specific payload.

Payload structure

PUT

+------------+-----+--------------+-------+
| Key Length | Key | Value Length | Value |
+------------+-----+--------------+-------+

GET

+------------+-----+
| Key Length | Key |
+------------+-----+

DELETE

+------------+-----+
| Key Length | Key |
+------------+-----+

3SUM

The operation does not require additional command-specific payload.

Protocol limitations

The current protocol is experimental.

- No stable public wire-protocol contract exists.
- Compatibility between future protocol versions is not guaranteed.
- The application protocol and experimental Raft transport are not unified.
- The Raft implementation currently uses a separate 5-byte message header.
- A complete protocol-versioning and compatibility strategy has not been implemented.

The format described here documents the current application request layout, not a finalized protocol specification.

---

Requirements

- Go 1.22 or later.
- Git.
- A terminal with Go available in the system PATH.

Important: The Go module is currently located inside the "titan-db/" subdirectory.

Run Go commands from that directory rather than the repository root.

Installation

Clone the repository:

git clone https://github.com/kauand3v/titan-db.git
cd titan-db/titan-db

Verify the Go installation:

go version

Download module dependencies:

go mod download

Quick Start

The following example demonstrates the current single-node TCP application path.

It does not demonstrate distributed writes, replicated state, or Raft consensus.

1. Start the server

Open a terminal inside the Go module directory:

go run ./cmd/server/main.go \
  --id=1 \
  --address=127.0.0.1:8001 \
  --cluster=1@127.0.0.1:8001

Keep the server running.

The example uses:

- Node ID: "1"
- Server address: "127.0.0.1:8001"
- Cluster configuration: "1@127.0.0.1:8001"

The cluster argument is configuration input; it does not mean that replicated consensus is operational.

2. Insert a value

Open a second terminal in the same Go module directory.

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=PUT \
  --key=greeting \
  --val=hello

This sends a PUT request to the running server.

3. Read a value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=GET \
  --key=greeting

The request reads the value from the server's in-memory store.

4. Delete a value

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=DELETE \
  --key=greeting

This sends a DELETE request to the server.

5. Run the experimental 3SUM operation

go run ./cmd/client/main.go \
  --server=127.0.0.1:8001 \
  --cmd=3SUM

The experimental operation reads values from the in-memory store, parses integer values, and searches for unique triplets whose sum is zero.

It is a demonstration feature rather than part of the core key-value storage interface.

---

Repository Structure

Current repository layout:

titan-db/
├── .vscode/
├── ruby/
├── titan-db/
│   ├── cmd/
│   │   ├── client/
│   │   └── server/
│   ├── internal/
│   │   ├── kv/
│   │   ├── raft/
│   │   └── rpc/
│   ├── pkg/
│   │   └── logger/
│   └── go.mod
├── .gitignore
└── README.md

Directory responsibilities

Directory| Purpose
".vscode/"| Editor configuration.
"ruby/"| Experimental Ruby client material.
"titan-db/"| Current Go module.
"cmd/"| Executable applications.
"internal/"| Internal implementation packages.
"pkg/"| Reusable project utilities.

The repository currently keeps the Go module in a nested directory. Moving it to the repository root is a potential future structural improvement.

---

Storage Model and Data Lifecycle

TitanDB currently uses an in-memory key-value storage model.

The implementation relies on:

- A Go map for storing key-value pairs.
- "sync.RWMutex" for concurrent access protection.
- Direct local application of write operations.

Concurrency

The read/write mutex allows synchronized access to the in-memory map.

It provides local concurrency protection but does not establish distributed consistency or durability.

Data lifecycle

The current storage lifecycle is:

flowchart TD
    A["Client PUT Request"] --> B["RPC Handler"]
    B --> C["Local In-Memory Store"]
    C --> D["Value Available While Process Runs"]
    D --> E["Server Process Terminates"]
    E --> F["In-Memory Data Lost"]

Persistence status

Persistent storage is not implemented.

Currently:

- Data is stored only in process memory.
- There is no write-ahead log.
- There are no durable snapshots.
- There is no recovery mechanism.
- Restarting the server loses previously stored values.
- The storage "Close" operation does not provide persistence.

A future persistent storage implementation will require a defined durability model, recovery behavior, and appropriate crash-consistency tests.

---

Raft and Distributed Systems

TitanDB includes an experimental Raft package intended to explore distributed consensus.

Raft is a consensus algorithm that coordinates replicated logs across multiple nodes.

A complete implementation requires substantially more than leader election.

Current implementation state

The existing Raft package contains early infrastructure for:

- Node roles.
- Election-related state.
- Election timers.
- RequestVote-related handling.
- Heartbeat-related functionality.
- Early AppendEntries handling.

These components are not sufficient to establish a functioning consensus system.

Missing consensus functionality

The following capabilities remain unimplemented:

Capability| Current state
Complete leader election correctness| Not verified
Term and log-freshness validation| Incomplete / not verified
Raft transport integration| Not implemented
Replicated log entries| Not implemented
Previous-log consistency validation| Not implemented
Majority-based commit| Not implemented
State machine application of committed entries| Not implemented
Durable term and vote state| Not implemented
Write-ahead log| Not implemented
Snapshotting| Not implemented
Crash recovery| Not implemented
Multi-node failure validation| Not verified

Why the current implementation is not consensus

There are several architectural gaps.

1. Separate transport paths

The application RPC server uses a 9-byte request header, while the experimental Raft implementation uses a separate 5-byte message header.

The normal accepted-connection handler does not dispatch connections to "HandleRaftConnection".

2. Application writes bypass Raft

The current application request handler applies writes directly to the local key-value store.

It does not submit writes through the Raft command path or wait for a majority commit.

3. AppendEntries is incomplete

The existing handler does not implement actual log-entry validation and replication. Updating node state and returning a response is not equivalent to replicating a Raft log.

4. No durable consensus state

There is no demonstrated persistent storage for Raft terms, votes, logs, or committed application data.

5. No demonstrated recovery guarantees

There are no verified multi-node failure scenarios establishing recovery, consistency, or fault tolerance.

Development direction

The intended direction is to evolve the experimental subsystem into a properly integrated consensus layer.

That requires transport integration, correct Raft protocol behavior, replicated log management, majority-based commitment, durable state, and multi-node testing.

Until those capabilities exist and are validated, TitanDB must be treated as a single-node storage prototype with experimental distributed-systems code.

---

Testing and Quality

Testing distributed infrastructure requires more than confirming that the project compiles.

The following Go commands are useful validation checks.

Run them from the Go module directory.

Run tests

go test ./...

Run the race detector

go test -race ./...

Run static analysis

go vet ./...

Important testing considerations

These commands are validation procedures, not evidence of existing test coverage.

A successful compilation does not establish:

- Correct TCP behavior under concurrent connections.
- Correct binary protocol handling.
- Absence of all concurrency bugs.
- Correct Raft elections.
- Correct log replication.
- Majority-based commitment.
- Linearizability.
- Crash recovery.
- Fault tolerance.

Distributed-systems correctness requires dedicated multi-process integration tests.

In particular, Raft testing should include leader elections, competing candidates, network failures, stale terms, log inconsistencies, leader replacement, and recovery after process termination.

No benchmark or test-coverage figures are claimed by this project.

---

Limitations and Security

TitanDB is an experimental implementation and currently has significant limitations.

Storage

- In-memory-only storage.
- No persistent data files.
- No write-ahead log.
- No snapshots.
- No crash recovery.
- No durable write guarantees.

Networking

- Experimental binary protocol.
- No stable compatibility contract.
- No demonstrated protocol version negotiation.
- No verified production-grade connection management.
- No demonstrated resource and connection limits.

Distributed systems

- Raft is not integrated into application writes.
- No actual replicated log.
- No majority-based commit.
- No verified multi-node failover.
- No demonstrated distributed consistency guarantees.

Security

The current implementation does not provide verified production security mechanisms, including:

- Client authentication.
- Authorization.
- Transport encryption.
- Secure credential management.
- Production-grade resource isolation.
- Operational security hardening.

TitanDB should only be exposed in controlled development and testing environments.

---

Development Roadmap

The roadmap prioritizes correctness, reliability, and architectural foundations before performance optimization.

1. Repository and module structure

- [ ] Move the Go module to the repository root.
- [ ] Update the module path to "github.com/kauand3v/titan-db".
- [ ] Organize optional language clients under a dedicated directory.
- [ ] Standardize project configuration and development instructions.

2. Testing foundations

- [ ] Add unit tests for key-value operations.
- [ ] Add binary protocol encoding and decoding tests.
- [ ] Add TCP server integration tests.
- [ ] Add client-server communication tests.
- [ ] Add concurrency and race-condition tests.

3. Protocol correctness

- [ ] Define a consistent message framing strategy.
- [ ] Validate payload lengths and malformed requests.
- [ ] Define protocol error handling.
- [ ] Introduce protocol versioning.
- [ ] Document request and response behavior.

4. TCP reliability

- [ ] Improve connection lifecycle management.
- [ ] Add configurable timeouts.
- [ ] Handle partial reads and writes correctly.
- [ ] Introduce resource and connection limits.
- [ ] Add network failure tests.

5. Raft transport integration

- [ ] Integrate Raft message handling with the networking layer.
- [ ] Define explicit transport handling for Raft messages.
- [ ] Establish a consistent node communication model.
- [ ] Test communication between multiple server processes.

6. Consensus correctness

- [ ] Complete RequestVote term validation.
- [ ] Implement log-freshness checks.
- [ ] Implement correct AppendEntries handling.
- [ ] Add previous-log consistency validation.
- [ ] Implement actual log replication.
- [ ] Implement majority-based commit.
- [ ] Apply committed entries to the state machine.

7. Persistence and recovery

- [ ] Implement a write-ahead log.
- [ ] Persist Raft term and vote state.
- [ ] Implement durable log storage.
- [ ] Implement snapshots.
- [ ] Implement crash recovery.
- [ ] Validate recovery after unexpected termination.

8. Distributed validation

- [ ] Add multi-node integration tests.
- [ ] Test leader election.
- [ ] Test leader failure and replacement.
- [ ] Test network partitions.
- [ ] Test node restarts.
- [ ] Verify replicated state consistency.

9. Performance and observability

- [ ] Add storage benchmarks.
- [ ] Benchmark request throughput and latency.
- [ ] Measure memory and resource usage.
- [ ] Improve structured logging.
- [ ] Add operational metrics.

10. Security and hardening

- [ ] Implement client authentication.
- [ ] Introduce authorization policies.
- [ ] Add transport encryption.
- [ ] Define secure configuration practices.
- [ ] Improve operational safeguards.
- [ ] Perform security-oriented testing.

---

Non-Goals

In its current state, TitanDB does not aim to:

- Replace Redis, etcd, or established database systems.
- Provide production-grade availability.
- Guarantee durable writes.
- Provide a stable public protocol or SDK.
- Guarantee linearizability.
- Guarantee fault tolerance.
- Serve as a production distributed database.

The project's current objective is to explore and implement the engineering foundations required by storage systems and distributed infrastructure.

---

Contributing

Contributions, bug reports, and technical discussions are welcome.

Before implementing a significant architectural change, consider opening an issue describing:

1. The problem being addressed.
2. The expected behavior.
3. The proposed design.
4. Relevant technical trade-offs.
5. The tests required to validate the implementation.

Pull requests

Keep pull requests focused and include appropriate tests whenever possible.

Changes involving networking, protocol encoding, storage, concurrency, or distributed-systems behavior should include reproducible validation.

Please avoid introducing undocumented functionality or making claims that cannot be demonstrated by the implementation.

---

License

No "LICENSE" file is currently present in the repository.

Until a license is declared, no permission to reuse, modify, or redistribute the project should be assumed.

---

Author

Kauan Oliveira

GitHub: "@kauand3v" (https://github.com/kauand3v)

---

<p align="center">
  <sub>TitanDB — An experimental journey into storage systems and distributed systems engineering.</sub>
</p>