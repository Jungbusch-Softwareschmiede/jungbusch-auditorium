// +build windows

package permissions

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/pkg/errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Überprüft die Read/Write-Rechte des ausführenden Users einer Datei oder Directory
func Permission(in string) (read bool, write bool, isDir bool, err error) {
	// Pfad bereinigen
	in = filepath.Clean(filepath.ToSlash(in))
	if !filepath.IsAbs(in) {
		in, err = filepath.Abs(filepath.Dir(os.Args[0]) + static.PATH_SEPERATOR + filepath.Clean(in))
		if err != nil {
			return
		}
	}

	// Herausfinden, ob der übergebene Pfad eine Directory ist
	info, err := os.Stat(in)
	if err != nil {
		return
	}
	isDir = info.IsDir()

	if isDir {
		read, write, err = getDirectoryPerms(in)
	} else {
		read, write, err = getFilePerms(in)
	}
	return
}

// Überprüft die Permissions des ausführenden Users einer File
func getFilePerms(file string) (read bool, write bool, err error) {
	// Rechte der Directory getten
	var r bool
	r, err = directoryHasRead(filepath.Dir(file))

	// Wenn wir keine read-Rechte auf der Directory haben, dann können wir die File weder lesen, noch schreiben
	if err != nil || !r {
		return false, false, nil
	}

	// Wenn die Rechte der Directory stimmen, checken wir nun die Rechte der File
	read, err = fileHasRead(file)
	write, err = fileHasWrite(file)
	return
}

// Überprüft die Permissions des ausführenden Users einer Directory
// Der User kann read-, aber keine write-Rechte haben - oder andersrum
func getDirectoryPerms(dir string) (read bool, write bool, err error) {
	// Pfad bereinigen
	in := filepath.Clean(filepath.ToSlash(dir + static.PATH_SEPERATOR + "permissiontest"))
	if !filepath.IsAbs(in) {
		in, err = filepath.Abs(filepath.Dir(os.Args[0]) + static.PATH_SEPERATOR + filepath.Clean(in))
		if err != nil {
			return
		}
	}

	// Read-Permissions überprüfen
	read, err = directoryHasRead(dir)
	if err != nil {
		return
	}

	// Read/Write-Permissions testen
	write, err = directoryHasWrite(in)
	return
}

func fileHasWrite(in string) (read bool, err error) {
	return fileHelper(in, false, os.O_WRONLY)
}

func fileHasRead(in string) (read bool, err error) {
	return fileHelper(in, false, os.O_RDONLY)
}

// Überprüft ob der aktuelle User read-Berechtigungen für die angegebene Directory hat, indem wir die Directory selbst auslesen
func directoryHasRead(in string) (read bool, err error) {
	_, err = ioutil.ReadDir(in)
	if err != nil {
		if !errors.Is(err, fs.ErrPermission) {
			return false, errors.New("Unbekannter Fehler: " + err.Error())
		} else {
			return false, nil
		}
	}

	return true, nil
}

// Überprüft, ob eine Directory Write-Permissions hat. Es ist möglich dass wir write haben, aber keine read
func directoryHasWrite(in string) (write bool, err error) {
	return fileHelper(in, true, os.O_RDWR|os.O_CREATE)
}

func fileHelper(in string, remove bool, args int) (perm bool, err error) {
	var file *os.File
	file, err = os.OpenFile(in, args, 0666)

	if err != nil {
		// Wenn der err nicht vom Typ "access denied" ist, ist was blödes passiert
		if !errors.Is(err, fs.ErrPermission) {
			// Wenn Access nicht denied wurde, ist ein unbekannter Fehler aufgetreten.
			err = errors.New("Unbekannter Fehler: " + err.Error())
		} else {
			return false, nil
		}
	} else {
		perm = true
		file.Close()
		if remove {
			err = os.Remove(in)
		}
	}

	return
}
