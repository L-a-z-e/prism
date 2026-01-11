package github

import "fmt"

type GitHubClient interface {
	CreatePullRequest(title, body, head, base string) (string, error)
}

type MockGitHubClient struct{}

func NewMockGitHubClient() GitHubClient {
	return &MockGitHubClient{}
}

func (c *MockGitHubClient) CreatePullRequest(title, body, head, base string) (string, error) {
	// Mock implementation
	return fmt.Sprintf("https://github.com/example/repo/pull/%d", 123), nil
}
