package ciscoterm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func (t *terminal) Connect(ciscoDev CiscoDevice) error {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	config := &ssh.ClientConfig{
		Config: ssh.Config{
			KeyExchanges: ciscoDev.KeyExchanges,
		},
		User: ciscoDev.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(ciscoDev.Password),
		},
		Timeout:         time.Second * time.Duration(ciscoDev.Timeout),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For development, use ssh.FixedHostKey or ssh.KnownHosts for production
	}

	conn, err := ssh.Dial("tcp", ciscoDev.Hostname, config)
	if err != nil {
		return fmt.Errorf("connect dial error: %w", err)
	}
	t.session, err = conn.NewSession()
	if err != nil {
		return fmt.Errorf("connect session error: %w", err)
	}

	err = t.session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return fmt.Errorf("connect pty error: %w", err)
	}

	t.stdoutBuf, err = t.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("connect stdout pipe error: %w", err)
	}

	t.stdinBuf, err = t.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("connect stdin pipe error: %w", err)
	}

	err = t.session.Shell()
	if err != nil {
		return fmt.Errorf("connect shell error: %w", err)
	}

	var byteCount int
	stdOutBytes := make([]byte, 100000)
	for {
		byteCount, err = t.stdoutBuf.Read(stdOutBytes)
		if err != nil {
			break
		}
		if byteCount <= 0 {
			break
		}
		s := string(stdOutBytes[:byteCount])
		lines := strings.Split(s, "\r\n")
		for _, line := range lines {
			if strings.HasSuffix(strings.TrimSpace(line), ">") {
				t.cmdPrompt = strings.TrimRight(strings.TrimSpace(line), ">")
				break
			}
		}
		if t.cmdPrompt != "" {
			break
		}
	}

	return err
}

func (t *terminal) Close() error {
	if t.session != nil {
		return t.session.Close()
	}
	return errors.New("no session")
}

func (t *terminal) EnableTerm(ctx context.Context, enablePasswd string) error {
	_, err := t.stdinBuf.Write([]byte("enable\n" + enablePasswd + "\n"))
	if err != nil {
		return fmt.Errorf("enable terminal error: %w", err)
	}

	t.isEnabled = true
	isFinished, _, err := t.readCommandOutput(ctx)
	if err != nil {
		return fmt.Errorf("enable terminal error: %w", err)
	}
	if !isFinished {
		return errors.New("enable terminal error: no command prompt received")
	}
	return err
}

func (t *terminal) DisablePagination(ctx context.Context) error {
	_, err := t.stdinBuf.Write([]byte("terminal pager 0\n"))
	if err != nil {
		return fmt.Errorf("disable pagination error: %w", err)
	}
	isFinished, _, err := t.readCommandOutput(ctx)
	if err != nil {
		return fmt.Errorf("disable pagination error: %w", err)
	}
	if !isFinished {
		return errors.New("disable pagination: no command prompt received")
	}
	return err
}

func (t *terminal) ExecuteCommand(ctx context.Context, cmd string) ([]string, error) {
	_, err := t.stdinBuf.Write([]byte(cmd + "\n"))
	if err != nil {
		return nil, fmt.Errorf("execute command error: %w", err)
	}

	isFinished, output, err := t.readCommandOutput(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute command error: %w", err)
	}
	if !isFinished {
		return nil, errors.New("execute command error: no command prompt received")
	}
	return output, err
}

func (t *terminal) readCommandOutput(ctx context.Context) (bool, []string, error) {
	time.Sleep(time.Millisecond * 100)
	stdOutBytes := make([]byte, 100000)
	var byteCount int
	var err error
	var lines []string
	var cmdFinished bool
	var s string

	for {
		select {
		case <-ctx.Done():
			return cmdFinished, lines, errors.New("exceeded timeout")
		default:
			byteCount, err = t.stdoutBuf.Read(stdOutBytes)
			if err != nil {
				return cmdFinished, lines[1:(len(lines) - 1)], err
			}
			if byteCount <= 0 {
				return cmdFinished, lines[1:(len(lines) - 1)], err
			}
			s = string(stdOutBytes[:byteCount])
			lines = append(lines, strings.Split(s, "\r\n")...)
			if t.isEnabled && strings.TrimSpace(lines[len(lines)-1]) == t.cmdPrompt+"#" {
				cmdFinished = true
				return cmdFinished, lines[1:(len(lines) - 1)], err
			}
			if !t.isEnabled && strings.TrimSpace(lines[len(lines)-1]) == t.cmdPrompt+">" {
				cmdFinished = true
				return cmdFinished, lines[1:(len(lines) - 1)], err
			}
			if t.isEnabled && strings.TrimSpace(lines[len(lines)-1]) == "Password:" {
				cmdFinished = false
				return cmdFinished, lines[1:(len(lines) - 1)], err
			}
		}
	}
}
