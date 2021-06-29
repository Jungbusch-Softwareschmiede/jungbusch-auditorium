// +build windows

package privilege

import "os"

// Prüfen ob Prgramm mit elevated privileges gestartet ist (Ausführen als Administrator)
func HasRootPrivileges() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}
