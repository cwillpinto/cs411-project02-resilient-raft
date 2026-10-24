package raft

import (
	"path/filepath"
	"testing"

	pb "resilient_raft/proto"
)

// This first check defines the persistence and voting contract. Students add
// five-node election, replication, and fault-injection tests to this package.
func TestRequestVotePersistsBeforeReply(t *testing.T) {
	path := filepath.Join(t.TempDir(), "node.json")
	persister := NewFilePersister(path)
	node := NewNode(1, []int32{2, 3, 4, 5}, persister, NewKVStore())
	reply := node.HandleRequestVote(&pb.RequestVoteArgs{
		Term: 1, CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0,
	})
	if !reply.GetVoteGranted() || reply.GetTerm() != 1 {
		t.Fatalf("expected vote for candidate 2 in term 1; got %+v", reply)
	}
	term, votedFor, _, err := persister.LoadState()
	if err != nil {
		t.Fatal(err)
	}
	if term != 1 || votedFor != 2 {
		t.Fatalf("vote reply preceded durable state: term=%d votedFor=%d", term, votedFor)
	}
}
