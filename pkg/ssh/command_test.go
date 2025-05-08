package ssh

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func whichSSH(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("which", "ssh")
	output, err := cmd.Output()
	require.NoError(t, err)
	p := strings.TrimSpace(string(output))
	return p
}

func TestNew(t *testing.T) {
	h, err := New("")
	require.NoError(t, err)
	sshPath := whichSSH(t)
	require.NotEmpty(t, sshPath)
	require.Equal(t, h.SSHPath, sshPath)
}
