// Dieses Package verwaltet die Commandline-Flags, sucht und liest eine Konfigurations-Datei ein und liefert ein Objekt
// vom Typ models.ConfigStruct mit allen gesetzten Werten zurück. Eine Besonderheit des Packages ist, dass es anders
// als alle anderen JBA-Packages nicht den Logger verwendet sondern stattdessen ein Slice aus models.LogMsg-Objekten
// zurückgibt. Dies liegt daran, dass der von anderen Packages verwendete Logger erst nach dem Parsen der Konfiguration
// initialisiert werden kann, da hier beispielsweise das Log-Level gesetzt wird. Daher werden Log-Einträge "gesammelt"
// und bei der Initialisierung des Loggers sozusagen nachträglich in den Log geschrieben.
package config_parser

import (
	"flag"
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strconv"
	s "strings"
)

var (
	tempLogger []LogMsg
	flags      []cliFlag
)

// Diese Methode ist die erste Methode, die im Jungbusch-Auditorium aufgerufen wird. Sie erzeugt ein
// models.ConfigStruct-Objekt mit den in static definierten Default-Werten. Ein Pointer zu diesem Objekt wird dann an
// die configHandler-Methode übergeben, welche die gesetzten Werte mit denen aus der config.ini-Datei überschreibt.
// Daraufhin wird die commandlineHandler-Methode mit demselben Pointer aufgerufen, welche alle auf der Commandline
// angegebenen Parameter parsed und in das Struct setzt.
func LoadConfig() (ConfigStruct, []LogMsg) {
	tempLogger = make([]LogMsg, 0)
	add(LogMsg{Message: "Das Jungbusch-Auditorium wurde gestartet.", Level: 3})
	add(LogMsg{Message: logger.SeperateTitle("Commandline-Parser"), Level: 3})
	initializeFlag()

	cs := ConfigStruct{
		AuditConfig:      static.EXPECTED_AUDIT_PATH,
		OutputPath:       static.DEFAULT_OUTPUT_PATH,
		VerbosityLog:     static.DEFAULT_LOG_VERBOSITY,
		VerbosityConsole: static.DEFAULT_CONSOLE_VERBOSITY,
		Config:           config_Flag,
	}

	// Zuerst wird die Config-Datei gesucht, eingelesen und in das Struct gesetzt (wenn vorhanden)
	err := configHandler(&cs)
	if err != nil {
		add(LogMsg{Message: err.Error(), Level: 1})
	}

	// Anschließend werden die Commandline-Parameter geparsed.
	// Sie überschreiben dabei die Werte, welche die Config-Datei gesetzt hat
	err = commandlineHandler(&cs)
	if err != nil {
		add(LogMsg{Message: err.Error(), Level: 1})
	}

	return cs, tempLogger
}

// Diese Funktion nimmt ein ConfigStruct entgegen, welches bisher nur mit Default-Werten gefüllt ist.
// Es wird eine Konfigurations-Datei gesucht und eingelesen.
// Für eine Beschreibung, welche Konfigurations-Datei eingelesen wird, siehe Handbuch
func configHandler(cs *ConfigStruct) (err error) {
	// Nach einer JBA suchen, wenn genau eine vorhanden ist die auch setzten
	if err = findAuditConfigs(cs); err != nil {
		// Wir müssen hier keinen Error schmeißen, da sonst das Programm unnötig beendet wird, wir geben stattdessen nur eine Warnung aus
		add(LogMsg{Message: "Unbekannter Fehler beim Suchen von .jba-Dateien im Pfad der Executable: " + err.Error(), Level: 2})
	}

	if cs.Config != "" {
		cs.Config, err = util.GetAbsolutePath(cs.Config)
		if err != nil {
			return
		}

		if !s.HasSuffix(cs.Config, ".ini") {
			return errors.New("Die Konfigurations-Datei muss die Dateiendung .ini haben.")
		}

		// Wenn explizit eine Config-File angegeben wurde, diese aber nicht existiert beenden wir das Programm mit einem error
		add(LogMsg{Message: "Die Konfigurationsdatei am Pfad \"" + cs.Config + "\" wird verwendet.", Level: 4})
		err = readConfigFile(cs)
		if err != nil {
			return err
		}
	} else {
		if cs.Config == "" {
			// Es wurde keine Config angegeben, wir setzen den Pfad auf den Default-Wert
			cs.Config = static.EXPECTED_CONFIG_PATH
		}

		cs.Config, err = util.GetAbsolutePath(cs.Config)
		if err != nil {
			return
		}

		err = readConfigFile(cs)

		// Falls die Datei nicht gefunden wurde, ignorieren wir den Error und machen erstmal mit den Default-Werten weiter
		if err != nil {
			if s.Contains(err.Error(), "existiert nicht") {
				add(LogMsg{Message: "Es wurde keine Konfigurations-Datei gefunden oder angegeben. Die Default-Werte werden verwendet.", Level: 2})
			} else {
				add(LogMsg{Message: "Fehler in der config.ini: " + err.Error(), Level: 1})
				return
			}
			err = nil
		} else {
			add(LogMsg{Message: "Die Konfigurationsdatei am Pfad \"" + static.EXPECTED_CONFIG_PATH + "\" wird eingelesen.", Level: 4})
		}
	}
	return err
}

// Diese Funktion nimmt ein ConfigStruct entgegen, welches mit Werten aus der Konfigurations-Datei
// und/oder Default-Werten gefüllt ist.
// Diese Werte werden nun von den per Commandline angegebenen Parametern überschrieben, wenn vorhanden.
func commandlineHandler(cs *ConfigStruct) (err error) {
	if flag.NFlag() > 0 {
		add(LogMsg{Message: "Es wurden " + strconv.Itoa(flag.NFlag()) + " Commandline-Parameter angegeben.", Level: 4})
		err = setCliParameters(cs)
		if err != nil {
			return err
		}
	} else {
		add(LogMsg{Message: "Es wurden keine Commandline-Parameter angegeben.", Level: 4})
	}
	return err
}

// In dieser Funktion werden die Werte aller Commandline-Parameter ausgelesen, ggf. in den korrekten Datentyp
// konvertiert und in das Config-Struct gesetzt
func setCliParameters(cs *ConfigStruct) (err error) {
	if auditConfig_Flag != "" {
		cs.AuditConfig = auditConfig_Flag
	}

	if outputPath_Flag != "" {
		cs.OutputPath = outputPath_Flag
	}

	if verbosityLog_Flag != -1 || argsContainFlag("verbosityLog") || argsContainFlag("vl") {
		cs.VerbosityLog = verbosityLog_Flag
	}

	if verbosityConsole_Flag != -1 || argsContainFlag("verbosityConsole") || argsContainFlag("vc") {
		cs.VerbosityConsole = verbosityConsole_Flag
	}

	fmt.Println(os.Args)
	if skipModuleCompatibilityCheck_Flag || argsContainFlag("skipModuleCompatibilityCheck") {
		fmt.Println(skipModuleCompatibilityCheck_Flag)
		cs.SkipModuleCompatibilityCheck = skipModuleCompatibilityCheck_Flag
	}

	if keepConsoleOpen_Flag || argsContainFlag("keepConsoleOpen") {
		cs.KeepConsoleOpen = keepConsoleOpen_Flag
	}

	if forceOS_Flag != "" {
		cs.ForceOS = forceOS_Flag
	}

	if ignoreMissingPrivileges_Flag || argsContainFlag("ignoreMissingPrivileges") {
		cs.IgnoreMissingPrivileges = ignoreMissingPrivileges_Flag
	}

	if zip_Flag || argsContainFlag("zip") {
		cs.Zip = zip_Flag
	}

	if zipOnly_Flag || argsContainFlag("zipOnly") {
		cs.ZipOnly = zipOnly_Flag
	}

	if version_Flag || argsContainFlag("version") {
		cs.Version = version_Flag
	}

	if createDefaultConfig_Flag || argsContainFlag("createDefault") {
		cs.CreateDefaultConfig = createDefaultConfig_Flag
	}

	if showModules_Flag || argsContainFlag("showModules") {
		cs.ShowModule = "all"
	}

	if showModuleInfo_Flag != "" || argsContainFlag("showModuleInfo") {
		cs.ShowModule = showModuleInfo_Flag
	}

	if checkConfiguration_Flag || argsContainFlag("checkConfiguration") {
		cs.CheckConfiguration = checkConfiguration_Flag
	}

	if checkSyntax_Flag || argsContainFlag("checkSyntax") || argsContainFlag("syntax") {
		cs.CheckSyntax = checkSyntax_Flag
	}

	if alwaysPrintProgress_Flag || argsContainFlag("alwaysPrintProgress") {
		cs.AlwaysPrintProgress = alwaysPrintProgress_Flag
	}

	if saveConfiguration_Flag || argsContainFlag("saveConfiguration") || argsContainFlag("s") {
		cs.SaveConfiguration = saveConfiguration_Flag
	}

	return err
}

func argsContainFlag(in string) bool {
	for _, arg := range os.Args {
		if s.Contains(arg, "--") || s.Contains(arg, "-") {
			if s.HasPrefix(s.ToLower(s.Trim(arg, "-")), s.ToLower(in)) {
				return true
			}
		}
	}
	return false
}

// Diese Funktion findet am Pfad der Executable alle Dateien, welche die Dateiendung .jba haben.
// Wenn genau eine jba-Datei gefunden und keine andere explizit angegeben ist,
// wird diese automatisch vom Programm verwendet.
func findAuditConfigs(cs *ConfigStruct) error {
	var files []string

	// Alle Dateinamen aus dem Verzeichnis holen
	path, err := util.GetAbsolutePath("./")
	if err != nil {
		return err
	}
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// Nach .jba Dateien suchen
	for _, file := range fileInfo {
		if file.Name()[s.LastIndex(file.Name(), ".")+1:] == "jba" {
			files = append(files, file.Name())
		}
	}

	// Wenn genau eine Datei im Pfad vorhanden ist, verwenden wir diese
	if len(files) == 1 {
		cs.AuditConfig, err = util.GetAbsolutePath(files[0])
		if err != nil {
			return err
		}
	}
	return nil
}

// Liest die übergebene Konfigurations-Datei vom System ein. Dabei werden Fehler erzeugt, wenn Key oder Wert fehlen,
// oder wenn ein ungültiger Key angegeben wurde. Alle Werte werden in das übergebene ConfigStruct geschrieben.
func readConfigFile(cs *ConfigStruct) (err error) {
	// Config Datei laden
	configSlice, err := util.ReadFile(cs.Config)
	if err != nil {
		return errors.New("Unbekannter Fehler beim Einlesen der Konfigurations-Datei: " + err.Error())
	}

	// Durch die eingelesene Datei iterieren und die Parameter setzen
	for _, line := range configSlice {
		if s.TrimSpace(line) != "" {
			if s.Contains(line, "=") {
				line = s.Replace(line, "\"", "", -1)
				values := s.SplitN(line, "=", 2)

				values[0] = s.Replace(values[0], " ", "", -1)
				values[1] = s.TrimSpace(values[1])

				if values[1] == "" {
					return errors.New("Der Wert eines Keys darf nicht leer sein.")
				}

				switch s.ToLower(values[0]) {

				case "auditconfig":
					cs.AuditConfig = values[1]

				case "config":
					// nichts

				case "outputpath":
					cs.OutputPath = values[1]

				case "verbositylog":
					verbosityLog, err := getLogLevelFromString(values[1])
					if err != nil {
						return err
					}
					cs.VerbosityLog = verbosityLog

				case "verbosityconsole":
					verbosityLog, err := getLogLevelFromString(values[1])
					if err != nil {
						return err
					}
					cs.VerbosityConsole = verbosityLog

				case "skipmodulecompatibilitycheck":
					if err = parseAndSetBool(&cs.SkipModuleCompatibilityCheck, values[1]); err != nil {
						return err
					}

				case "keepconsoleopen":
					if err = parseAndSetBool(&cs.KeepConsoleOpen, values[1]); err != nil {
						return err
					}

				case "forceOS":
					cs.ForceOS = values[1]

				case "ignoremissingprivileges":
					if err = parseAndSetBool(&cs.IgnoreMissingPrivileges, values[1]); err != nil {
						return err
					}

				case "zip":
					if err = parseAndSetBool(&cs.Zip, values[1]); err != nil {
						return err
					}

				case "ziponly":
					if err = parseAndSetBool(&cs.ZipOnly, values[1]); err != nil {
						return err
					}

				case "version":
					if err = parseAndSetBool(&cs.Version, values[1]); err != nil {
						return err
					}

				case "showmodules":
					if rslt, err := util.ParseStringToBool(values[1]); err == nil && rslt {
						cs.ShowModule = "all"
					} else {
						return err
					}

				case "showmoduleinfo":
					if values[1] != "" {
						cs.ShowModule = showModuleInfo_Flag
					}

				case "checkconfiguration":
					if err = parseAndSetBool(&cs.CheckConfiguration, values[1]); err != nil {
						return err
					}

				case "checksyntax":
					if err = parseAndSetBool(&cs.CheckSyntax, values[1]); err != nil {
						return err
					}

				case "alwaysprintprogress":
					if err = parseAndSetBool(&cs.AlwaysPrintProgress, values[1]); err != nil {
						return err
					}

				case "saveconfiguration":
					if err = parseAndSetBool(&cs.SaveConfiguration, values[1]); err != nil {
						return err
					}

				case "createdefaultconfig":
					if err = parseAndSetBool(&cs.CreateDefaultConfig, values[1]); err != nil {
						return err
					}

				default:
					return errors.New("Ungültiger Key: " + values[0])
				}
			} else {
				if !(s.HasPrefix(s.TrimSpace(line), "[") && s.HasSuffix(s.TrimSpace(line), "]")) {
					return errors.New("Ungültiger Wert (fehlender Key): " + line)
				}
			}
		}
	}

	return nil
}

// Parsed den Wert eines übergebenen Strings in einen Boolean und setzt das Ergebnis
// in den übergebenen boolean-Pointer. Verwendet parseStringToBool.
func parseAndSetBool(toSet *bool, value string) (err error) {
	var rslt bool
	if rslt, err = util.ParseStringToBool(value); err == nil {
		*toSet = rslt
	}
	return err
}

// Konvertiert das übergebene Loglevel (String) in einen int
func getLogLevelFromString(in string) (int, error) {
	out, err := strconv.Atoi(in)
	if err != nil {
		return 0, errors.New("Das angegebene Loglevel ist keine ganze Zahl.")
	} else {
		return out, nil
	}
}

// Wrapper für Slice-Append
func add(msg LogMsg) {
	tempLogger = append(tempLogger, msg)
}

// Wird aufgrund eines Golang-Quirks in den Tests benötigt
func ResetFlags() {
	auditConfig_Flag = ""
	config_Flag = ""
	outputPath_Flag = ""

	verbosityLog_Flag = -1
	verbosityConsole_Flag = -1
	skipModuleCompatibilityCheck_Flag = false
	keepConsoleOpen_Flag = false
	forceOS_Flag = ""
	ignoreMissingPrivileges_Flag = false
	alwaysPrintProgress_Flag = false
	version_Flag = false
	createDefaultConfig_Flag = false
	showModules_Flag = false
	showModuleInfo_Flag = ""
	checkConfiguration_Flag = false
	checkSyntax_Flag = false
	zipOnly_Flag = false
	zip_Flag = false

	saveConfiguration_Flag = false
}
