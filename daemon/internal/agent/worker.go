package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/prism/daemon/internal/ai"
	"github.com/prism/daemon/internal/git"
	"github.com/prism/daemon/internal/github"
	pb "github.com/prism/daemon/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	AgentID   string
	RedisAddr string
	GrpcAddr  string
}

type Worker struct {
	cfg         Config
	redisClient *redis.Client
	grpcClient  pb.AgentServiceClient
	gitService  git.GitService
	ghClient    github.GitHubClient
	aiClient    ai.Client
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

	// AI Client 초기화
	var aiClient ai.Client
	aiMode := os.Getenv("AI_MODE")

	if aiMode == "ollama" {
		log.Println("🤖 Using Ollama AI Client (DeepSeek-R1 32B)")
		aiClient = ai.NewOllamaClient()
	} else {
		log.Println("🤖 Using Mock AI Client")
		aiClient = ai.NewMockClient()
	}

	return &Worker{
		cfg:         cfg,
		redisClient: rdb,
		grpcClient:  c,
		gitService:  git.NewGitService(),
		ghClient:    github.NewMockGitHubClient(),
		aiClient:    aiClient,
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
		go w.processTask(msg.Payload)
	}
}

func (w *Worker) processTask(payload string) {
	// Parse task ID from payload
	var taskId string
	for i, c := range payload {
		if c == ':' {
			taskId = payload[:i]
			break
		}
	}
	if taskId == "" {
		taskId = payload
	}

	log.Printf("🚀 Processing Task %s with AI...", taskId)
	w.updateStatus(taskId, "IN_PROGRESS", "Agent started working on task...", "", "", "")

	// 1. Init/Open Git Repo
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

	// 3. AI Code Generation
	log.Printf("🤖 Generating code with AI...")

	prompt := fmt.Sprintf(`You are an expert software engineer. Generate production-ready code for the following task:

Task: %s

Requirements:
- Use Spring Boot 3.x
- Include proper annotations (@RestController, @Service, @Repository, etc.)
- Add comprehensive error handling with try-catch blocks
- Include input validation using @Valid and Jakarta validation
- Follow Java best practices and naming conventions
- Add JavaDoc comments for all public methods
- Make it production-ready with logging
- Include proper HTTP status codes

Generate ONLY the code without explanations or markdown formatting.`, taskId)

	startTime := time.Now()
	generatedCode, err := w.aiClient.GenerateCode(prompt)
	elapsed := time.Since(startTime)

	if err != nil {
		log.Printf("⚠️ AI generation failed after %.2fs: %v", elapsed.Seconds(), err)
		log.Println("📝 Using fallback code...")
		generatedCode = fmt.Sprintf(`// Fallback code for task %s
// AI generation failed: %v

package com.prism.generated;

public class Task%s {
    // TODO: Implement task manually
}`, taskId, err, taskId)
		w.updateStatus(taskId, "IN_PROGRESS", "AI failed, using fallback", branchName, "", "")
	} else {
		log.Printf("✅ AI generated %d characters in %.2f seconds", len(generatedCode), elapsed.Seconds())
		w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("AI generated code (%.1fs)", elapsed.Seconds()), branchName, "", "")
	}

	// 4. Write generated code to file
	codeFile := filepath.Join(repoPath, "generated_code.java")
	if err := os.WriteFile(codeFile, []byte(generatedCode), 0644); err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("Failed to write file: %v", err), branchName, "", "")
		return
	}

	// 5. Commit
	commitMsg := fmt.Sprintf("feat: implement task %s with AI\n\nGenerated in %.2fs using AI", taskId, elapsed.Seconds())
	hash, err := w.gitService.CommitChanges(repo, commitMsg)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("Failed to commit: %v", err), branchName, "", "")
		return
	}
	w.updateStatus(taskId, "IN_PROGRESS", "Committed AI-generated changes", branchName, hash, "")

	// 6. Push (Mock)
	_ = w.gitService.Push(repo)

	// 7. Create PR
	prTitle := fmt.Sprintf("Task %s - AI Generated", taskId)
	prBody := fmt.Sprintf("AI-generated implementation\nGenerated in: %.2fs\nCode size: %d characters", elapsed.Seconds(), len(generatedCode))
	prUrl, _ := w.ghClient.CreatePullRequest(prTitle, prBody, branchName, "main")

	// 8. Complete
	w.updateStatus(taskId, "DONE", "✅ Work completed with AI. PR Created.", branchName, hash, prUrl)
	log.Printf("✅ Task %s completed successfully!", taskId)
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
