package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/prism/daemon/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AgentID     string
	RedisAddr   string
	GrpcAddr    string
}

type Worker struct {
	cfg        Config
	redisClient *redis.Client
	grpcClient  pb.AgentServiceClient
}

func NewWorker(cfg Config) *Worker {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	conn, err := grpc.NewClient(cfg.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewAgentServiceClient(conn)

	return &Worker{
		cfg:        cfg,
		redisClient: rdb,
		grpcClient:  c,
	}
}

func (w *Worker) Start() {
	log.Printf("Starting Agent Worker for ID: %s", w.cfg.AgentID)

	// Register first
	_, err := w.grpcClient.RegisterAgent(context.Background(), &pb.RegisterAgentRequest{
		AgentId: w.cfg.AgentID,
		Version: "0.0.1",
	})
	if err != nil {
		log.Printf("Failed to register agent: %v. Continuing anyway...", err)
	}

	// Subscribe to Redis
	channel := fmt.Sprintf("tasks:%s", w.cfg.AgentID)
	pubsub := w.redisClient.Subscribe(context.Background(), channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Received task: %s", msg.Payload)
		// Payload format: "taskId:title"
		// In real world, this would be JSON
		go w.processTask(msg.Payload)
	}
}

func (w *Worker) processTask(payload string) {
	// Parse simplified payload
	var taskId string
	// Split by first colon
	if n, err := fmt.Sscanf(payload, "%s", &taskId); err != nil || n == 0 {
		taskId = payload // fallback
	}
	// Clean up taskId (remove title part if simple scan didn't work well)
	// For MVP, let's assume payload is just ID or "ID:Title" and we take the first part
	for i, c := range payload {
		if c == ':' {
			taskId = payload[:i]
			break
		}
	}

	log.Printf("Processing Task %s...", taskId)
	w.updateStatus(taskId, "IN_PROGRESS", "Agent started working on task...")

	// Simulate work
	time.Sleep(5 * time.Second)

	w.updateStatus(taskId, "DONE", "Work completed successfully by local daemon.")
	log.Printf("Task %s completed.", taskId)
}

func (w *Worker) updateStatus(taskId, status, details string) {
	_, err := w.grpcClient.UpdateTaskStatus(context.Background(), &pb.UpdateTaskStatusRequest{
		TaskId:  taskId,
		AgentId: w.cfg.AgentID,
		Status:  status,
		Details: details,
	})
	if err != nil {
		log.Printf("Failed to update status: %v", err)
	}
}
