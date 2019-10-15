package main

import (
	"fmt"
	"os"
	"os/exec"
)

// connectSSH is function that connects to the desired host, taking the name of
// the host and the host information as input. If the connection fails, the function
// will return an error. This function depends on the ssh binary to work.
func connectSSH(name string, info sshHostInfo) error {

	fmt.Printf("\n>>> Connecting to %s\n", name)

	// Search for the ssh binary in the $PATH variable. If the ssh binary cannot be found
	// an error will be returned.
	if _, err := exec.LookPath("ssh"); err != nil {
		return fmt.Errorf("The ssh binary could not be found, make sure ssh is installed and the $PATH variable is set appropriately")
	}

	// Prepare execution of ssh
	cmdArgs := fmt.Sprintf("%s@%s", info.Username, info.Host)
	cmd := exec.Command("ssh", cmdArgs, "-p", info.Port)

	// Tunnel stdin, stdout and stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the ssh connection
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for the ssh connection to complete
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
