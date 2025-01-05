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

func cg() {
	pids := "/sys/fs/cgroup/pids"
	err := os.Mkdir(filepath.Join(pids, "liz"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	must(os.WriteFile(filepath.Join(pids, "liz/pids.max"), []byte("20"), 0700))
	must(os.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700))
	must(os.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
