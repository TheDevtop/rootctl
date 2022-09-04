package main

/*
	Prog: Launch chroot environments
	Vers: 0.5
	Auth: Thijs Haker
*/

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

type entry struct {
	Path string
	Cmd  string
	Args []string
}

const (
	EC_ERR    = 1
	EC_DEF    = 0
	CONF_FILE = "/etc/rootctl.conf"
	ROOT_PATH = "/"
)

// Chroot and chdir
func switchRoot(Path string) error {
	var err error

	if err = unix.Chroot(Path); err != nil {
		return err
	}
	if err = unix.Chdir(ROOT_PATH); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		err       error
		name      string
		buf       []byte
		ok        bool
		cmd       *exec.Cmd
		rootEntry = new(entry)
		confMap   = make(map[string]entry, 2)
	)

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: rootctl [name]")
		os.Exit(EC_ERR)
	}
	name = os.Args[1]

	// Read config file
	if buf, err = os.ReadFile(CONF_FILE); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EC_ERR)
	}

	// Unmarshal to map with configuration entries
	if err = json.Unmarshal(buf, &confMap); err != nil {
		fmt.Fprintln(os.Stderr, "Error: Can't parse file contents!")
		os.Exit(EC_ERR)
	}

	// Check if specified entry is in map
	if *rootEntry, ok = confMap[name]; !ok {
		fmt.Fprintf(os.Stderr, "Error: Entry %s could not be resolved!", name)
		os.Exit(EC_ERR)
	}

	if err = switchRoot(rootEntry.Path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EC_ERR)
	}

	// Create new process and clean environtment
	cmd = exec.Command(rootEntry.Cmd, rootEntry.Args...)
	cmd.Env = nil
	cmd.Run()
	os.Exit(EC_DEF)
}
