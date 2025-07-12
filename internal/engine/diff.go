package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// saveDiff saves the diff data as a JSON file in the branch-specific diff directory
func saveDiff(root, branchName, commitHash, msg, diff string) error {
	diffDir := filepath.Join(root, ".snapver", "diffs", branchName)
	if err := os.MkdirAll(diffDir, 0755); err != nil {
		return fmt.Errorf("failed to create branch diff directory: %w", err)
	}

	data := map[string]string{
		"commit":  commitHash,
		"message": msg,
		"diff":    diff,
	}
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")

	utcTime := time.Now().UTC().Format("20060102T150405Z")
	fileName := fmt.Sprintf("snapver_%s.json", utcTime)
	outFile := filepath.Join(diffDir, fileName)

	if err := os.WriteFile(outFile, jsonBytes, 0644); err != nil {
		return fmt.Errorf("failed to save diff JSON: %w", err)
	}

	log.Println("✓ Diff saved:", outFile)
	return nil
}
