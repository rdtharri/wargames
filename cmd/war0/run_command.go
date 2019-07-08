package main

import (
	"fmt"
	"os"

	"github.com/rdtharri/wargames/sshtools"
	"golang.org/x/crypto/ssh"
)

func main() {

	sshConfig := &ssh.ClientConfig{
		User: "bandit0",
		Auth: []ssh.AuthMethod{
			ssh.Password("bandit0"),
			//	sshtools.SSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client := &sshtools.SSHClient{
		Config: sshConfig,
		Host:   "bandit.labs.overthewire.org",
		Port:   2220,
	}

	cmd := &sshtools.SSHCommand{
		Path: "ls -al",
		//Env:    []string{"THIS=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Command Complete.")

}
