package host

import (
	"fmt"
	"log"

	"github.com/amavrin/practikum-yc-setup/pkg/ssh"
)

var xdgFileContent = `#!/bin/sh\necho $* > /dev/tty`

func SetupXdgOpen(sshHost *ssh.Host) error {
	xdgFile := "/usr/local/bin/xdg-open"
	tmpFile := "/tmp/xdg-open"
	command := fmt.Sprintf("echo -e '%s' > %s", xdgFileContent, tmpFile)
	command2 := fmt.Sprintf("cp %s %s", tmpFile, xdgFile)
	command3 := fmt.Sprintf("chmod 755 %s", xdgFile)
	_, err := sshHost.Command(command, false)
	if err != nil {
		return err
	}
	_, err = sshHost.Command(command2, true)
	if err != nil {
		return err
	}
	_, err = sshHost.Command(command3, true)
	if err != nil {
		return err
	}
	return err
}

func RemoveYCProfile(sshHost *ssh.Host) error {
	command := "rm -rf ./.config/yandex-cloud"
	_, err := sshHost.Command(command, false)
	return err
}

func InstallYC(sshHost *ssh.Host) error {
	command := "curl -sSL https://storage.yandexcloud.net/yandexcloud-yc/install.sh | bash"
	out, err := sshHost.Command(command, true)
	if err != nil {
		log.Printf("error installing yc: %s", out)
		return err
	}
	return nil
}
