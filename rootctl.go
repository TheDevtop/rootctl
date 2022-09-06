package main

/*
	Prog: Launch chroot environments
	Vers: 1.1
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
	Path string   // Path to new root
	Cmd  string   // Path to command
	Args []string // Argument list
	Env  []string // Environment list
}

const (
	EC_ERR    = 1 // Exit with errors
	EC_DEF    = 0 // Exit normal
	CONF_FILE = "/etc/rootctl.conf"
	ROOT_PATH = "/"
)

// Generate utsname
func uname() string {
	var buf = new(unix.Utsname)
	unix.Uname(buf)
	return fmt.Sprintf("%s %s %s", string((*buf).Sysname[:]), string(((*buf).Release[:])), string((*buf).Machine[:]))
}

// Chroot and chdir
func switchRoot(path string) error {
	var err error

	if err = unix.Chroot(path); err != nil {
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
		fmt.Fprintf(os.Stderr, "Error: Entry %s couldn't be resolved!\n", name)
		os.Exit(EC_ERR)
	}

	if err = switchRoot(rootEntry.Path); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(EC_ERR)
	}

	// Print utsname
	fmt.Fprintln(os.Stdout, uname())

	// Create new process and change properties
	cmd = exec.Command(rootEntry.Cmd, rootEntry.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = rootEntry.Env

	// Launch process
	if err = cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(EC_DEF)
}
