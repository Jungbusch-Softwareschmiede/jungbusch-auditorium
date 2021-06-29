// +build windows

package util

import (
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"strconv"
	s "strings"
	"syscall"
)

// Führt einen Befehl aus
func ExecCommand(command string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("powershell", "-command", command)
	fmt.Println(cmd)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return "", errors.New(err.Error() + ": " + string(out))
	}
	return string(out), nil
}

// Registry Abfrage ausführen
func RegQuery(keyPath string, value string) (string, error) {
	// Rootkey checken + aus keyPath löschen
	rootKeys := map[string]registry.Key{
		"HKEY_CLASSES_ROOT":     registry.CLASSES_ROOT,
		"HKEY_CURRENT_USER":     registry.CURRENT_USER,
		"HKEY_LOCAL_MACHINE":    registry.LOCAL_MACHINE,
		"HKEY_USERS":            registry.USERS,
		"HKEY_CURRENT_CONFIG":   registry.CURRENT_CONFIG,
		"HKEY_PERFORMANCE_DATA": registry.PERFORMANCE_DATA,
	}

	var rootKey registry.Key
	for rootKeyString, rootKeyConst := range rootKeys {
		if s.Contains(keyPath, rootKeyString) {
			rootKey = rootKeyConst
			keyPath = s.Replace(keyPath, rootKeyString+"\\", "", 1)
		}
	}

	// Key öffnen
	key, err := registry.OpenKey(rootKey, keyPath, registry.QUERY_VALUE)
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return "", static.ERROR_KEY_NOT_FOUND
		} else {
			return "", err
		}
	}
	defer key.Close()

	// Valuetyp ermitteln
	_, valueType, err := key.GetValue(value, []byte{})
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return "", static.ERROR_VALUE_NOT_FOUND
		} else {
			return "", err
		}
	}

	// Value je nach Valuetyp zurückgeben
	switch valueType {
	// REG_SZ + REG_EXPAND_SZ
	case 1, 2:
		res, _, err := key.GetStringValue(value)
		return res, err
	// REG_DWORD + REG_DWORD_BIG_ENDIAN + REG_QWORD
	case 4, 5, 11:
		res, _, err := key.GetIntegerValue(value)
		return strconv.Itoa(int(res)), err
	// REG_MULTI_SZ
	case 7:
		res, _, err := key.GetStringsValue(value)
		resString := ""
		for n := range res {
			resString += res[n] + ","
		}
		return resString, err
		// Eventull noch weitere Typen
	}
	return "", errors.New("Unbekannter Fehler.")
}
