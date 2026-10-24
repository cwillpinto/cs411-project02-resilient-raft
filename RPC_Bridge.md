# RPC Bridge for Project 2

Week 7's RPC lesson was not held. Complete this bridge before building the Raft network layer. It uses the same client, server, protobuf message, and timeout ideas as Project 2 without implementing any Raft requirement for you.

**Instructor starter:** https://github.com/cwillpinto/cs411-project02-resilient-raft

1. From the repository root, run `go run ./cmd/rpc_smoke`. Expect `RPC_BRIDGE_READY`.
2. In `cmd/rpc_smoke/main.go`, identify the listening address, server registration, client stub, request message, response, and deadline.
3. Change the request's service name to `missing`. Predict and observe the response. Restore it afterward.
4. Reduce the context timeout to a tiny value. Observe that a deadline failure may depend on timing; explain why every outbound Raft RPC needs a bounded timeout.
5. Read `proto/raft.proto`. Match `RequestVoteArgs` and `RequestVoteReply` to the health request and response roles. Generate the Raft stubs with the command in the README; `go build ./...` then confirms the generated code compiles.

The bridge proves only that one local gRPC request and response can complete. It does not prove election safety, log replication, partition handling, or linearizable reads. Implement and test those in the five-node project.
