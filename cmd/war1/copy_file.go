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

	scpConfig := &sshtools.SCPConfig{
		SourcePath: "readme",
		DestPath:   "bandit0_pass",
	}

	fmt.Printf("Moving file from: %s to %s\n", scpConfig.SourcePath, scpConfig.DestPath)
	if err := client.GrabFile(scpConfig); err != nil {
		fmt.Fprintf(os.Stderr, "file transfer error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Transfer Complete.")

}
