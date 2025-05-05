/// Used in `justfile`
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if runtime.GOOS == "darwin" {
		homeDir, _ := os.UserHomeDir()
		fmt.Printf("%s/Library/Caches/merlion/merlion.log", homeDir)
	} else {
		fmt.Print("~/.cache/merlion/merlion.log")
	}
}
