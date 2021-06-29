// +build linux darwin

package util

import (
	"errors"
	"os/exec"
	s "strings"
)

// Dummy Funktion, wird benötigt da diese Funktion im OS-Detector ausschließlich für Windows aufgerufen wird und Go meckert wenn sie für Unix nicht vorhanden ist.
func RegQuery(keyPath string, value string) (string, error) {
	return "", nil
}

// Führt den übergebenen Befehl aus
func ExecCommand(command string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil && !s.Contains(command, ">/dev/null") {
		return "", errors.New(err.Error() + ": " + string(out))
	}
	return string(out), nil
}
