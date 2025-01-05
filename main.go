package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
)

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
	fmt.Printf("Running %v as %v\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runtime.GOOS != "linux" {
		fmt.Println("this operating system is not supported")
		os.Exit(0)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v as %v\n", os.Args[2:], os.Getpid())

	cg()

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

// cg sets up the cgroup
func cg() error {
	cgroupPath := "/sys/fs/cgroup"
	containerPath := filepath.Join(cgroupPath, "liz")

	// create the cgroup directory
	if err := os.MkdirAll(containerPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %v", err)
	}

	// set the maximum number of processes
	if err := os.WriteFile(filepath.Join(containerPath, "pids.max"), []byte("20"), 0644); err != nil {
		return fmt.Errorf("failed to set pids.max: %v", err)
	}

	// add current process to the cgroup
	if err := os.WriteFile(filepath.Join(containerPath, "cgroup.procs"),
		[]byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %v", err)
	}

	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
