package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

/// What we are building?
// docker         run image <cmd> <params>
// building equivalent:
// go run main.go run       <cmd> <params>

func main() {
	if len(os.Args) < 2 {
		panic("specify a command")
	}

	switch os.Args[1] {
	case "run":
		run()

	default:
		panic("bad comment")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runtime.GOOS == "linux" {
		fmt.Println("Warning: If you're on Linux, uncomment the code manually")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			// Cloneflags: syscall.CLONE_NEWUTS,
		}
	} else if runtime.GOOS == "darwin" {
		fmt.Println("Warning: macOS does not support Linux namespaces")
		cmd.SysProcAttr = &syscall.SysProcAttr{
    }
	}

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
