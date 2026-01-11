package main

import (
	"flag"
	"log"

	"github.com/prism/daemon/internal/agent"
)

func main() {
	agentID := flag.String("agent-id", "", "The UUID of the agent this daemon represents")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis address")
	grpcAddr := flag.String("grpc-addr", "localhost:9090", "Prism Server gRPC address")

	flag.Parse()

	if *agentID == "" {
		log.Fatal("Please provide --agent-id")
	}

	cfg := agent.Config{
		AgentID:   *agentID,
		RedisAddr: *redisAddr,
		GrpcAddr:  *grpcAddr,
	}

	w := agent.NewWorker(cfg)
	w.Start()
}
