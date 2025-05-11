package browser

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

func Open(url string) error {
	var path string
	var err error
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		path, err = exec.LookPath("open")
		if err != nil {
			return err
		}
		cmd = exec.Command(path, url)
	case "linux":
		path, err = exec.LookPath("xdg-open")
		if err != nil {
			return err
		}
		cmd = exec.Command(path, url)
	case "windows":
		winURL := strings.ReplaceAll(url, "&", "^&")
		cmd = exec.Command("cmd", "/c", "start", "", winURL)
	default:
		return errors.New("unsupported OS")
	}

	return cmd.Run()
}
