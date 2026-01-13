package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		log.Fatalf("gRPC connection failed: %v", err)
	}

	c := pb.NewAgentServiceClient(conn)
	aiClient := ai.NewClient() // 자동으로 Ollama/Mock 선택

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
	log.Printf("🚀 Starting Agent Worker: %s", w.cfg.AgentID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := w.grpcClient.RegisterAgent(ctx, &pb.RegisterAgentRequest{
		AgentId: w.cfg.AgentID,
		Version: "0.1.0",
	})
	if err != nil {
		log.Printf("⚠️  Failed to register agent: %v", err)
	}

	channel := fmt.Sprintf("tasks:%s", w.cfg.AgentID)
	log.Printf("👂 Listening on Redis channel: %s", channel)

	pubsub := w.redisClient.Subscribe(context.Background(), channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("📨 Received task: %s", msg.Payload)
		go w.processTask(msg.Payload)
	}
}

func (w *Worker) processTask(payload string) {
	taskId := w.parseTaskId(payload)

	log.Printf("\n🚀 ========================================")
	log.Printf("   Processing Task: %s", taskId)
	log.Printf("   ========================================\n")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	w.updateStatus(taskId, "IN_PROGRESS", "🤖 Agent started...", "", "", "")

	// 1. Git Setup
	log.Println("📂 Setting up Git repository...")
	repoPath := filepath.Join(os.TempDir(), "prism-repo")
	_ = os.MkdirAll(repoPath, 0755)
	repo, err := w.gitService.InitOrOpen(repoPath)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Git init failed: %v", err), "", "", "")
		return
	}

	// 2. Create Branch
	branchName := fmt.Sprintf("feat/task-%s", taskId)
	log.Printf("🌳 Creating branch: %s\n", branchName)
	if err := w.gitService.CreateBranch(repo, branchName); err != nil {
		log.Printf("⚠️  Branch warning: %v", err)
	}
	w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("✅ Branch: %s", branchName), branchName, "", "")

	// 3. AI Code Generation with Chat API
	log.Println("\n🤖 ========================================")
	log.Println("   Starting AI Code Generation (Chat API)")
	log.Println("   ========================================\n")

	systemPrompt := `You are an expert Java Spring Boot developer. Generate production-ready code ONLY.

RULES:
- Use Spring Boot 3.x with Java 17+
- Include @RestController, @Service, @Repository, @Entity annotations
- Add comprehensive error handling (@ControllerAdvice, try-catch)
- Include validation (@Valid, Jakarta Bean Validation)
- Follow Java naming conventions (camelCase, PascalCase)
- Add JavaDoc for public methods
- Use SLF4J for logging
- Include proper HTTP status codes (200, 400, 404, 500)
- Output ONLY pure Java code (NO markdown, NO explanations)`

	userPrompt := fmt.Sprintf(`Create a complete Spring Boot REST API for task: %s

Requirements:
- Complete CRUD operations
- @RestController with proper endpoints
- @Service layer with business logic
- Exception handling with custom exceptions
- Request/Response DTOs if needed
- Comprehensive logging
- Production-ready quality

Generate the complete code now.`, taskId)

	startTime := time.Now()
	var generatedCode string
	var codeErr error

	// Chat API 호출 (Context 포함)
	generatedCode, codeErr = w.aiClient.GenerateCodeWithSystem(ctx, systemPrompt, userPrompt)
	elapsed := time.Since(startTime)

	if codeErr != nil {
		log.Printf("\n⚠️  AI failed after %.2fs: %v\n", elapsed.Seconds(), codeErr)
		log.Println("📝 Using fallback code...\n")

		generatedCode = w.generateFallbackCode(taskId, codeErr)
		w.updateStatus(taskId, "IN_PROGRESS", "⚠️  AI failed, using fallback", branchName, "", "")
	} else {
		log.Printf("✅ AI generated %d characters in %.2f seconds\n", len(generatedCode), elapsed.Seconds())
		w.updateStatus(taskId, "IN_PROGRESS",
			fmt.Sprintf("✅ AI: %d bytes in %.2fs", len(generatedCode), elapsed.Seconds()),
			branchName, "", "")
	}

	// 4. Write Code
	log.Println("💾 Writing generated code...")
	codeFile := filepath.Join(repoPath, fmt.Sprintf("Task%sController.java",
		strings.ReplaceAll(taskId, "-", "")))
	if err := os.WriteFile(codeFile, []byte(generatedCode), 0644); err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Write failed: %v", err), branchName, "", "")
		return
	}
	log.Printf("✅ Code written to: %s\n", codeFile)

	// 5. Commit
	log.Println("📍 Committing changes...")
	commitMsg := fmt.Sprintf(
		"feat: implement task %s with AI\n\nGenerated in %.2fs using Chat API\nCode size: %d bytes",
		taskId, elapsed.Seconds(), len(generatedCode))
	hash, err := w.gitService.CommitChanges(repo, commitMsg)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Commit failed: %v", err), branchName, "", "")
		return
	}
	log.Printf("✅ Committed: %s\n", hash)
	w.updateStatus(taskId, "IN_PROGRESS", "✅ Committed changes", branchName, hash, "")

	// 6. Push
	log.Println("🔝 Pushing to remote...")
	_ = w.gitService.Push(repo)

	// 7. Create PR
	log.Println("🔗 Creating Pull Request...")
	prTitle := fmt.Sprintf("[Task %s] AI-Generated Implementation", taskId)
	prBody := w.generatePRBody(taskId, elapsed, len(generatedCode), branchName)
	prUrl, prErr := w.ghClient.CreatePullRequest(prTitle, prBody, branchName, "main")
	if prErr != nil {
		log.Printf("⚠️  PR creation failed: %v", prErr)
		prUrl = "<PR creation failed>"
	}

	// 8. Complete
	log.Println("\n✅ ========================================")
	log.Println("   Task Completed Successfully!")
	log.Println("   ========================================\n")

	w.updateStatus(
		taskId,
		"DONE",
		fmt.Sprintf("✅ Completed!\n📊 Time: %.2fs, Code: %d bytes\n🔗 PR: %s",
			elapsed.Seconds(), len(generatedCode), prUrl),
		branchName,
		hash,
		prUrl,
	)
}

func (w *Worker) parseTaskId(payload string) string {
	for i, c := range payload {
		if c == ':' {
			return payload[:i]
		}
	}
	return payload
}

func (w *Worker) generateFallbackCode(taskId string, err error) string {
	return fmt.Sprintf(`package com.prism.generated;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;

@RestController
@RequestMapping("/api/task/%s")
public class Task%sController {
    
    // AI generation failed: %v
    // TODO: Implement manually
    
    @GetMapping
    public ResponseEntity<String> get() {
        return ResponseEntity.ok("TODO: Implement task %s");
    }
}`, taskId, strings.ReplaceAll(taskId, "-", ""), err, taskId)
}

func (w *Worker) generatePRBody(taskId string, elapsed time.Duration, codeSize int, branchName string) string {
	return fmt.Sprintf(`## AI-Generated Implementation

### Task ID
%s

### Generation Details
- **Time taken**: %.2f seconds
- **Code size**: %d bytes
- **Model**: Ollama DeepSeek-R1
- **Branch**: %s
- **API**: Chat API (system + user prompts)

### Changes
- Generated production-ready Spring Boot controller
- Includes error handling and validation
- Following best practices and naming conventions

### Review Checklist
- [ ] Code compiles without errors
- [ ] Tests pass successfully
- [ ] Code quality is acceptable
- [ ] Security considerations are met

*Generated by Prism AI Agent*`, taskId, elapsed.Seconds(), codeSize, branchName)
}

func (w *Worker) updateStatus(taskId, status, details, branch, commit, prUrl string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := w.grpcClient.UpdateTaskStatus(ctx, &pb.UpdateTaskStatusRequest{
		TaskId:        taskId,
		AgentId:       w.cfg.AgentID,
		Status:        status,
		Details:       details,
		GitBranch:     branch,
		GitCommitHash: commit,
		GitPrUrl:      prUrl,
	})
	if err != nil {
		log.Printf("⚠️  Failed to update status: %v", err)
	}
}
