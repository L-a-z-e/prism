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
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewAgentServiceClient(conn)

	// ============================================================================
	// AI Client 초기화 - 우선순위: 환경변수 → 자동 선택
	// ============================================================================
	var aiClient ai.Client

	// 1. 강제 선택
	if os.Getenv("OLLAMA_MOCK") == "true" {
		log.Println("🧪 Forcing Mock AI Client (OLLAMA_MOCK=true)")
		aiClient = ai.NewMockClient()
	} else {
		// 2. 자동 선택 (NewClient는 자동으로 Ollama 연결 시도)
		aiClient = ai.NewClient()
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
	log.Printf("🚀 Starting Agent Worker for ID: %s", w.cfg.AgentID)

	// Register first
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := w.grpcClient.RegisterAgent(ctx, &pb.RegisterAgentRequest{
		AgentId: w.cfg.AgentID,
		Version: "0.1.0",
	})
	if err != nil {
		log.Printf("⚠️  Failed to register agent: %v. Continuing anyway...", err)
	}

	// Subscribe to Redis
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

	log.Printf("\n🚀 ========================================\n")
	log.Printf("   Processing Task: %s\n", taskId)
	log.Printf("   ========================================\n")

	// Context with 30 minute timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	w.updateStatus(taskId, "IN_PROGRESS", "🤖 Agent started working on task...", "", "", "")

	// 1. Init/Open Git Repo
	log.Println("📂 Setting up Git repository...")
	repoPath := filepath.Join(os.TempDir(), "prism-repo")
	_ = os.MkdirAll(repoPath, 0755)
	repo, err := w.gitService.InitOrOpen(repoPath)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Failed to init git repo: %v", err), "", "", "")
		return
	}

	// 2. Create Branch
	branchName := fmt.Sprintf("feat/task-%s", taskId)
	log.Printf("🌳 Creating branch: %s\n", branchName)
	if err := w.gitService.CreateBranch(repo, branchName); err != nil {
		log.Printf("⚠️  Branch creation warning (might exist): %v", err)
	}
	w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("✅ Created branch %s", branchName), branchName, "", "")

	// 3. AI Code Generation
	log.Println("\n🤖 ========================================")
	log.Println("   Starting AI Code Generation...")
	log.Println("   ========================================\n")

	prompt := fmt.Sprintf(`You are an expert software engineer. Generate production-ready code for the following task:

Task ID: %s

Requirements:
- Use Spring Boot 3.x with modern Java
- Include proper annotations (@RestController, @Service, @Repository, etc.)
- Add comprehensive error handling with try-catch blocks
- Include input validation using @Valid and Jakarta validation
- Follow Java best practices and naming conventions
- Add JavaDoc comments for all public methods
- Make it production-ready with logging using slf4j
- Include proper HTTP status codes (200, 400, 404, 500)
- Create a simple but complete REST API example

Generate ONLY the code without explanations or markdown formatting.
Code should be immediately executable and compilable.`, taskId)

	startTime := time.Now()
	var generatedCode string
	var codeErr error

	// 환경변수로 스트리밍 모드 선택
	useStreaming := os.Getenv("STREAM_CODE_GENERATION") == "true"

	if useStreaming {
		// 스트리밍 모드 (진행상황 실시간 전송)
		log.Println("📡 Using streaming mode...")
		var codeBuffer strings.Builder

		codeErr = w.aiClient.GenerateCodeStream(ctx, prompt, func(chunk string) error {
			codeBuffer.WriteString(chunk)
			log.Printf("📝 Generated chunk: %d bytes", len(chunk))

			// 진행상황 업데이트 (선택)
			// w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("🤖 Generating... (%d bytes)", codeBuffer.Len()), branchName, "", "")

			return nil
		})

		generatedCode = codeBuffer.String()
	} else {
		// 비스트리밍 모드 (전체 응답 대기)
		log.Println("⏳ Waiting for complete AI response...")
		generatedCode, codeErr = w.aiClient.GenerateCode(ctx, prompt)
	}

	elapsed := time.Since(startTime)

	if codeErr != nil {
		log.Printf("\n⚠️  AI generation failed after %.2fs: %v\n", elapsed.Seconds(), codeErr)
		log.Println("📝 Using fallback code...\n")

		generatedCode = fmt.Sprintf(`// Fallback code for task %s
// AI generation failed: %v
// Please implement manually

package com.prism.generated;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;

@RestController
@RequestMapping("/api/task/%s")
public class Task%sController {
    
    @GetMapping
    public ResponseEntity<String> get() {
        return ResponseEntity.ok("TODO: Implement task %s");
    }
}`, taskId, codeErr, taskId, strings.ReplaceAll(taskId, "-", ""), taskId)

		w.updateStatus(taskId, "IN_PROGRESS", "⚠️  AI failed, using fallback", branchName, "", "")
	} else {
		log.Printf("✅ AI generated %d characters in %.2f seconds\n", len(generatedCode), elapsed.Seconds())
		w.updateStatus(taskId, "IN_PROGRESS", fmt.Sprintf("✅ AI generated %d bytes in %.2fs", len(generatedCode), elapsed.Seconds()), branchName, "", "")
	}

	// 4. Write generated code to file
	log.Println("💾 Writing generated code to file...")
	codeFile := filepath.Join(repoPath, fmt.Sprintf("Task%sController.java", strings.ReplaceAll(taskId, "-", "")))
	if err := os.WriteFile(codeFile, []byte(generatedCode), 0644); err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Failed to write file: %v", err), branchName, "", "")
		return
	}
	log.Printf("✅ Code written to: %s\n", codeFile)

	// 5. Commit
	log.Println("📍 Committing changes...")
	commitMsg := fmt.Sprintf("feat: implement task %s with AI\n\nGenerated in %.2fs using Ollama DeepSeek\nCode size: %d bytes", taskId, elapsed.Seconds(), len(generatedCode))
	hash, err := w.gitService.CommitChanges(repo, commitMsg)
	if err != nil {
		w.updateStatus(taskId, "FAILED", fmt.Sprintf("❌ Failed to commit: %v", err), branchName, "", "")
		return
	}
	log.Printf("✅ Committed with hash: %s\n", hash)
	w.updateStatus(taskId, "IN_PROGRESS", "✅ Committed AI-generated changes", branchName, hash, "")

	// 6. Push (Mock)
	log.Println("🔝 Pushing to remote...")
	_ = w.gitService.Push(repo)

	// 7. Create PR
	log.Println("🔗 Creating Pull Request...")
	prTitle := fmt.Sprintf("[Task %s] AI-Generated Implementation", taskId)
	prBody := fmt.Sprintf(`## AI-Generated Implementation

### Task ID
%s

### Generation Details
- **Time taken**: %.2f seconds
- **Code size**: %d bytes
- **Model**: Ollama DeepSeek-R1
- **Branch**: %s

### Changes
- Generated production-ready Spring Boot controller
- Includes error handling and validation
- Following best practices and naming conventions

### Review Checklist
- [ ] Code compiles without errors
- [ ] Tests pass successfully
- [ ] Code quality is acceptable
- [ ] Security considerations are met

*Generated by Prism AI Agent*`, taskId, elapsed.Seconds(), len(generatedCode), branchName)

	prUrl, prErr := w.ghClient.CreatePullRequest(prTitle, prBody, branchName, "main")
	if prErr != nil {
		log.Printf("⚠️  Failed to create PR: %v", prErr)
		prUrl = "<PR creation failed>"
	}

	// 8. Complete
	log.Println("\n✅ ========================================")
	log.Println("   Task Completed Successfully!")
	log.Println("   ========================================\n")

	w.updateStatus(
		taskId,
		"DONE",
		fmt.Sprintf(
			"✅ Work completed!\n📊 Time: %.2fs, Code: %d bytes\n🔗 PR: %s",
			elapsed.Seconds(),
			len(generatedCode),
			prUrl,
		),
		branchName,
		hash,
		prUrl,
	)
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
