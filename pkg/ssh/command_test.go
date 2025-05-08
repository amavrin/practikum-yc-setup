package ssh

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

var TestSSHServer string

func TestMain(m *testing.M) {
	flag.StringVar(&TestSSHServer, "server", "localhost", "user@server to use in test")
	flag.Parse()

	code := m.Run()

	os.Exit(code)
}

func TestNew(t *testing.T) {
	h, err := New("")
	require.NoError(t, err)
	sshPath := whichSSH(t)
	require.NotEmpty(t, sshPath)
	require.Equal(t, h.SSHPath, sshPath)
}

func TestSSHCommand(t *testing.T) {
	h, err := New(TestSSHServer)
	require.NoError(t, err)

	out, err := h.Command("id -u", false)
	require.NoError(t, err)
	assert.Equal(t, out, "1001")

	out, err = h.Command("id -u", true)
	require.NoError(t, err)
	assert.Equal(t, out, "0")
}

func TestSSHTunnel(t *testing.T) {
	h, err := New(TestSSHServer)
	disconnect, err := h.Tunnel(2022, 22)
	defer disconnect()
	require.NoError(t, err)
	conn, err := net.Dial("tcp", "localhost:2022")
	require.NoError(t, err)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	resp, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, resp, "OpenSSH")
}

func TestSSHChat(t *testing.T) {
	h, err := New(TestSSHServer)
	require.NoError(t, err)
	err = h.StartChat()
	require.NoError(t, err)
	err = h.Send("id -u")
	require.NoError(t, err)
	out, err := h.WaitFor("1001", 20)
	require.NoError(t, err)
	fmt.Println(out)
}
