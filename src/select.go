package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/hashicorp/hcl"
)

// sshHosts is a type holding the hosts that can be jumped to
type sshHosts struct {
	Hosts     map[string]sshHostInfo `hcl:"ssh_hosts"`
	hostNames []string
}

// sshHostInfo holds the information required to connect to the host over SSH
type sshHostInfo struct {
	Username string `hcl:"username"`
	Host     string `hcl:"host"`
	Port     string `hcl:"port"`
}

// getHosts retrieves the hosts that can be jumped to using the hostsfile.
func getHosts(path string) (*sshHosts, error) {

	// Read the configuration file
	confFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal file, expecting HCL format
	var hosts sshHosts
	err = hcl.Unmarshal(confFile, &hosts)
	if err != nil {
		return nil, err
	}

	// Read all hostnames and sort in alphabetical order
	for name := range hosts.Hosts {
		hosts.hostNames = append(hosts.hostNames, name)
	}
	sort.Strings(hosts.hostNames)

	return &hosts, nil
}

// selectHost is a functio to select the desired host from the list of hosts, which
// are extracted from the hostsfile
func (sh *sshHosts) selectHost() error {

	// Clear scree
	clearScreen()
	fmt.Print(bastionWelcome)

	// Print the connection options
	fmt.Printf("    0 - Close bastion\n\n")
	for index, name := range sh.hostNames {
		fmt.Printf("    %d - %s\n", index+1, name)
	}

	// Read user input
	input := getHostInput(len(sh.Hosts))

	if input == 0 {
		os.Exit(0)
	}

	// Retrieve information of the host from the selected input
	hostName := sh.hostNames[input-1]
	hostInfo := sh.Hosts[hostName]

	if globalUser != "" {
		hostInfo.Username = globalUser
	}

	// Connect to the selected host with the host information
	err := connectSSH(hostName, hostInfo)
	if err != nil {
		return fmt.Errorf("Failed to connect to %s: %v", hostName, err)
	}

	return nil
}

// getHostInput is a function that retrieves user input, for host selection
func getHostInput(max int) int {

	var input int
	fmt.Printf("\n  Your choice [1 - %d]: ", max)

	if _, err := fmt.Scan(&input); err != nil {
		return getHostInput(max)
	}

	if input < 0 || input > max {
		fmt.Printf("\nWrong input, please try again")
		return getHostInput(max)
	}
	return input
}
