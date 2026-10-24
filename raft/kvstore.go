// Project 2: Resilient Raft -- KV state machine (given, complete).
// This is deliberately simple -- the interesting engineering in this
// project is getting commands to this state machine correctly and in
// the right order (that's node.go's job), not the state machine
// itself.
package raft

import "sync"

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{data: make(map[string]string)}
}

// Apply is called ONLY with committed, in-order commands (node.go's
// job to guarantee that).
func (kv *KVStore) Apply(cmd KVCommand) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	switch cmd.Op {
	case "PUT":
		kv.data[cmd.Key] = cmd.Value
	case "DELETE":
		delete(kv.data, cmd.Key)
	}
}

func (kv *KVStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, ok := kv.data[key]
	return v, ok
}
