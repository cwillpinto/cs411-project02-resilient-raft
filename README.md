# CS411 Project 2 Resilient Raft

**Instructor starter:** https://github.com/cwillpinto/cs411-project02-resilient-raft  
**Due:** Wednesday, November 4, 2026 at 11:59 PM Central
**Points:** 150 total, including the 15-point Distributed Phase Gate  
**Task:** [Project_2_Resilient_Raft.md](Project_2_Resilient_Raft.md)

This Go scaffold is intentionally incomplete. It includes protocol definitions, node TODOs, a KV store, and persistence helpers. It does not implement a working five-node cluster or fault-injection harness. Students build and verify those components as part of the project.

## Work in your own repository

Choose **Use this template** on the instructor repository and create a repository under your own GitHub account. Clone your repository, work there, and push to your `origin`. If you clone the instructor starter directly, create an empty repository in your account, then run:

```sh
git remote rename origin upstream
git remote add origin <your-repository-url>
git remote -v
git push -u origin <your-branch>
```

## RPC bridge

Week 7 RPC was missed. Read [RPC_Bridge.md](RPC_Bridge.md), then run `go run ./cmd/rpc_smoke` and expect `RPC_BRIDGE_READY` before building the Raft network layer.

## Generate and check

Install Go 1.25, `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc`. From the repository root:

```sh
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/raft.proto
go mod tidy
go build ./...
go test ./raft
go test -race ./raft
```

A scaffold build is only a setup check. Your own tests must cover leader election, log matching, leader completeness, and the specified crash, partition, and stale-rejoin cases in a five-node cluster. Document the exact launch and fault-injection commands.

Submit your student-owned repository URL and commit SHA, launch instructions, test and race-detector output, fault commands and traces, and the Distributed Phase Gate command, trace, and supported claim through Blackboard.

