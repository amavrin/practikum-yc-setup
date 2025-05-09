package host

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

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

func ConfigureYCProfile(sshHost *ssh.Host, fedID string) error {
	command := fmt.Sprintf("yc init --federation-id=%s", fedID)
	err := sshHost.StartChat()
	if err != nil {
		return err
	}
	_, err = sshHost.GetPrompt()
	if err != nil {
		return err
	}
	log.Printf("sending '%s' command...", command)
	err = sshHost.Send(command)
	if err != nil {
		return err
	}
	lines := 10
	log.Printf(`Waiting for response...`)
	_, err = sshHost.WaitFor("After your successful authentication", lines)
	if err != nil {
		return err
	}
	log.Print("sending Enter...")
	err = sshHost.Send("")
	if err != nil {
		return err
	}
	log.Print("waiting for URL...")
	urlString, err := sshHost.WaitFor("https://auth.yandex.cloud", lines)
	if err != nil {
		return err
	}
	fmt.Println("got URL: ", urlString)
	return nil
}

func getRedirectPort(input string) (int, error) {
	urlStart := strings.Index(input, "https://auth.yandex.cloud/oauth")
	if urlStart == -1 {
		return 0, errors.New("url not found in the input")
	}
	urlStr := strings.Fields(input[urlStart:])[0]
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, err
	}
	queryParams := parsedURL.Query()
	redirectURIEncoded := queryParams.Get("redirect_uri")
	if redirectURIEncoded == "" {
		return 0, errors.New("redirect_uri not found in URL")
	}

	redirectURI, err := url.QueryUnescape(redirectURIEncoded)
	if err != nil {
		return 0, err
	}
	parsedRedirectURI, err := url.Parse(redirectURI)
	if err != nil {
		return 0, err
	}

	redirectHost := parsedRedirectURI.Host
	if redirectHost == "" {
		return 0, errors.New("redirect host not found in redirect uri")
	}
	parts := strings.Split(redirectHost, ":")
	if len(parts) != 2 {
		return 0, errors.New("port not found in redirect host string")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	return port, nil
}
