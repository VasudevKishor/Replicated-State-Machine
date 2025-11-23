# Distributed Key-Value Store (Raft Consensus)

## 1. Introduction
This project is a Distributed Key-Value Store implemented in Go. It achieves fault tolerance and strong consistency by replicating a state machine across a cluster of servers using the **Raft Consensus Algorithm**.

The system is designed to be resilient: as long as a majority of nodes (N/2 + 1) are operational, the cluster remains available to serve reads and writes.

## 2. System Architecture
The system follows the standard "Replicated State Machine" architecture defined in the Raft Extended Paper.

### High-Level Components
1.  **Consensus Module (Raft):** The core engine. It manages the leader election, log replication, and safety properties.
2.  **Write-Ahead Log (WAL):** A persistent append-only file that records every command accepted by the Leader.
3.  **State Machine (KV Store):** An in-memory map that applies committed commands from the log to update the current state.
4.  **RPC Layer:** Handles communication between nodes (RequestVote, AppendEntries).

![Raft Architecture Diagram](docs/architecture.png)
*(Note: Use a standard Raft architecture diagram here)*

## 3. Functional Requirements

### 3.1 Leader Election
* **Role Transitions:** Nodes must transition correctly between Follower, Candidate, and Leader states based on timeouts.
* **Election Safety:** A term can have at most one leader.
* **Heartbeats:** The Leader must send periodic heartbeats to suppress new elections.
* **Randomized Timeouts:** Election timeouts must be randomized (e.g., 150-300ms) to prevent split votes.

### 3.2 Log Replication
* **Append-Only:** The Leader appends client commands to its log and replicates them to Followers in parallel.
* **Consistency Check:** Followers must reject `AppendEntries` requests if their previous log entry does not match the Leader's (Log Matching Property).
* **Commit Index:** Once a command is replicated to a majority of servers, it is considered "Committed" and applied to the State Machine.

### 3.3 Persistence
To survive crash-restarts, each server must persist the following metadata to disk before responding to RPCs:
* `currentTerm` (The latest term server has seen)
* `votedFor` (CandidateId that received vote in current term)
* `log[]` (Log entries)

### 3.4 Client API
* **`GET(key)`**: Returns the value for a key.
* **`PUT(key, value)`**: Sets the value for a key. Returns an error or redirects if the node is not the Leader.

## 4. Technical Constraints
* **Language:** Go (Golang) 1.20+
* **Communication Protocol:** Go `net/rpc` over TCP.
* **Concurrency Model:**
    * Use `sync.Mutex` for shared state protection.
    * Use `channels` for inter-module signaling (e.g., commit notifications).
    * Use `time.Ticker` for heartbeats and election timers.
* **No External Dependencies:** The core consensus logic must be written from scratch (no `hashicorp/raft`).

## 5. RPC Specifications
The nodes communicate using strictly these two RPC methods.

### 5.1 `RequestVote`
Invoked by Candidates to gather votes.

**Arguments:**
| Field | Type | Description |
| :--- | :--- | :--- |
| `Term` | `int` | Candidate's term |
| `CandidateId` | `int` | Candidate requesting vote |
| `LastLogIndex` | `int` | Index of candidate’s last log entry |
| `LastLogTerm` | `int` | Term of candidate’s last log entry |

**Results:**
| Field | Type | Description |
| :--- | :--- | :--- |
| `Term` | `int` | CurrentTerm, for candidate to update itself |
| `VoteGranted` | `bool` | True means candidate received vote |

### 5.2 `AppendEntries`
Invoked by Leader to replicate log entries; also used as heartbeat.

**Arguments:**
| Field | Type | Description |
| :--- | :--- | :--- |
| `Term` | `int` | Leader’s term |
| `LeaderId` | `int` | So follower can redirect clients |
| `PrevLogIndex` | `int` | Index of log entry immediately preceding new ones |
| `PrevLogTerm` | `int` | Term of prevLogIndex entry |
| `Entries` | `[]Entry` | Log entries to store (empty for heartbeat) |
| `LeaderCommit` | `int` | Leader’s commitIndex |

**Results:**
| Field | Type | Description |
| :--- | :--- | :--- |
| `Term` | `int` | CurrentTerm, for leader to update itself |
| `Success` | `bool` | True if follower contained entry matching PrevLogIndex |

## 6. Directory Structure
```text
.
├── go.mod
├── main.go             # Entry point
├── raft/
│   ├── raft.go         # Core consensus logic
│   ├── rpc.go          # RPC struct definitions
│   └── util.go         # Debugging and helpers
└── kv/
    ├── server.go       # KV Store logic
    └── client.go       # Client CLI
