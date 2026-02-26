package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func executeCommand(args ...string) (string, string, error) {
	command := NewRootCmdForTest()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command.SetOut(stdout)
	command.SetErr(stderr)
	command.SetArgs(args)

	err := command.Execute()

	return stdout.String(), stderr.String(), err
}

func TestRootHelpShowsMainGroups(t *testing.T) {
	stdout, stderr, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(stdout, "auth") || !strings.Contains(stdout, "order") {
		t.Fatalf("expected auth and order in help output, got: %q", stdout)
	}

	if stderr != "" {
		t.Fatalf("expected empty stderr, got: %q", stderr)
	}
}

func TestUnknownCommandReturnsError(t *testing.T) {
	_, _, err := executeCommand("unknown")
	if err == nil {
		t.Fatal("expected an error for unknown command")
	}
}

func TestAuthSetRequiresSecrets(t *testing.T) {
	_, _, err := executeCommand("auth", "set", "--profile", "default")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "--api-key is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLegacyKeysAliasStillWorks(t *testing.T) {
	_, _, err := executeCommand("keys", "set", "--profile", "default")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "--api-key is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOrderPlaceStub(t *testing.T) {
	stdout, _, err := executeCommand(
		"order", "place",
		"--profile", "default",
		"--market", "BTC_PERP",
		"--side", "buy",
		"--amount", "0.01",
		"--price", "50000",
		"--expiration", "0",
		"--client-order-id", "my-order-001",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(stdout, "not implemented yet") {
		t.Fatalf("expected stub output, got: %q", stdout)
	}
}

func TestOrderRangeStub(t *testing.T) {
	stdout, _, err := executeCommand(
		"order", "range",
		"--profile", "default",
		"--market", "BTC_PERP",
		"--side", "buy",
		"--start-price", "49000",
		"--end-price", "50000",
		"--step", "50",
		"--amount-mode", "constant",
		"--base-amount", "0.005",
		"--dry-run",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(stdout, "not implemented yet") {
		t.Fatalf("expected stub output, got: %q", stdout)
	}
}
