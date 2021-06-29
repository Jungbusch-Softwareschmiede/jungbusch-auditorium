package util

import (
	"bufio"
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/permissions"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	s "strings"
)

// Gibt den absoluten Pfad des übergebenen Pfads zurück. Geht vom Pfad der Executable aus.
// Die aktuell working-directory wird ignoriert.
func GetAbsolutePath(currentPath string) (path string, err error) {
	path = filepath.ToSlash(currentPath)
	if !filepath.IsAbs(currentPath) {
		path, err = filepath.Abs(filepath.Dir(os.Args[0]) + static.PATH_SEPERATOR + filepath.Clean(currentPath))
		return path, err
	}
	return currentPath, nil
}

// Castet den übergebenen String entweder zu einem Boolean, Float oder String, wenn nichts anderes zutrifft
func CastToAppropriateType(value string) interface{} {
	value = s.Trim(value, "\"")
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	} else if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num
	} else {
		return value
	}
}

// Returned true, wenn der Array den übergebenen String enthält
func ArrayContainsString(arr []string, val string) bool {
	if arr != nil {
		for _, a := range arr {
			if s.ToLower(a) == s.ToLower(val) {
				return true
			}
		}
	}
	return false
}

// Entfernt alle Leerzeichen in einem String
func CompressString(in string) string {
	in = s.Replace(in, " ", "", -1)
	return in
}

// Entfernt alle im Array übergebenen Strings aus dem String
func RemoveFromString(in string, toReplace []string) string {
	for n := range toReplace {
		in = s.Replace(in, toReplace[n], "", -1)
	}
	return in
}

// Liest eine utf-8-Datei ein
func ReadFile(filename string) ([]string, error) {
	return readFileWorker(filename, "utf8")
}

// Liest eine utf16-Datei ein
func ReadUTF16File(filename string) ([]string, error) {
	return readFileWorker(filename, "utf16")
}

// Hilfs-Funktion zum Einlesen von Dateien
func readFileWorker(filename string, mode string) ([]string, error) {
	// Pfad in absoluten Pfad konvertieren
	filename, err := GetAbsolutePath(filename)
	if err != nil {
		return nil, err
	}

	// Lese-Berechtigungen am angegebenen Pfad bestimmen
	r, _, isDir, err := permissions.Permission(filename)

	switch {
	case err != nil:
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New("Die angegebene Datei existiert nicht: " + filename)
		}
		return nil, errors.New(err.Error())

	case isDir:
		return nil, errors.New("Es muss eine Datei angegeben werden, kein Ordner.")

	case !r:
		return nil, errors.New("Dem ausführenden Benutzer fehlen die Lese-Rechte für die Datei am folgenden Pfad: " + filename)
	}

	audFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer audFile.Close()

	var audLines []string

	var scanner *bufio.Scanner
	switch mode {
	case "utf16":
		scanner = bufio.NewScanner(transform.NewReader(audFile, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()))
	default:
		scanner = bufio.NewScanner(audFile)
	}

	for scanner.Scan() {
		audLines = append(audLines, s.Replace(scanner.Text(), "\t", "", -1))
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return audLines, nil
}

// Gibt true zurück, wenn der übergebene Pfad eine Datei ist.
func IsFile(name string) bool {
	name, err := GetAbsolutePath(name)
	if err != nil {
		return false
	}
	if info, err := os.Stat(name); err == nil {
		if info != nil {
			if !info.IsDir() {
				return true
			}
			return false
		}
		return false
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// Gibt true zurück, wenn der übergebene Pfad ein Verzeichniss ist.
func IsDir(name string) bool {
	name, err := GetAbsolutePath(name)
	if err != nil {
		return false
	}
	if info, err := os.Stat(name); err == nil {
		if info != nil {
			if info.IsDir() {
				return true
			}
			return false
		}
		return false
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

// Erstellt eine Datei. Ist dieselbe Datei bereits vorhanden, wird sie überschrieben.
func CreateFile(data []string, fileName string) error {
	fileName, err := GetAbsolutePath(fileName)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, static.CREATE_FILE_PERMISSIONS)

	if err != nil {
		return err
	}

	datawriter := bufio.NewWriter(file)

	for _, line := range data {
		_, _ = datawriter.WriteString(line + "\n")
	}

	datawriter.Flush()
	file.Close()

	return nil
}

// Parsed den Wert eines übergebenen Strings in einen Boolean und gibt diesen zurück.
func ParseStringToBool(in string) (bool, error) {
	in = s.TrimSpace(s.ToLower(in))

	if in == "1" || in == "t" || in == "true" {
		return true, nil
	}

	if in == "0" || in == "f" || in == "false" {
		return false, nil
	}

	return false, errors.New("Der Wert für <" + in + "> darf nur <true/false> sein.")
}

// Gibt den übergebenen string-Array auf der Konsole aus
func PrintStrArray(in []string) string {
	var out string

	for n := range in {
		out += in[n] + ", "
	}

	if len(out) > 2 {
		out = out[:len(out)-2]
	}

	return out
}

// Returned den string aus str der zwischen strings (rem) steht
func GetStringInBetween(str string, rem string) string {
	start := s.Index(str, rem)
	if start == -1 {
		return ""
	}
	end := s.LastIndex(str, rem)
	if end == -1 {
		return ""
	}
	return str[start+1 : end]
}

//
//
//
// Konsolenausgaben (Debug)
//
//
//
func AuditModulePrinter(mod models.AuditModule, tabs int) {
	currTabs := tabs
	tabs += 14

	fmt.Println(fmt.Sprintf("%"+strconv.Itoa(tabs)+"v", "Modul-Name: ") + fmt.Sprintf("%"+strconv.Itoa(len(mod.ModuleName)+3)+"v", mod.ModuleName))
	fmt.Println(fmt.Sprintf("%"+strconv.Itoa(tabs)+"v", "Passed: ") + fmt.Sprintf("%"+strconv.Itoa(len(mod.Passed)+3)+"v", mod.Passed))
	fmt.Print(fmt.Sprintf("%"+strconv.Itoa(tabs)+"v", "Variables: ") + fmt.Sprintf("%"+strconv.Itoa(3)+"v", ""))
	fmt.Println(mod.Variables)
	fmt.Print(fmt.Sprintf("%"+strconv.Itoa(tabs)+"v", "Modul-Param: ") + fmt.Sprintf("%"+strconv.Itoa(3)+"v", ""))
	fmt.Println(mod.ModuleParameters)

	for n := range mod.NestedModules {
		fmt.Println()
		fmt.Println("-----------Nested")
		AuditModulePrinter(mod.NestedModules[n], currTabs+5)
	}
}
