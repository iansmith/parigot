package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

const (
	desiredMemory = 50 // The whole memory size you want to consume (in megabytes)
)

func main() {
	// Calculate the number of iterations needed
	numIterations := desiredMemory / getProcessMemory()

	fmt.Printf("Forking %d times to consume approximately %dMB of memory.\n", numIterations, desiredMemory)

	// forks itself repeatedly to consume the memory.
	for i := 0; i < numIterations; i++ {
		forkAndExec()
	}

	fmt.Println("All forks completed.")
}

func getProcessMemory() int {
	var usage runtime.MemStats
	runtime.ReadMemStats(&usage)

	// Convert bytes to megabytes
	usageMb := int(usage.Sys / 1024 / 1024)

	fmt.Printf("Sys = %v Mb\n", usageMb)

	return usageMb
}

func forkAndExec() {
	// fork the current process
	child, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to fork: %v\n", err)
	}

	// Prepare command and arguments for the child process
	cmd := exec.Command(child)

	// Set environment variables, if necessary
	cmd.Env = os.Environ()

	// Start the child process
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start child process: %v\n", err)
		return
	}

	// Print the PID of the child process
	fmt.Printf("Child process started with PID: %d\n", cmd.Process.Pid)
}
