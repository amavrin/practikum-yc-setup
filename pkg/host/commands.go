package host

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/amavrin/practikum-yc-setup/pkg/browser"
	"github.com/amavrin/practikum-yc-setup/pkg/cons"
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
		log.Printf("error installing yc: %s\n", out)
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

	cons.Log("sending ", command, " command...")
	err = sshHost.Send(command)
	if err != nil {
		return err
	}
	lines := 10
	cons.Log(`Waiting for response...`)
	_, err = sshHost.WaitFor("After your successful authentication", lines)
	if err != nil {
		return err
	}
	cons.Log("sending Enter...")
	err = sshHost.Send("")
	if err != nil {
		return err
	}
	cons.Log("waiting for URL...")
	urlString, err := sshHost.WaitFor("https://auth.yandex.cloud", lines)
	if err != nil {
		return err
	}
	port, url, err := getRedirectPort(urlString)
	if err != nil {
		return err
	}

	cons.Log("setting up SSH tunnel...")
	_, err = sshHost.Tunnel(port, port)
	if err != nil {
		return err
	}

	cons.Log("opening browser...")
	err = browser.Open(url)
	if err != nil {
		return err
	}

	cons.Log("configuring profile settings...")
	_, err = sshHost.WaitFor("Please choose folder to use", 20)
	if err != nil {
		return err
	}

	cons.Log("using 1-st available folder...")
	err = sshHost.Send("1")
	if err != nil {
		return err
	}

	_, err = sshHost.WaitFor("Your current folder has been set", 20)
	if err != nil {
		return err
	}
	cons.Log("not choosing default zone...")
	err = sshHost.Send("n")
	if err != nil {
		return err
	}

	cons.Log("waiting for prompt...")
	_, err = sshHost.GetPrompt()
	if err != nil {
		return err
	}
	return nil
}

func CheckYC(sshHost *ssh.Host) (string, error) {
	return sshHost.Command("/home/student/yandex-cloud/bin/yc resource-manager cloud list", false)
}

func getRedirectPort(input string) (int, string, error) {
	urlStart := strings.Index(input, "https://auth.yandex.cloud/oauth")
	if urlStart == -1 {
		return 0, "", errors.New("url not found in the input")
	}
	urlStr := strings.Fields(input[urlStart:])[0]
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", err
	}
	queryParams := parsedURL.Query()
	redirectURIEncoded := queryParams.Get("redirect_uri")
	if redirectURIEncoded == "" {
		return 0, "", errors.New("redirect_uri not found in URL")
	}

	redirectURI, err := url.QueryUnescape(redirectURIEncoded)
	if err != nil {
		return 0, "", err
	}
	parsedRedirectURI, err := url.Parse(redirectURI)
	if err != nil {
		return 0, "", err
	}

	redirectHost := parsedRedirectURI.Host
	if redirectHost == "" {
		return 0, "", errors.New("redirect host not found in redirect uri")
	}
	parts := strings.Split(redirectHost, ":")
	if len(parts) != 2 {
		return 0, "", errors.New("port not found in redirect host string")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, "", err
	}
	return port, urlStr, nil
}
