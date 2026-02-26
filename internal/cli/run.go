package cli

import (
	"fmt"
	"io"
)

const rootUsage = `wbcli is a CLI for WhiteBIT trading workflows.

Usage:
  wbcli <command> [subcommand] [flags]

Available Commands:
  keys    Manage API credential profiles
  order   Place and manage orders

Use "wbcli <command> --help" for more information about a command.`

const keysUsage = `Manage API credential profiles.

Usage:
  wbcli keys <command>

Available Commands:
  set     Store credentials for a profile
  list    List configured profiles
  remove  Remove credentials for a profile
  test    Validate credentials for a profile`

const orderUsage = `Place and manage orders.

Usage:
  wbcli order <command>

Available Commands:
  place   Place a single collateral limit order
  range   Build or submit a range order plan`

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, rootUsage)
		return 0
	}

	switch args[0] {
	case "keys":
		return runKeys(args[1:], stdout, stderr)
	case "order":
		return runOrder(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n%s\n", args[0], rootUsage)
		return 2
	}
}

func runKeys(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, keysUsage)
		return 0
	}

	switch args[0] {
	case "set", "list", "remove", "test":
		fmt.Fprintf(stdout, "wbcli keys %s is not implemented yet\n", args[0])
		return 0
	default:
		fmt.Fprintf(stderr, "unknown keys command: %s\n\n%s\n", args[0], keysUsage)
		return 2
	}
}

func runOrder(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, orderUsage)
		return 0
	}

	switch args[0] {
	case "place", "range":
		fmt.Fprintf(stdout, "wbcli order %s is not implemented yet\n", args[0])
		return 0
	default:
		fmt.Fprintf(stderr, "unknown order command: %s\n\n%s\n", args[0], orderUsage)
		return 2
	}
}
