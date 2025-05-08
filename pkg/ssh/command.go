package ssh

import "os/exec"

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

func (h Host) Command(cmd string, sudo bool) (stdout, stderr string, e error) {
	return "", "", nil
}
