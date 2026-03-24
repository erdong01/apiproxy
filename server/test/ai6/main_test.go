package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadDataFromCSV(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.csv")
	csvContent := "\ufeffid,task_id,input_prompt\n" +
		"1,  cgt-123  ,\"{\"\"content\"\":[{\"\"type\"\":\"\"text\"\",\"\"text\"\":\"\"hello\"\"},{\"\"type\"\":\"\"video_url\"\",\"\"video_url\"\":{\"\"url\"\":\"\"https://example.com/demo.mp4\"\"}}]}\"\n" +
		"2,,\n"

	if err := os.WriteFile(filePath, []byte(csvContent), 0o644); err != nil {
		t.Fatalf("write temp csv failed: %v", err)
	}

	var data []Data
	if err := ReadDataFromCSV(filePath, &data); err != nil {
		t.Fatalf("ReadDataFromCSV failed: %v", err)
	}

	if len(data) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data))
	}

	if data[0].TaskId != "cgt-123" {
		t.Fatalf("expected task id to be trimmed, got %q", data[0].TaskId)
	}

	if len(data[0].InputPrompt.Content) != 2 {
		t.Fatalf("expected 2 content items, got %d", len(data[0].InputPrompt.Content))
	}

	if data[0].InputPrompt.Content[1].VideoURL == nil || data[0].InputPrompt.Content[1].VideoURL.URL != "https://example.com/demo.mp4" {
		t.Fatalf("expected parsed video url, got %+v", data[0].InputPrompt.Content[1].VideoURL)
	}

	if data[1].TaskId != "" {
		t.Fatalf("expected empty task id for second row, got %q", data[1].TaskId)
	}
}

func TestQueryTasksByIDsFallsBackToDetailQuery(t *testing.T) {
	originalURL := queryTaskListURL
	originalClientFactory := newTaskQueryHTTPClient
	t.Cleanup(func() {
		queryTaskListURL = originalURL
		newTaskQueryHTTPClient = originalClientFactory
	})

	newTaskQueryHTTPClient = func() *http.Client {
		return &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch {
				case r.URL.Path == "/tasks":
					if got := r.URL.Query().Get("filter.task_ids"); got != "cgt-1,cgt-2" {
						t.Fatalf("unexpected filter.task_ids: %q", got)
					}
					if got := r.URL.Query().Get("page_num"); got != "1" {
						t.Fatalf("unexpected page_num: %q", got)
					}
					return jsonResponse(t, http.StatusOK, taskListResponse{
						Items: []task{},
						Total: 0,
					}), nil
				case strings.HasPrefix(r.URL.Path, "/tasks/"):
					taskID := strings.TrimPrefix(r.URL.Path, "/tasks/")
					return jsonResponse(t, http.StatusOK, task{
						ID:     taskID,
						Model:  "doubao-test",
						Status: "succeeded",
						Usage: Usage{
							TotalTokens: 123,
						},
					}), nil
				default:
					return jsonResponse(t, http.StatusNotFound, map[string]string{"error": "not found"}), nil
				}
			}),
		}
	}

	queryTaskListURL = "http://example.invalid/tasks"

	items, err := QueryTasksByIDs("demo-key", []string{"cgt-1", "cgt-2"})
	if err != nil {
		t.Fatalf("QueryTasksByIDs failed: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].ID != "cgt-1" || items[1].ID != "cgt-2" {
		t.Fatalf("unexpected item order: %+v", items)
	}
}

func TestQueryTasksByIDsSkipsNonArkTaskIDs(t *testing.T) {
	originalURL := queryTaskListURL
	originalClientFactory := newTaskQueryHTTPClient
	t.Cleanup(func() {
		queryTaskListURL = originalURL
		newTaskQueryHTTPClient = originalClientFactory
	})

	newTaskQueryHTTPClient = func() *http.Client {
		return &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				switch {
				case r.URL.Path == "/tasks":
					if got := r.URL.Query().Get("filter.task_ids"); got != "cgt-1" {
						t.Fatalf("unexpected filter.task_ids: %q", got)
					}
					return jsonResponse(t, http.StatusOK, taskListResponse{Items: []task{}, Total: 0}), nil
				case r.URL.Path == "/tasks/cgt-1":
					return jsonResponse(t, http.StatusOK, task{ID: "cgt-1", Model: "doubao-test"}), nil
				case r.URL.Path == "/tasks/cgt-2":
					return jsonResponse(t, http.StatusOK, task{ID: "cgt-2", Model: "doubao-test"}), nil
				default:
					t.Fatalf("unexpected path: %s", r.URL.Path)
					return nil, nil
				}
			}),
		}
	}

	queryTaskListURL = "http://example.invalid/tasks"

	items, err := QueryTasksByIDs("demo-key", []string{"021774029669670ee0befab39272c2daa6a928a2ff2cde194c47a", "cgt-1", "2035163017636945921"})
	if err != nil {
		t.Fatalf("QueryTasksByIDs failed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != "cgt-1" {
		t.Fatalf("unexpected item: %+v", items[0])
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func jsonResponse(t *testing.T, status int, payload any) *http.Response {
	t.Helper()

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal response failed: %v", err)
	}

	return &http.Response{
		StatusCode: status,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}
}
