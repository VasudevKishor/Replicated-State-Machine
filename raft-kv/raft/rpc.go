package raft

// LogEntry represents a single command (e.g., "SET x=5")
type LogEntry struct {
	Term    int         // The term when this entry was received
	Command interface{} // The actual command (payload)
}

// --- RequestVote RPC ---

// RequestVoteArgs is sent by a Candidate to ask for a vote.
type RequestVoteArgs struct {
	Term         int // Candidate's Term
	CandidateId  int // ID of the candidate requesting vote
	LastLogIndex int // Index of candidate’s last log entry
	LastLogTerm  int // Term of candidate’s last log entry
}

// RequestVoteReply is the response from a Follower.
type RequestVoteReply struct {
	Term        int  // CurrentTerm, for candidate to update itself
	VoteGranted bool // True means candidate received vote
}

// --- AppendEntries RPC ---

// AppendEntriesArgs is sent by the Leader to replicate log entries.
// Also acts as a Heartbeat if Entries is empty.
type AppendEntriesArgs struct {
	Term         int        // Leader’s term
	LeaderId     int        // So follower can redirect clients
	PrevLogIndex int        // Index of log entry immediately preceding new ones
	PrevLogTerm  int        // Term of prevLogIndex entry
	Entries      []LogEntry // Log entries to store (empty for heartbeat)
	LeaderCommit int        // Leader’s commitIndex
}

// AppendEntriesReply is the response from a Follower.
type AppendEntriesReply struct {
	Term    int  // CurrentTerm, for leader to update itself
	Success bool // True if follower contained entry matching PrevLogIndex
}