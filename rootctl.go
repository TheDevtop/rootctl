package main

/*
	Prog: Manage chroot environments
	Vers: 0.3
	Auth: Thijs Haker
*/

import (
	"flag"
	"fmt"
	"os"

	"github.com/TheDevtop/rootctl/vkern"
	"golang.org/x/sys/unix"
)

const (
	default_command = "/usr/bin/login"
	default_root    = "/"
)

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

	// Allocate and fill kconf
	kconf := new(vkern.KernConf)
	kconf.CmdStr = *cmdFlag

	vkern.Main(kconf)
	os.Exit(0)
}
