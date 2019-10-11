package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"

	"github.com/hashicorp/hcl"
)

const bastionWelcome = `
   ____                     __
  /\  _ \                  /\ \__  __
  \ \ \L\ \     __      ___\ \ ,_\/\_\     ___     ___
   \ \  _ <'  /'__ \   /',__\\ \ \/\/\ \  / __ \ /' _  \
    \ \ \L\ \/\ \L\.\_/\__,  \\ \ \_\ \ \/\ \L\ \/\ \/\ \
     \ \____/\ \__/.\_\/\____/ \ \__\\ \_\ \____/\ \_\ \_\
      \/___/  \/__/\/_/\/___/   \/__/ \/_/\/___/  \/_/\/_/

  Please select the desired host to connect:
`

type sshConns struct {
	SSHConnections map[string]sshConnInfo
	ConnNames      []string
}

type sshConnInfo struct {
	Username string
	Host     string
	Port     string
}

func getSSHConnections(path string) (*sshConns, error) {

	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conns sshConns
	err = hcl.Unmarshal(confFile, &conns)
	if err != nil {
		return nil, err
	}

	for name := range conns.SSHConnections {
		conns.ConnNames = append(conns.ConnNames, name)
	}
	sort.Strings(conns.ConnNames)

	return &conns, nil
}

func (sc *sshConns) selectConn() error {

	clearScreen()
	fmt.Print(bastionWelcome)

	fmt.Printf("    0 - Close bastion\n\n")
	for index, name := range sc.ConnNames {
		fmt.Printf("    %d - %s\n", index+1, name)
	}

	input := getConnInput(len(sc.SSHConnections))

	if input == 0 {
		os.Exit(0)
	}

	connName := sc.ConnNames[input-1]
	connInfo := sc.SSHConnections[connName]

	if globalUser != "" {
		connInfo.Username = globalUser
	}

	err := connectSSH(connName, connInfo)
	if err != nil {
		return fmt.Errorf("Failed to connect to %s: %v", connName, err)
	}

	return nil
}

func connectSSH(name string, info sshConnInfo) error {

	fmt.Printf("Connecting to %s\n", name)

	cmdArgs := fmt.Sprintf("%s@%s", info.Username, info.Host)

	cmd := exec.Command("ssh", cmdArgs, "-p", info.Port)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func getConnInput(max int) int {

	var input int
	fmt.Printf("\n  Your choice [1 - %d]: ", max)
	fmt.Scan(&input)

	if input < 0 || input > max {
		fmt.Printf("\nWrong input, please try again")
		return getConnInput(max)
	}
	return input
}
