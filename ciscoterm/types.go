package ciscoterm

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type CiscoDevice struct {
	Hostname     string
	Username     string
	Password     string
	Enable       string
	KeyExchanges []string
	Timeout      int32
}

type terminal struct {
	stdinBuf  io.WriteCloser
	stdoutBuf io.Reader
	session   *ssh.Session
	cmdPrompt string
	isEnabled bool
}
