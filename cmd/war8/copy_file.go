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
	pass, err := ioutil.ReadFile("cmd/war7/bandit7_pass")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Password file open error: %s\n", err)
		os.Exit(1)
	}

	sshConfig := &ssh.ClientConfig{
		User: "bandit7",
		Auth: []ssh.AuthMethod{
			ssh.Password(string(pass)),
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
		Path: "cat data.txt | grep -i 'millionth' | awk '{print $2}'",
		//Env:    []string{"THIS=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)

	passwd, err := client.RunCommandGetOutput(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Command Complete. \n")

	// Write Pass to file
	err = ioutil.WriteFile("cmd/war8/bandit8_pass", []byte(strings.TrimRight(passwd, "\r\n")), 0664)
	if err != nil {
		panic(err)
	}

}
