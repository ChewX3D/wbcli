package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRootHelp(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer

	exitCode := Run(nil, &out, &err)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(out.String(), "Available Commands:") {
		t.Fatalf("expected help output, got: %q", out.String())
	}

	if err.Len() != 0 {
		t.Fatalf("expected no stderr output, got: %q", err.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer

	exitCode := Run([]string{"does-not-exist"}, &out, &err)

	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}

	if !strings.Contains(err.String(), "unknown command") {
		t.Fatalf("expected unknown command error, got: %q", err.String())
	}
}

func TestRunKeysStub(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer

	exitCode := Run([]string{"keys", "list"}, &out, &err)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(out.String(), "not implemented yet") {
		t.Fatalf("expected stub output, got: %q", out.String())
	}

	if err.Len() != 0 {
		t.Fatalf("expected no stderr output, got: %q", err.String())
	}
}

func TestRunOrderStub(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer

	exitCode := Run([]string{"order", "place"}, &out, &err)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	if !strings.Contains(out.String(), "not implemented yet") {
		t.Fatalf("expected stub output, got: %q", out.String())
	}

	if err.Len() != 0 {
		t.Fatalf("expected no stderr output, got: %q", err.String())
	}
}
