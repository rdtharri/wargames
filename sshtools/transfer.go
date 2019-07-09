package sshtools

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHCommand struct
type SSHCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// SSHClient struct
type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

// SCPConfig File Transfer struct
type SCPConfig struct {
	SourcePath string
	DestPath   string
}

// GrabFile from Client
func (client *SSHClient) GrabFile(scpConfig *SCPConfig) error {

	var (
		err error
	)

	// Conection
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return err
	}

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(connection)
	if err != nil {
		return err
	}
	defer sftp.Close()

	// Open the source file
	srcFile, err := sftp.Open(scpConfig.SourcePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(scpConfig.DestPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the file
	srcFile.WriteTo(dstFile)

	return err

}

// RunCommand on client
func (client *SSHClient) RunCommand(cmd *SSHCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	// Start a new Session
	if session, err = client.newSession(); err != nil {
		return err
	}
	defer session.Close()

	// Setup standards for command
	if err = client.prepareCommand(session, cmd); err != nil {
		return err
	}

	// Run Command
	err = session.Run(cmd.Path)
	return err
}

// RunCommandGetOutput on client
func (client *SSHClient) RunCommandGetOutput(cmd *SSHCommand) (string, error) {
	var (
		session *ssh.Session
		err     error
		buff    bytes.Buffer
	)

	// Start Pipe
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	// Change Output to Pipe
	old := cmd.Stdout
	cmd.Stdout = w

	// Start a new Session
	if session, err = client.newSession(); err != nil {
		return "", err
	}
	defer session.Close()

	// Setup standards for command
	if err = client.prepareCommand(session, cmd); err != nil {
		return "", err
	}

	// Run Command
	err = session.Run(cmd.Path)

	// Close Pipe Writer, Copy buffer
	w.Close()
	cmd.Stdout = old
	io.Copy(&buff, r)

	return buff.String(), err
}

// Start a new session for the client
func (client *SSHClient) newSession() (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for psuedo terminal failed: %s", err)
	}

	return session, nil
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *SSHCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

// SSHAgent
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
