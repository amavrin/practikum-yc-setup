package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/amavrin/practikum-yc-setup/pkg/cons"
	"github.com/amavrin/practikum-yc-setup/pkg/host"
	"github.com/amavrin/practikum-yc-setup/pkg/ssh"
)

func run(server string, federationID string) error {
	cons.Log("running for ", server)
	sshHost, err := ssh.New(server)
	if err != nil {
		return err
	}

	cons.Log("setting up xdg-open...")
	err = host.SetupXdgOpen(&sshHost)
	if err != nil {
		return fmt.Errorf("host.SetupXdgOpen: %w", err)
	}
	cons.Log("removing yc profile...")
	err = host.RemoveYCProfile(&sshHost)
	if err != nil {
		return fmt.Errorf("host.RemoveYCProfile: %w", err)
	}
	cons.Log("installing yc...")
	err = host.InstallYC(&sshHost)
	if err != nil {
		return fmt.Errorf("host.InstallYC: %w", err)
	}
	cons.Log("configuring YC profile...")
	err = host.ConfigureYCProfile(&sshHost, federationID)
	if err != nil {
		return fmt.Errorf("configure YC profile: %w", err)
	}

	out, err := host.CheckYC(&sshHost)
	cons.Log("check YC working...")
	if err != nil {
		return fmt.Errorf("host.CheckYC: %w", err)
	}

	cons.Log("\r\n", out)
	cons.Log("successfully configured YC profile")
	return nil
}

func main() {
	var sshServer string
	var federationID string
	flag.StringVar(&sshServer, "server", "", "ssh server to connect to")
	flag.StringVar(&federationID, "federation-id", "", "federation id")
	flag.Parse()
	if sshServer == "" {
		log.Fatal("SSH server is not set")
	}
	if federationID == "" {
		log.Fatal("federation ID is not set")
	}
	err := run(sshServer, federationID)
	if err != nil {
		log.Fatal(err)
	}
}
