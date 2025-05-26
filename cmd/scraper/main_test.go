package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRun_NoArgs(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	args := []string{}

	exitCode := run(out, errOut, args)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if errOut.Len() == 0 {
		t.Fatal("expected error but got none")
	}
}
func TestRun_InvalidURL(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	args := []string{"this/Is/An/Invalid/Url"}

	exitCode := run(out, errOut, args)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if errOut.Len() == 0 {
		t.Fatal("expected error but got none")
	}
}
func TestRun_LoadRulesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html><body><h1>Test</h1></body></html>`))
	}))
	defer server.Close()

	rulesJSON := `[
  		{ "query": "h1 riority": 1 },
  		{ "qery" "p", "pririty":  }
	]`
	tmpFile, err := os.CreateTemp("", "rules-*.json")
	if err != nil {
    		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up
	if _, err := tmpFile.Write([]byte(rulesJSON)); err != nil {
    		t.Fatal(err)
	}
	tmpFile.Close()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	args := []string{server.URL, tmpFile.Name()}

	exitCode := run(out, errOut, args)

	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if errOut.Len() == 0 {
		t.Fatal("expected error but got none")
	}
}
//func TestRun_ScrapError(t *testing.T) {}
//func TestRun_Success(t *testing.T)    {}
