package git

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitService interface {
	InitOrOpen(path string) (*git.Repository, error)
	CreateBranch(repo *git.Repository, branchName string) error
	CommitChanges(repo *git.Repository, message string) (string, error)
	Push(repo *git.Repository) error
}

type GitServiceImpl struct{}

func NewGitService() GitService {
	return &GitServiceImpl{}
}

func (s *GitServiceImpl) InitOrOpen(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err == nil {
		return repo, nil
	}
	return git.PlainInit(path, false)
}

func (s *GitServiceImpl) CreateBranch(repo *git.Repository, branchName string) error {
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Create new branch
	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(refName, headRef.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return err
	}

	return w.Checkout(&git.CheckoutOptions{
		Branch: refName,
	})
}

func (s *GitServiceImpl) CommitChanges(repo *git.Repository, message string) (string, error) {
	w, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// Add all changes
	if _, err := w.Add("."); err != nil {
		return "", err
	}

	// Commit
	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Prism Agent",
			Email: "agent@prism.local",
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", err
	}

	return commit.String(), nil
}

func (s *GitServiceImpl) Push(repo *git.Repository) error {
	// Mock Push for sandbox environment since we don't have remote credentials
	return nil
	// Real implementation:
	// return repo.Push(&git.PushOptions{})
}
