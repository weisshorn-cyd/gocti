//go:build ignore

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("expected arguments: [path to 'tools' folder], received: %s", os.Args[1:])
	}

	toolsFolder := os.Args[1]

	generator := filepath.Join(toolsFolder, "gocti_type_generator/gocti_type_generator/main.py")

	args := []string{generator}

	command := exec.Command("python3", args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Run()
}
