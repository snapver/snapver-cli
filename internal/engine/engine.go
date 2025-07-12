package engine

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	root         string
	branchName   string
	pendingFiles map[string]struct{}
	pendingMu    sync.Mutex
}

// NewEngine creates a new Engine with the given root directory
func NewEngine(root string) (*Engine, error) {
	branchName, err := getCurrentBranch()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %v", err)
	}

	return &Engine{
		root:         root,
		branchName:   branchName,
		pendingFiles: make(map[string]struct{}),
	}, nil
}

// StartInBackground runs the watch loop and batch flusher in goroutines
func (d *Engine) StartInBackground() {
	go d.watchLoop()
	go d.flushLoop()
	fmt.Println("📡 Snapver is watching...")
}

// flushLoop periodically commits all pending files
func (d *Engine) flushLoop() {
	for {
		time.Sleep(2 * time.Second)
		d.flushPendingFiles()
	}
}

func (d *Engine) flushPendingFiles() {
	d.pendingMu.Lock()
	files := make([]string, 0, len(d.pendingFiles))
	for f := range d.pendingFiles {
		files = append(files, f)
	}
	d.pendingFiles = make(map[string]struct{})
	d.pendingMu.Unlock()

	if len(files) == 0 {
		return
	}

	var added []string
	for _, filePath := range files {
		if strings.Contains(filePath, ".snapver") || strings.Contains(filePath, ".git") {
			continue
		}
		gitignorePath := filepath.Join(d.root, ".gitignore")
		ensureSnapverInGitignore(gitignorePath)
		if isIgnored(filePath) {
			continue
		}
		relPath, err := filepath.Rel(d.root, filePath)
		if err != nil {
			log.Println("❌ Failed to get relative path:", err)
			continue
		}
		if err := run("git", "add", relPath); err != nil {
			log.Println("❌ git add failed:", err)
			continue
		}
		added = append(added, relPath)
	}

	if len(added) == 0 {
		return
	}

	msg := fmt.Sprintf("snapver: auto-commit %d files", len(added))
	if err := run("git", "commit", "-m", msg); err != nil {
		log.Println("❌ git commit failed (maybe no diff):", err)
		return
	}
	log.Println("✓ Committed:", msg)

	commitHash, err := latestCommitHash()
	if err != nil {
		log.Println("❌ failed to get latest commit hash:", err)
		return
	}
	diffOut, err := exec.Command("git", "diff", fmt.Sprintf("%s^", commitHash), commitHash).Output()
	if err != nil {
		log.Println("❌ git diff failed:", err)
		return
	}
	if err := saveDiff(d.root, d.branchName, commitHash, msg, string(diffOut)); err != nil {
		log.Println(err)
	}
}

// handleFileChange just adds the file to pendingFiles
func (d *Engine) handleFileChange(filePath string) {
	d.pendingMu.Lock()
	d.pendingFiles[filePath] = struct{}{}
	d.pendingMu.Unlock()
}
