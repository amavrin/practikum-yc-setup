package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/amavrin/practikum-yc-setup/pkg/host"
	"github.com/amavrin/practikum-yc-setup/pkg/ssh"
)

func run(server string) error {
	log.Printf("running for %s", server)
	sshHost, err := ssh.New(server)
	if err != nil {
		return err
	}

	log.Println("setting up xdg-open...")
	err = host.SetupXdgOpen(&sshHost)
	if err != nil {
		return fmt.Errorf("host.SetupXdgOpen: %w", err)
	}
	log.Println("removing yc profile...")
	err = host.RemoveYCProfile(&sshHost)
	if err != nil {
		return fmt.Errorf("host.RemoveYCProfile: %w", err)
	}
	log.Println("installing yc...")
	err = host.InstallYC(&sshHost)
	if err != nil {
		return fmt.Errorf("host.InstallYC: %w", err)
	}
	return nil
}

func main() {
	var sshServer string
	flag.StringVar(&sshServer, "server", "example.com", "ssh server to connect to")
	flag.Parse()
	err := run(sshServer)
	if err != nil {
		log.Fatal(err)
	}
}
