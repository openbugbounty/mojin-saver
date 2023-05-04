package util

import "os"

func IsPipedInput() bool {
	fi, _ := os.Stdin.Stat()

	return fi.Mode()&os.ModeNamedPipe != 0
}
