package e2e_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type TestResult struct {
	File   string
	Status string
	Msg    string
}

var results []TestResult

func apply(t *testing.T, file string) (string, error) {
	cmd := exec.Command("kubectl", "apply", "-f", file)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func discoverFiles(t *testing.T, dir string) []string {
	files := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Logf("Skipping file due to error: %v", err)
			return nil // Skip file instead of failing
		}
		if info == nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to discover files in %s: %v", dir, err)
	}

	if len(files) == 0 {
		t.Logf("No test files found in %s", dir)
	}

	return files
}

func TestWebhookScenarios(t *testing.T) {
	t.Cleanup(func() {
		printSummary()
	})

	runScenario(t, "pass-scenario", true)
	runScenario(t, "fail-scenario", false)
}

func runScenario(t *testing.T, dir string, shouldPass bool) {
	files := discoverFiles(t, dir)
	for _, f := range files {
		name := filepath.Base(f)
		t.Run(name, func(t *testing.T) {
			out, err := apply(t, f)

			if shouldPass && err != nil {
				results = append(results, TestResult{f, "FAIL", "Should have passed, but failed"})
				t.Errorf("FAIL (expected success): %s\n%s", f, out)
				return
			}

			if !shouldPass && err == nil {
				results = append(results, TestResult{f, "FAIL", "Should have been rejected, but passed"})
				t.Errorf("FAIL (expected rejection): %s\n%s", f, out)
				return
			}

			msg := "Applied successfully"
			if !shouldPass {
				msg = "Rejected as expected"
			}

			results = append(results, TestResult{f, "✅ PASS", msg})
			t.Logf("PASS: %s", f)
		})
	}
}

func printSummary() {
	fmt.Println("\nTest Result Summary")
	fmt.Printf("┌ %-30s ┬ %-6s ┬ %-35s ┐\n", "Test File", "Status", "Message")
	fmt.Println("├" + strings.Repeat("─", 32) + "┼" + strings.Repeat("─", 8) + "┼" + strings.Repeat("─", 37) + "┤")
	for _, r := range results {
		fmt.Printf("│ %-30s │ %-6s │ %-35s │\n", filepath.Base(r.File), r.Status, r.Msg)
	}
	fmt.Println("└" + strings.Repeat("─", 32) + "┴" + strings.Repeat("─", 8) + "┴" + strings.Repeat("─", 37) + "┘")
}
