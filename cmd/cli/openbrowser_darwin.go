// go:build darwin

package cli

import "os/exec"

// openBrowser tries to open the default web browser at the specified location
func OpenBrowser(location string) error {
	return exec.Command("open", location).Start()
}
