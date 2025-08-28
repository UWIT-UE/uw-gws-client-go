package gws

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.Handler) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	cfg := &Config{
		APIUrl:        srv.URL,
		Timeout:       5,
		SkipTLSVerify: true,
	}
	c, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	return c
}

func TestRenameGroup_ResolvesRegidAndCallsMove(t *testing.T) {
	calledGet := false
	calledPut := false
	var putPath string
	var putQuery url.Values

	handler := http.NewServeMux()

	handler.HandleFunc("/group/u_dept_team", func(w http.ResponseWriter, r *http.Request) {
		calledGet = true
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Etag", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"regid":"abc123","id":"u_dept_team"}}`))
	})

	handler.HandleFunc("/groupMove/abc123", func(w http.ResponseWriter, r *http.Request) {
		calledPut = true
		putPath = r.URL.Path
		putQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	})

	client := newTestClient(t, handler)

	if err := client.RenameGroup("u_dept_team", "team_new"); err != nil {
		t.Fatalf("RenameGroup error: %v", err)
	}

	if !calledGet || !calledPut {
		t.Fatalf("expected both GET and PUT to be called")
	}
	if putPath != "/groupMove/abc123" {
		t.Fatalf("unexpected PUT path: %s", putPath)
	}
	if got := putQuery.Get("newext"); got != "team_new" {
		t.Fatalf("expected newext query param, got: %s", got)
	}
}

func TestMoveGroup_ResolvesRegidAndCallsMove(t *testing.T) {
	calledGet := false
	calledPut := false
	var putQuery url.Values

	handler := http.NewServeMux()
	handler.HandleFunc("/group/abc123", func(w http.ResponseWriter, r *http.Request) {
		// Allow passing regid directly too
		calledGet = true
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Etag", "1")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"regid":"abc123","id":"u_old_stem_leaf"}}`))
	})
	handler.HandleFunc("/groupMove/abc123", func(w http.ResponseWriter, r *http.Request) {
		calledPut = true
		putQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
	})

	client := newTestClient(t, handler)

	if err := client.MoveGroup("abc123", "u_new_stem"); err != nil {
		t.Fatalf("MoveGroup error: %v", err)
	}

	if !calledGet || !calledPut {
		t.Fatalf("expected both GET and PUT to be called")
	}
	if got := putQuery.Get("newstem"); got != "u_new_stem" {
		t.Fatalf("expected newstem query param, got: %s", got)
	}
}

func TestRenameGroup_Validation(t *testing.T) {
	client := &Client{config: &Config{Timeout: time.Second}}
	if err := client.RenameGroup("", "x"); err == nil || !strings.Contains(err.Error(), "groupID") {
		t.Fatalf("expected groupID validation error")
	}
	if err := client.RenameGroup("u_name", ""); err == nil || !strings.Contains(err.Error(), "newLeaf") {
		t.Fatalf("expected newLeaf validation error")
	}
}

func TestMoveGroup_Validation(t *testing.T) {
	client := &Client{config: &Config{Timeout: time.Second}}
	if err := client.MoveGroup("", "x"); err == nil || !strings.Contains(err.Error(), "groupID") {
		t.Fatalf("expected groupID validation error")
	}
	if err := client.MoveGroup("u_name", ""); err == nil || !strings.Contains(err.Error(), "newStem") {
		t.Fatalf("expected newStem validation error")
	}
}
