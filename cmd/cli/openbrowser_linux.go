// go:build linux

package cli

import "os/exec"

// openBrowser tries to open the default web browser at the specified location
func OpenBrowser(location string, exitChan chan bool) error {
	defer func() { exitChan <- true }()
	cmd := exec.Command("xdg-open", location)
	return cmd.Start()
}
