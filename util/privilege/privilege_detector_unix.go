// +build linux darwin

// Überprüft, ob der auszuführende Nutzer Administrator/ Root-Privilegien hat
package privilege

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"os"
	"syscall"
)

// Prüfen ob Prozess als root gestartet wurde
func HasRootPrivileges() bool {
	if os.Getegid() == 0 {
		syscall.Umask(0)
		static.CREATE_DIRECTORY_PERMISSIONS = 0767
		static.CREATE_FILE_PERMISSIONS = 0646
		return true
	}
	return false
}
