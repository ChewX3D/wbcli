package cli

import "os"

// IsTerminalInput reports whether provided file is attached to a terminal.
func IsTerminalInput(file *os.File) bool {
	if file == nil {
		return false
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
