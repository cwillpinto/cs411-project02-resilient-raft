# ☁️ CS 411 Project 2: Resilient Raft

**Instructor starter repository:** https://github.com/cwillpinto/cs411-project02-resilient-raft
**Points:** 150 total, including the 15-point Distributed Phase Gate.
**Due:** Wednesday, November 4, 2026 at 11:59 PM Central.
**Objective:** Design and implement a Raft-replicated, in-memory key-value store in Go, communicating over gRPC, that survives node crashes and network partitions while preserving linearizable reads and writes — a real, unabridged consensus node, not the guided simulation you wrangled in Lab 8.

---

## RPC bridge before the Raft work

The Week 7 RPC class was missed. Start with [RPC_Bridge.md](https://github.com/cwillpinto/cs411-project02-resilient-raft/blob/main/RPC_Bridge.md) in the starter repository and run `go run ./cmd/rpc_smoke`. The bridge demonstrates a local gRPC request, response, and timeout. It prepares you for Phase 1; it does not replace the five-node cluster.

## 🌐 The Technical Challenge
A single database server is a single point of failure. Real distributed datastores (etcd, CockroachDB, Consul) replicate every write across a cluster of nodes and use a consensus protocol to guarantee that, even when machines crash or the network partitions, every node that *does* respond agrees on the same sequence of committed operations — never two conflicting histories.

You will build a 5-node cluster where each node runs the full Raft protocol: leader election with randomized timeouts, log replication via `AppendEntries`, and a replicated state machine exposing a simple `Get`/`Put`/`Delete` key-value API. Your cluster must keep serving client requests correctly as long as a majority (3 of 5) of nodes remain reachable, and must refuse to silently lose or fork committed writes when the minority is crashed, partitioned, or lagging.

---

## 🛠️ Detailed Student Tasks

1. **Phase 1 — gRPC Scaffold & Protocol Definitions.** Define a `.proto` schema for the internal Raft RPCs (`RequestVote`, `AppendEntries`) and the external client API (`Get`, `Put`, `Delete`). Generate Go stubs and stand up 5 independent node processes (or goroutine-isolated instances) that can dial each other over gRPC on distinct ports.
2. **Phase 2 — Leader Election.** Implement the Follower/Candidate/Leader state machine: randomized election timeouts (150–300ms), term increment and persistence on candidacy, vote-granting rules (including the log up-to-date check), and heartbeat-driven suppression of new elections while a leader is alive. A killed leader must trigger a new election among the remaining nodes within one timeout window.
3. **Phase 3 — Log Replication & State Machine.** Implement `AppendEntries` for both heartbeats and log replication, including the log-matching consistency check (`prevLogIndex`/`prevLogTerm`), conflict detection and truncation of divergent follower logs, and commit-index advancement once an entry is replicated to a majority. Apply committed entries, in order, to each node's local key-value state machine.
4. **Phase 4 — Fault Injection & Cluster Verification.** Build a test harness that runs a local 5-node cluster and injects: (a) leader crash mid-term, (b) a minority-partition (2 nodes cut off from the other 3), and (c) a crashed node rejoining with a stale log. Write `go test ./raft` covering election safety (never two leaders in the same term), log matching, and leader completeness (no committed entry is ever lost). Verify client `Put`/`Get` correctness continues to hold throughout every injected failure.

---

## 💡 Pro-Tips & Implementation Hints
- Model each node's state transitions with a single goroutine reading from channels (timer fires, RPC arrivals, client requests) rather than protecting shared fields with ad hoc mutexes scattered across handlers — it's easier to reason about and easier to keep race-free.
- Persist `currentTerm`, `votedFor`, and the log to local disk (a simple append-only file is fine) before responding to any RPC that depends on them — a node that "forgets" its vote after a simulated crash-and-restart is a safety bug, not a performance detail.
- Use `context.WithTimeout` on every outbound RPC. A follower that never times out on a dead peer will stall your entire replication loop.
- Run everything under `go test -race ./raft` continuously during development, not just before submission. A race in the term/state fields will produce a cluster that looks correct 95% of the time and forks silently the other 5%.
- Compare a partitioned-minority leader (still thinks it's leader, can't reach a majority) against the majority side electing a new leader — this is the classic split-brain scenario your fault injection tests must explicitly cover.

---

## ⚖️ Grading Rubric
| Category | Weight | Points | Description |
| :--- | :---: | ---: | :--- |
| Leader Election Correctness | 20% | 30 | Exactly one leader per term is ever elected; randomized timeouts prevent livelock; crashed leaders are replaced within one timeout window. |
| Log Replication & Consistency | 25% | 37.5 | `AppendEntries` correctly detects and resolves log conflicts; commit index only advances on majority replication; no committed entry is ever lost or reordered. |
| Fault Tolerance Under Induced Failures | 25% | 37.5 | Cluster continues serving correct `Get`/`Put` results through leader crash, minority partition, and stale-node rejoin scenarios exercised by the test harness. |
| gRPC API & KV Store Functionality | 15% | 22.5 | `Get`/`Put`/`Delete` correctly routed to the current leader, correctly applied only after commit, correctly rejected/redirected when sent to a non-leader. |
| Code Quality, Testing & Documentation | 15% | 22.5 | Race-free (`go test -race`), meaningful test coverage of the safety properties above, and a README documenting how to launch and stress-test the cluster. |

The Distributed Phase Gate is **15 points within** the 37.5-point fault-tolerance criterion, not a separate assignment. Submit one reproducible fault-injection command, a timestamped term/leader/log or commit trace, and a short safety or recovery claim that the trace supports. The other 22.5 fault-tolerance points cover the required crash, minority-partition, and stale-rejoin behavior.

---

## Project Timeline

The full Project 2 cluster, tests, README, and Distributed Phase Gate evidence are due **Wednesday, November 4, 2026 at 11:59 PM Central**, under the approved schedule amendment.

---

**Terminal-First Warning:** This project is graded by an automated harness that builds and stress-tests your cluster entirely from the command line (`go build ./...`, `go test ./raft`, and a scripted multi-process cluster launch with injected `SIGKILL`s and simulated network drops). If your implementation only behaves correctly inside your IDE's debugger, or depends on manual steps not captured in your README, it will be treated as non-functional. Verify your submission from a clean terminal before you turn it in.
