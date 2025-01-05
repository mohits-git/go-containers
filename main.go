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
  case "child":
    child()
	default:
		panic("bad comment")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:], "as", os.Getpid())

  cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runtime.GOOS != "linux" {
		fmt.Println("this operating system is not supported")
		os.Exit(0)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
    Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v\n", os.Args[2:], "as", os.Getpid())

  syscall.Sethostname([]byte("mycontainer"))
  // pwd
  path, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  syscall.Chroot(path + "/rootfs")
  syscall.Chdir("/")
  syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

  syscall.Unmount("/proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
