// Project 2: Resilient Raft -- file-based Persister (given, complete).
// Simple whole-state-overwrite JSON file. Not the most efficient
// possible design (a real Raft implementation appends incrementally),
// but simple, correct, and enough to prove the "must persist before
// replying" safety property -- see node.go's Persister interface.
package raft

import (
	"encoding/json"
	"os"
)

type FilePersister struct {
	path string
}

func NewFilePersister(path string) *FilePersister {
	return &FilePersister{path: path}
}

type persistedState struct {
	CurrentTerm int64
	VotedFor    int32
	Log         []LogEntry
}

func (p *FilePersister) SaveState(currentTerm int64, votedFor int32, log []LogEntry) error {
	state := persistedState{CurrentTerm: currentTerm, VotedFor: votedFor, Log: log}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	tmpPath := p.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}
	// Atomic rename so a crash mid-write never leaves a corrupt state file.
	return os.Rename(tmpPath, p.path)
}

func (p *FilePersister) LoadState() (int64, int32, []LogEntry, error) {
	data, err := os.ReadFile(p.path)
	if os.IsNotExist(err) {
		return 0, -1, nil, nil // fresh node, no prior state
	}
	if err != nil {
		return 0, -1, nil, err
	}
	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return 0, -1, nil, err
	}
	return state.CurrentTerm, state.VotedFor, state.Log, nil
}
