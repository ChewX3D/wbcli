package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func executeCommandWithInput(input string, args ...string) (string, string, error) {
	command := NewRootCmdForTest()
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	command.SetOut(stdout)
	command.SetErr(stderr)
	if input != "" {
		command.SetIn(strings.NewReader(input))
	}
	command.SetArgs(args)

	err := command.Execute()

	return stdout.String(), stderr.String(), err
}

func executeCommand(args ...string) (string, string, error) {
	return executeCommandWithInput("", args...)
}

func assertUnknownAuthSubcommand(t *testing.T, subcommand string) {
	t.Helper()

	_, _, err := executeCommand("auth", subcommand)
	if err == nil {
		t.Fatalf("expected unknown command error for %q", subcommand)
	}
	if !strings.Contains(err.Error(), "unknown command \""+subcommand+"\"") {
		t.Fatalf("unexpected error for %q: %v", subcommand, err)
	}
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

func TestLegacyAuthSetCommandRemoved(t *testing.T) {
	assertUnknownAuthSubcommand(t, "set")
}

func TestLegacyAuthUseCommandRemoved(t *testing.T) {
	assertUnknownAuthSubcommand(t, "use")
}

func TestLegacyAuthListCommandRemoved(t *testing.T) {
	assertUnknownAuthSubcommand(t, "list")
}

func TestLegacyAuthCurrentCommandRemoved(t *testing.T) {
	assertUnknownAuthSubcommand(t, "current")
}

func TestAuthLoginRejectsInvalidStdinContract(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	_, _, err := executeCommandWithInput("only-key\n", "auth", "login")
	if err == nil {
		t.Fatal("expected stdin parsing error")
	}
	if !strings.Contains(err.Error(), "exactly two non-empty lines") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAuthStatusWorksWhenLoggedOut(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	stdout, _, err := executeCommand("auth", "status")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(stdout, "logged_in=false") {
		t.Fatalf("expected logged_out status, got: %q", stdout)
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
