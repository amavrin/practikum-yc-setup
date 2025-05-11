package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/amavrin/practikum-yc-setup/pkg/cons"
)

type Host struct {
	Address string
	SSHPath string
	Prompt  byte
	StdIn   io.WriteCloser
	StdOut  io.ReadCloser
	StdErr  io.ReadCloser
	Scanner *bufio.Scanner
	Reader  *bufio.Reader
	Cmd     *exec.Cmd
	Debug   bool
}

func New(address string) (Host, error) {
	path, err := exec.LookPath("ssh")
	if err != nil {
		return Host{}, err
	}
	return Host{
		Address: address,
		SSHPath: path,
		Prompt:  '$',
	}, nil
}

func (h *Host) SetPrompt(p byte) {
	h.Prompt = p
}

func (h *Host) SetDebug(d bool) {
	h.Debug = d
}

func (h *Host) Command(c string, sudo bool) (string, error) {
	if c == "" {
		return "", errors.New("cmd should not be empty")
	}
	cmdLine := strings.Fields(c)
	args := []string{"-tt", h.Address}
	if sudo {
		args = append(args, "sudo")
	}
	args = append(args, cmdLine...)
	cmd := exec.Command(h.SSHPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	out := strings.TrimSpace(string(output))
	return out, nil
}

func (h *Host) Tunnel(local, remote int) (func(), error) {
	portSpec := fmt.Sprintf("%d:localhost:%d", local, remote)
	args := []string{"-L", portSpec, h.Address, "sleep", "infinity"}
	cmd := exec.Command(h.SSHPath, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	err = checkTCPConn(local)
	if err != nil {
		return nil, err
	}
	return func() { cmd.Process.Kill() }, nil
}

func (h *Host) StartChat() error {
	args := []string{"-tt", h.Address}
	cmd := exec.Command(h.SSHPath, args...)
	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	h.StdIn = stdIn
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	h.StdOut = stdOut
	h.Scanner = bufio.NewScanner(h.StdOut)
	h.Reader = bufio.NewReader(h.StdOut)
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	h.StdErr = stdErr
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (h *Host) WaitFor(wait string, maxLines int) (string, error) {
	ret := ""
	for h.Scanner.Scan() {
		line := h.Scanner.Text()
		if h.Debug {
			cons.Log("get: ", line)
		}
		ret += line + "\n"
		maxLines--
		if maxLines < 0 {
			return "", errors.New("max lines reached with no match")
		}
		if !strings.Contains(ret, wait) {
			continue
		}
		return ret, nil
	}
	return ret, errors.New("no more input and not match")
}

func (h *Host) Send(what string) error {
	_, err := io.WriteString(h.StdIn, what+"\n")
	return err
}

func (h *Host) GetPrompt() (string, error) {
	return h.Reader.ReadString(h.Prompt)
}

func checkTCPConn(port int) error {
	timeout := 1 * time.Second
	count := 10
	addr := fmt.Sprintf("localhost:%d", port)
	for range count {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			time.Sleep(timeout)
			continue
		}
		conn.Close()
		return nil
	}
	return fmt.Errorf("connection to %s not established after %d of %v", addr, count, timeout)
}
