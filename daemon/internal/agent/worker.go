package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prism/daemon/internal/git"
	"github.com/prism/daemon/internal/github"
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
	gitService  git.GitService
	ghClient    github.GitHubClient
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
		gitService:  git.NewGitService(),
		ghClient:    github.NewMockGitHubClient(),
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
	w.updateStatus(taskId, "IN_PROGRESS", "Agent started working on task...", "", "", "")

	// 1. Init/Open Git Repo (Simulate in /tmp)
	repoPath := filepath.Join(os.TempDir(), "prism-repo")
	_ = os.MkdirAll(repoPath, 0755)
	repo, err := w.gitService.InitOrOpen(repoPath)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("Failed to init git repo: %v", err), "", "", "")
		return
	}

	// 2. Create Branch
	branchName := fmt.Sprintf("feat/task-%s", taskId)
	if err := w.gitService.CreateBranch(repo, branchName); err != nil {
		log.Printf("Branch creation failed (might exist): %v", err)
	}
	w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("Created branch %s", branchName), branchName, "", "")

	// 3. Simulate Coding (Write file)
	dummyFile := filepath.Join(repoPath, "task.txt")
	_ = os.WriteFile(dummyFile, []byte(fmt.Sprintf("Work for task %s", taskId)), 0644)

	// 4. Commit
	hash, err := w.gitService.CommitChanges(repo, fmt.Sprintf("feat: implement task %s", taskId))
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("Failed to commit: %v", err), branchName, "", "")
		return
	}
	w.updateStatus(taskId, "IN_PROGRESS", "Committed changes", branchName, hash, "")

	// 5. Push (Mock)
	_ = w.gitService.Push(repo)

	// 6. Create PR
	prUrl, _ := w.ghClient.CreatePullRequest("Task "+taskId, "Implemented feature", branchName, "main")

	// 7. Complete
	w.updateStatus(taskId, "DONE", "Work completed successfully. PR Created.", branchName, hash, prUrl)
	log.Printf("Task %s completed.", taskId)
}

func (w *Worker) updateStatus(taskId, status, details, branch, commit, prUrl string) {
	_, err := w.grpcClient.UpdateTaskStatus(context.Background(), &pb.UpdateTaskStatusRequest{
		TaskId:        taskId,
		AgentId:       w.cfg.AgentID,
		Status:        status,
		Details:       details,
		GitBranch:     branch,
		GitCommitHash: commit,
		GitPrUrl:      prUrl,
	})
	if err != nil {
		log.Printf("Failed to update status: %v", err)
	}
}
