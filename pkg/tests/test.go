package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertBody(t *testing.T, body io.Reader, expectedFile string) {
	expectedData := readFile(t, expectedFile)
	actualData, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	actualData = normalizeJSON(t, actualData)
	expectedData = normalizeJSON(t, expectedData)

	if diff := cmp.Diff(string(expectedData), string(actualData)); diff != "" {
		t.Fatalf("response mismatch (-want +got):\n%s", diff)
	}
}

func readFile(t *testing.T, relPath string) []byte {
	projectRoot := getProjectRoot("")
	fullPath := filepath.Join(projectRoot, relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", fullPath, err)
	}
	return data
}

func getProjectRoot(pwd string) string {
	if pwd == "" {
		pwd, _ = os.Getwd()
	} else {
		pwd = fmt.Sprintf("%s/..", pwd)
	}
	if fileExists(fmt.Sprintf("%s/go.mod", pwd)) {
		return pwd
	}
	return getProjectRoot(pwd)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func normalizeJSON(t *testing.T, data []byte) []byte {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return trimmed
	}
	var payload any
	if err := json.Unmarshal(trimmed, &payload); err != nil {
		return trimmed
	}
	normalized, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("failed to normalize json: %v", err)
	}
	return normalized
}
