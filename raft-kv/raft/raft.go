package raft

import (
	"sync"
	"time"
)

// State represents the role of the node
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type Raft struct {
	mu    sync.Mutex // Lock to protect shared access to this peer's state
	peers []string   // RPC addresses of all peers (e.g., "localhost:8001")
	me    int        // This peer's index into peers[]

	// Persistent state on all servers
	currentTerm int
	votedFor    int // -1 if null (nobody voted for yet)
	log         []LogEntry

	// Volatile state on all servers
	commitIndex int // Index of highest log entry known to be committed
	lastApplied int // Index of highest log entry applied to state machine

	// Volatile state on leaders (Reinitialized after election)
	nextIndex  []int // For each server, index of the next log entry to send
	matchIndex []int // For each server, index of highest log entry known to be replicated

	// State Management
	state          State
	electionTimer  *time.Timer
	heartbeatTimer *time.Ticker

	// ApplyChannel: Sending committed messages to the "Client" (KV Store)
	applyCh chan ApplyMsg
}

// ApplyMsg is sent to the KV Store when a command is committed
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// NewRaft creates a new Raft server instance
func NewRaft(peers []string, me int, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.me = me
	rf.applyCh = applyCh

	// Initialize state
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.log = make([]LogEntry, 0)

	// Add a dummy entry at index 0 to make math easier (1-based indexing)
	rf.log = append(rf.log, LogEntry{Term: 0}) 

	rf.commitIndex = 0
	rf.lastApplied = 0

	// Initialize timers (we will start them in the next step)
	// We set a random election timeout later
	rf.electionTimer = time.NewTimer(time.Second) 
	rf.heartbeatTimer = time.NewTicker(100 * time.Millisecond)

	return rf
}