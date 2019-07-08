package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rdtharri/wargames/sshtools"
	"golang.org/x/crypto/ssh"
)

func main() {

	// Grab Previous Password
	pass, err := ioutil.ReadFile("cmd/war1/bandit1_pass")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Password file open errpr: %s\n", err)
		os.Exit(1)
	}

	sshConfig := &ssh.ClientConfig{
		User: "bandit1",
		Auth: []ssh.AuthMethod{
			ssh.Password(string(pass[:(len(pass) - 1)])),
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
		SourcePath: "./-",
		DestPath:   "cmd/war2/bandit2_pass",
	}

	fmt.Printf("Moving file from: %s to %s\n", scpConfig.SourcePath, scpConfig.DestPath)
	if err := client.GrabFile(scpConfig); err != nil {
		fmt.Fprintf(os.Stderr, "file transfer error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Transfer Complete.")

}
