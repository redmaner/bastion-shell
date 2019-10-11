package main

import (
	"log"
	"os"
)

var (
	globalUser = ""
)

func main() {

	// Get SSH connections file
	sshConnectionsFile, ok := os.LookupEnv("BASTION_SSH_CONNECTIONS")
	if !ok {
		log.Fatalf("Please set the BASTION_SSH_CONNECTIONS environment variable")
	}

	// get SSH global user
	sshGlobalUser, _ := os.LookupEnv("BASTION_GLOBAL_USER")
	globalUser = sshGlobalUser

	sshConnections, err := getSSHConnections(sshConnectionsFile)
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = sshConnections.selectConn()
		if err != nil {
			log.Fatal(err)
		}
	}
}
