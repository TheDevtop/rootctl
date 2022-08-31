package vkern

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"golang.org/x/sys/unix"
)

// Setup signal handling
func signal_init() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch)

	for {
		<-sigch
	}
}

// Returns utsname
func uname() string {
	var buf = new(unix.Utsname)
	unix.Uname(buf)
	return string((*buf).Version[:]) + string((*buf).Machine[:])
}

// Initialize environment and create child process
func Main(conf *KernConf) {
	var err error

	// Print uname and setup signals
	fmt.Println(uname())
	go signal_init()

	// Set root user
	if err = unix.Setuid(conf.Uid); err != nil {
		panic(err)
	}

	// Create new process
	cmd := exec.Command(conf.CmdStr, conf.ArgStr...)

	// Configure new process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if conf.ClearEnv {
		cmd.Env = nil
	}

	// Execute new process
	if err = cmd.Start(); err != nil {
		fmt.Println(err)
	}
}
