// Project 2: Resilient Raft -- node state machine scaffold.
//
// This is the Project 2 starter due November 4, 2026 -- this file is an
// architectural scaffold covering Phases 1-3's data model, not a
// finished implementation.
package raft

import (
	"sync"
	"time"

	pb "resilient_raft/proto"
)

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

// KVCommand is what gets replicated through the log and applied to
// the state machine once committed.
type KVCommand struct {
	Op    string // "PUT", "DELETE"
	Key   string
	Value string
}

type LogEntry struct {
	Term    int64
	Command KVCommand
}

// Persister is given: a minimal append-only-file persistence layer.
// Persist currentTerm/votedFor/log to disk BEFORE responding to any
// RPC that depends on them (Pro-Tip in the assignment) -- a node that
// "forgets" its vote after a simulated crash-and-restart is a safety
// bug, not a performance detail.
type Persister interface {
	SaveState(currentTerm int64, votedFor int32, log []LogEntry) error
	LoadState() (currentTerm int64, votedFor int32, log []LogEntry, err error)
}

type Node struct {
	mu sync.Mutex

	id      int32
	peerIDs []int32
	peers   map[int32]pb.RaftInternalClient // gRPC client stubs to each peer

	persister Persister

	// Persistent state (must be saved via persister before replying to RPCs)
	currentTerm int64
	votedFor    int32 // -1 if none
	log         []LogEntry

	// Volatile state
	role        Role
	commitIndex int64
	lastApplied int64

	// Volatile leader-only state (reset on election)
	nextIndex  map[int32]int64
	matchIndex map[int32]int64

	kv            *KVStore
	electionTimer *time.Timer
	stopped       bool
}

func NewNode(id int32, peerIDs []int32, persister Persister, kv *KVStore) *Node {
	n := &Node{
		id:        id,
		peerIDs:   peerIDs,
		persister: persister,
		votedFor:  -1,
		role:      Follower,
		kv:        kv,
	}
	// TODO (Phase 1/2): on startup, call persister.LoadState() to
	// recover currentTerm/votedFor/log from a prior crash, before
	// starting the election timer.
	return n
}

// --- Phase 2: Leader Election ---

func (n *Node) resetElectionTimer() {
	// TODO: randomized [150ms, 300ms) jitter, same pattern as Lab 8 --
	// but this time it's driving REAL gRPC RequestVote calls, not
	// in-process function calls.
}

func (n *Node) startElection() {
	// TODO: increment currentTerm, vote for self, persist state BEFORE
	// sending any RPCs, then fan out RequestVote RPCs to every peer
	// concurrently (with a context.WithTimeout on each -- Pro-Tip),
	// tally votes, become Leader on a majority.
}

// HandleRequestVote is the RequestVote RPC handler (called by the
// gRPC service implementation, not directly by peers).
func (n *Node) HandleRequestVote(args *pb.RequestVoteArgs) *pb.RequestVoteReply {
	// TODO (Phase 2): term check, vote-granting rule including the log
	// up-to-date check (same rule as Lab 8, Part A Q4) -- persist
	// currentTerm/votedFor BEFORE returning if either changed.
	return &pb.RequestVoteReply{}
}

// --- Phase 3: Log Replication & State Machine ---

// HandleAppendEntries is the AppendEntries RPC handler (heartbeats
// AND log replication both go through this).
func (n *Node) HandleAppendEntries(args *pb.AppendEntriesArgs) *pb.AppendEntriesReply {
	// TODO (Phase 3):
	// 1. Term check (reject stale terms; step down if args.Term > currentTerm).
	// 2. Log-matching consistency check at PrevLogIndex/PrevLogTerm --
	//    reject if the follower's log doesn't have a matching entry there.
	// 3. On conflict, truncate the follower's log from the conflict
	//    point and append the leader's entries.
	// 4. Advance commitIndex to min(LeaderCommit, index of last new
	//    entry), then apply newly-committed entries to n.kv in order.
	// 5. Reset the election timer (this IS the heartbeat suppression
	//    mechanism).
	return &pb.AppendEntriesReply{}
}

// replicateToPeer sends AppendEntries to one follower, retrying with
// backtracked nextIndex on log-mismatch rejection. Only the leader
// calls this (one goroutine per follower is the recommended shape).
func (n *Node) replicateToPeer(peerID int32) {
	// TODO (Phase 3, Phase 4): also responsible for advancing
	// commitIndex once an entry is replicated to a majority (n.mu-protected).
}

// --- Client-facing KV ops (Phase 3) ---

// Propose appends a new command to the leader's log (if this node IS
// the leader) and blocks until it's committed and applied, or returns
// a not-leader indication.
func (n *Node) Propose(cmd KVCommand) (applied bool, isLeader bool, leaderHint int32) {
	// TODO: if role != Leader, return (false, false, <last known leader, if any>).
	// Otherwise append to log, persist, trigger replication, wait for
	// commit (e.g. via a condition variable or a per-index channel),
	// then return (true, true, n.id).
	return false, false, -1
}

