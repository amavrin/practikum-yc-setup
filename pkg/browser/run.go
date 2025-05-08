package browser

import (
	"errors"
	"os/exec"
	"runtime"
)

func Open(url string) error {
	var path string
	var err error
	switch runtime.GOOS {
	case "darwin":
		path, err = exec.LookPath("open")
		if err != nil {
			return err
		}
	case "linux":
		path, err = exec.LookPath("xdg-open")
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported OS")
	}

	cmd := exec.Command(path, url)
	return cmd.Run()
}
