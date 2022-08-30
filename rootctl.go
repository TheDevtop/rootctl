package main

/*
	Prog: Manage chroot environments
	Vers: 0.1
	Auth: Thijs Haker
*/

import (
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"os/exec"
)

const (
	default_command = "/usr/bin/login"
	default_root    = "/"
)

// Generate uname
func genUname() string {
	var buf = new(unix.Utsname)
	unix.Uname(buf)
	return fmt.Sprint(string((*buf).Version[:]) + string((*buf).Machine[:]))
}

// Boot system with new root
// Otherwise panic
func boot(cmdStr string) {
	var (
		err error
		cmd = exec.Command(cmdStr)
	)

	// Print uname
	fmt.Println(genUname())

	// Connect files, empty environment
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = make([]string, 0)

	// Start process
	if err = cmd.Run(); err != nil {
		panic(err)
	}
}

// Change root and directory
func switchRoot(path string) error {
	var err error

	if err = unix.Chroot(path); err != nil {
		return err
	}
	if err = unix.Chdir(default_root); err != nil {
		return err
	}
	return nil
}

// Print usage
func usage() {
	fmt.Println("rootctl -r PATH [-c CMD]")
	flag.PrintDefaults()
}

func main() {
	var (
		err      error
		rootFlag = flag.String("r", "", "Specify path to root")
		cmdFlag  = flag.String("c", default_command, "Specify command to execute")
	)

	flag.Usage = usage
	flag.Parse()

	// Switch root, or crash
	if err = switchRoot(*rootFlag); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	boot(*cmdFlag)
	os.Exit(0)
}
