package ssh

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Host struct {
	Address string
	SSHPath string
}

func New(address string) (Host, error) {
	path, err := exec.LookPath("ssh")
	if err != nil {
		return Host{}, err
	}
	return Host{
		Address: address,
		SSHPath: path,
	}, nil
}

func (h Host) Command(c string, sudo bool) (string, error) {
	if c == "" {
		return "", errors.New("cmd should not be empty")
	}
	cmdLine := strings.Fields(c)
	args := []string{"-t", h.Address}
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

func (h Host) Tunnel(local, remote int) (func(), error) {
	portSpec := fmt.Sprintf("%d:localhost:%d", local, remote)
	args := []string{"-L", portSpec, h.Address, "sleep", "infinity"}
	cmd := exec.Command(h.SSHPath, args...)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	return func() { cmd.Process.Kill() }, nil
}
