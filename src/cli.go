package main

import (
	"log"
	"os"
)

var (
	globalUser = ""
)

func main() {

	// Get SSH hostsfile
	sshHostsFile, ok := os.LookupEnv("BASTION_SSH_HOSTS")
	if !ok {
		log.Fatalf("Please set the BASTION_SSH_HOSTS environment variable")
	}

	// get SSH global user
	sshGlobalUser, _ := os.LookupEnv("BASTION_GLOBAL_USER")
	globalUser = sshGlobalUser

	sshConnections, err := getHosts(sshHostsFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = sshConnections.selectHost()
		if err != nil {
			log.Fatal(err)
		}
	}
}
