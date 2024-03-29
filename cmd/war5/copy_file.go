package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rdtharri/wargames/sshtools"
	"golang.org/x/crypto/ssh"
)

func main() {

	// Grab Previous Password
	pass, err := ioutil.ReadFile("cmd/war4/bandit4_pass")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Password file open errpr: %s\n", err)
		os.Exit(1)
	}

	sshConfig := &ssh.ClientConfig{
		User: "bandit4",
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

	// Setup Command
	cmd := &sshtools.SSHCommand{
		Path: "find inhere/* -type f -exec  file {} + | grep ASCII | awk -F: '{print $1}'",
		//Env:    []string{"THIS=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)

	filePath, err := client.RunCommandGetOutput(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Command Complete. \n")

	// Grab File
	scpConfig := &sshtools.SCPConfig{
		SourcePath: strings.TrimRight(filePath, "\r\n"),
		DestPath:   "cmd/war5/bandit5_pass",
	}

	fmt.Printf("Moving file from: %s to %s\n", scpConfig.SourcePath, scpConfig.DestPath)
	if err := client.GrabFile(scpConfig); err != nil {
		fmt.Fprintf(os.Stderr, "file transfer error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Transfer Complete.")

}
