// A small RPC bridge for students who missed the Week 7 RPC lesson.
// This health request is independent of the Raft implementation.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("rpc-bridge", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()

	connection, err := grpc.NewClient(listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	reply, err := healthpb.NewHealthClient(connection).Check(ctx,
		&healthpb.HealthCheckRequest{Service: "rpc-bridge"})
	if err != nil {
		log.Fatal(err)
	}
	if reply.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		log.Fatalf("unexpected status: %s", reply.GetStatus())
	}
	fmt.Println("RPC_BRIDGE_READY")
}
