// +build linux darwin

// Erlaubt das bestimmen von Datei- oder Directory-Permissions basierend auf den Berechtigungen des aktuellen Prozesses
package permissions

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

// Überprüft die Read/Write-Rechte des ausführenden Users einer Datei oder Directory
func Permission(in string) (read bool, write bool, isDir bool, err error) {
	// Pfad bereinigen
	path := filepath.Clean(filepath.ToSlash(in))
	if !filepath.IsAbs(in) {
		in, err = filepath.Abs(filepath.Dir(os.Args[0]) + static.PATH_SEPERATOR + filepath.Clean(in))
		if err != nil {
			return
		}
	}

	// Prüfen, ob Datei/Verzeichnis vorhanden ist
	var info os.FileInfo
	if info, err = os.Stat(path); err != nil {
		return
	} else {
		if info == nil {
			err = errors.New("Unbekannter Fehler!")
			return
		} else {
			// Prüfen, ob es ein Verzeichnis oder eine Datei ist
			isDir = info.IsDir()
		}

		// Prüfen, ob Lesen möglich ist
		if err = unix.Access(path, unix.R_OK); err == nil {
			read = true
		}
		// Prüfen, ob Schreiben möglich ist
		if err = unix.Access(path, unix.W_OK); err == nil {
			write = true
		}
	}

	return
}
