// In diesem Package wird die vorher eingelesene Konfiguration interpretiert und validiert. Das bedeutet, dass
// Parameter ohne Programmstart abgearbeitet werden und die angegebenen Pfade validiert werden.
package config_interpreter

import (
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/permissions"
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
	"reflect"
	s "strings"
)

// Diese Methode bekommt den Pointer zu einem ConfigStruct übergeben. Sie arbeitet ausstehende Parameter ab
// und validiert den Pfad zu der Audit-Konfiguration, sowie den Output-Pfad.
// Wenn keine Errors auftreten, gibt sie den Boolean continueExecution zurück. Ist dieser false,
// wird das Programm nach dem Ausführen dieser Methode beendet. Das ist beispielsweise dann der Fall,
// wenn nur die Version ausgegeben werden soll.
func InterpretConfig(cs *ConfigStruct) (continueExecution bool, err error) {
	Info(SeperateTitle("CLI-Interpreter"))
	switch {
	case cs.Version:
		InfoPrintAlways("Aktuelle Version: " + static.CURRENT_VERSION)
		return false, nil

	case cs.CreateDefaultConfig:
		if err = writeStructToConfigFile(*cs, "./config.ini"); err != nil {
			return false, errors.New("Beim Erstellen der Konfigurations-Datei ist Error aufgetreten: " + err.Error())
		} else {
			InfoPrintAlways("Es wurde erfolgreich eine Default-File erstellt.")
			return false, nil
		}

	case cs.SaveConfiguration:
		if err := writeStructToConfigFile(*cs, cs.Config); err != nil {
			return false, errors.New("Beim Speichern der Konfigurations-Datei ist ein Error augetreten: " + err.Error())
		} else {
			Info("Die Konfiguration wurde erfolgreich überschrieben!")
		}

	case cs.VerbosityLog == 0:
		Warn("Das Log-Level wurde auf NONE gesetzt. Es werden keine Nachrichten in den Log geschrieben.")

	case cs.VerbosityConsole == 0:
		Warn("Das Konsolen-Verbosity-Level wurde auf NONE gesetzt. In der Konsole werden keinerlei Ausgaben gemacht.")
	}

	// Pfad der Audit-Datei validieren
	err = validateAuditConfigPath(cs)
	if err != nil && cs.ShowModule == "" {
		return false, err
	} else {
		err = nil
	}

	if cs.Zip && cs.ZipOnly {
		return false, errors.New("-zip und -zipOnly widersprechen sich und dürfen nicht gleichzeitig gesetzt werden.")
	}

	return true, err
}

// Diese Methode validiert den Output-Pfad. Dafür wird überprüft ob der Pfad existiert und ob
// der Prozess Lese- und Schreibrechte hat.
func ValidateOutputPath(cs *ConfigStruct) (err error) {
	cs.OutputPath, err = util.GetAbsolutePath(cs.OutputPath)
	if err != nil {
		return err
	}

	read, write, isDir, err := permissions.Permission(cs.OutputPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("Der Output-Pfad existiert nicht oder ist kein Ordner: " + cs.OutputPath)
		}
		return err
	}

	if !(read && write) {
		return errors.New("Der Output-Pfad ist vom ausführenden Benutzer nicht les- und/oder schreibbar: " + cs.OutputPath)
	}

	if !isDir {
		return errors.New("Der Output-Pfad ist kein Ordner: " + cs.OutputPath)
	}

	return nil
}

// Diese Methode validiert den Pfad zur Audit-Konfiguration. Dafür wird überprüft ob der Pfad existiert und ob
// der Prozess Lese- und Schreibrechte hat. Des Weiteren wird überprüft, ob die angegebene Datei vom korrekten Typ ist.
func validateAuditConfigPath(cs *ConfigStruct) (err error) {
	// Relativen Pfad in absoluten Pfad konvertieren
	cs.AuditConfig, err = util.GetAbsolutePath(cs.AuditConfig)

	if err != nil {
		return err
	}
	read, _, isDir, err := permissions.Permission(cs.AuditConfig)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			path, _ := util.GetAbsolutePath(static.EXPECTED_AUDIT_PATH)
			if path == cs.AuditConfig {
				return errors.New("Es konnte keine Audit-Konfiguration gefunden werden.")
			} else {
				return errors.New("Es konnte keine Audit-Konfiguration an folgendem Pfad gefunden werden: " + cs.AuditConfig)
			}

		} else if errors.Is(err, fs.ErrPermission) {
			return errors.New("Der Prozess hat nicht die benötigten Berechtigungen um den Pfad der angegebenen Audit-Konfigurationsdatei zu öffnen. Wurde keine Datei explizit angegeben, hat der Prozess keine Berechtigungen für den Pfad der Executable, an welchem eine Audit-Konfigurationsdatei gesucht wird.")
		} else {
			return errors.New("Beim Bestimmen der Berechtigungen der Audit-Konfiguration ist ein unbekannter Fehler aufgetreten: " + err.Error())
		}
	}

	if isDir {
		return errors.New("Die angegebene Audit-Konfiguration ist keine Datei, sondern ein Pfad.")
	}

	if !read {
		return errors.New("Der ausführende Benutzer hat keine Lese-Rechte für die angegebene Audit-Konfiguration.")
	}

	// Dateiendung der Audit-Datei
	if cs.AuditConfig[s.LastIndex(cs.AuditConfig, ".")+1:] != "jba" {
		return errors.New("Die Audit-Config Datei ist nicht vom Format \".jba.\"")
	}

	Debug("Der Pfad zur Audit-Konfiguration ist valide.")
	return nil
}

// Diese Methode schreibt ein ConfigStruct in die Datei am übergebenen Pfad. Es werden nur Parameter in die Datei
// geschrieben, die keinen leeren Wert haben.
func writeStructToConfigFile(cs ConfigStruct, path string) error {
	r, w, isDir, err := permissions.Permission(".")
	if err != nil {
		return err
	}

	if r && w && isDir {
		var defaultConfig []string
		defaultConfig = append(defaultConfig, "[ENVIRONMENT]")

		val := reflect.ValueOf(&cs).Elem()
		for i := 0; i < val.NumField(); i++ {
			if val.Field(i).Interface() != "" && s.ToLower(val.Type().Field(i).Name) != "createdefaultconfig" && s.ToLower(val.Type().Field(i).Name) != "config" && s.ToLower(val.Type().Field(i).Name) != "saveconfiguration" {
				if s.ToLower(val.Type().Field(i).Name) == "outputpath" {
					defaultConfig = append(defaultConfig, fmt.Sprint(val.Type().Field(i).Name)+"="+filepath.Dir(fmt.Sprint(val.Field(i).Interface())))
				} else {
					defaultConfig = append(defaultConfig, fmt.Sprintf("%v=%v", val.Type().Field(i).Name, val.Field(i).Interface()))
				}
			}
		}
		path, err = util.GetAbsolutePath(path)
		if err != nil {
			return err
		}

		if err = util.CreateFile(defaultConfig, path); err != nil {
			return err
		}
	}
	return nil
}
