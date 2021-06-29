// Dieses package übernimmt das Schreiben der Log-Datei, das Ausgeben von Informationen auf der Konsole,
// sowie das Behandeln von Error-Nachrichten und das Beenden des Programms.
package logger

import (
	"bufio"
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/permissions"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	s "strings"
	"time"
)

// Das Logger-Objekt, welches im Hintergrund die Log-Datei schreibt
type logger struct {
	cmdVerbosity int      // Die Commandline-Stufe
	logVerbosity int      // Die Log-Stufe
	*log.Logger           // Der Logger selbst
	wait         bool     // True, wenn nach beenden des Programms gewartet werden soll
	file         *os.File // Die Log-Datei
}

// In diesem struct können die verfügbaren Log-Level gespeichert werden
type level struct {
	name string
	lv   int
}

var (
	l     logger
	info  level
	warn  level
	erro  level
	debug level
	none  level
)

// Initialisiert alle für die Verwendung des Loggers intern benötigten Werte
func InitializeLogger(cs *ConfigStruct, preLog []LogMsg) error {
	// Verbosity validieren damit wir so schnell wie möglich den Logger initialisieren können
	if cs.VerbosityConsole < 0 || cs.VerbosityConsole > 4 || cs.VerbosityLog < 0 || cs.VerbosityLog > 4 {
		return errors.New("Bitte ein Loglevel zwischen 0 und 4 angeben.")
	}

	// Log-Pfad generieren
	path, err := util.GetAbsolutePath(cs.OutputPath + static.PATH_SEPERATOR + static.LOG_NAME + ".log")

	// Datei öffnen
	l.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, static.CREATE_FILE_PERMISSIONS)
	if err != nil {
		return err
	}

	// Log-Objekt initialisieren
	l = logger{
		cmdVerbosity: cs.VerbosityConsole,
		logVerbosity: cs.VerbosityLog,
		wait:         cs.KeepConsoleOpen,
		file:         l.file,
		Logger:       log.New(l.file, "", log.Ltime),
	}

	// Verfügbare Log-Level initialisieren
	none = level{name: "", lv: 0}
	erro = level{name: "ERROR    ", lv: 1}
	warn = level{name: "WARNING  ", lv: 2}
	info = level{name: "INFO     ", lv: 3}
	debug = level{name: "DEBUG    ", lv: 4}

	// Die Log-Nachrichten, welche vor Initialisierung des Loggers gesammelt wurden ausgeben
	for n := range preLog {
		processPrelogMsg(preLog[n])
	}

	return nil
}

// Schließt den Logger
func CloseLogger() {
	err := l.file.Close()
	if err != nil {
		fmt.Println("Error beim schließen des Loggers: " + err.Error())
	}
}

// Diese Funktion wird dann aufgerufen, wenn vor Initialisierung des Loggers ein Fehler aufgetreten ist.
// Es wird versucht eine Log-Datei an unterschiedlichen Pfaden zu erstellen, bis dies erfolgreich ist.
func LogPanic(cs *ConfigStruct, preLog []LogMsg) {
	path, err := panicExit(cs, preLog)
	if err != nil {
		fmt.Println(panicExit(cs, preLog))
		fmt.Println("Das Jungbusch-Auditorium konnte keinen Log erstellen.\nDas Programm beendet nun.")
		Exit()
	} else {
		fmt.Println("Vor der initialisierung des Loggers ist ein Fehler aufgetreten. Eine Log-Datei wurde an folgendem Pfad erstellt: " + path)
	}
}

// Initialisiert den Output-Ordner. In diesem wird die Log-Datei abgelegt.
func InitializeOutput(cs *ConfigStruct) (err error) {
	// Pfad des Output-Ordners zusammenbauen
	if cs.OutputPath, err = util.GetAbsolutePath(cs.OutputPath + static.PATH_SEPERATOR + static.OUTPUT_FOLDER_NAME + "_" + time.Now().Format(static.OUTPUT_TIMESTAMP_FORMAT)); err != nil {
		return errors.New("Unbekannter Fehler beim initialisieren des Output-Ordners: " + err.Error())
	}

	// Output-Ordner erstellen
	if err = os.Mkdir(cs.OutputPath, static.CREATE_DIRECTORY_PERMISSIONS); err != nil {
		// Wenn der Ordner bereits existiert, wurden mehrere JBA-Instanzen zur selben Sekunde gestartet.
		// Wir lassen uns nicht anmerken dass irgendwas schief gelaufen ist und tun einfach so,
		// als wäre das ein Fehler den wir von Beginn an abfangen wollten ;)
		if errors.Is(err, fs.ErrExist) {
			return errors.New("Bitte nur jeweils eine Instanz der Jungbusch-Softwareschmiede starten.")
		} else {
			return errors.New("Fehler beim erstellen des Output-Pfads: " + err.Error())
		}
	}

	return
}

// Gibt eine Nachricht auf dem Info-Level aus. Wird always auf true gesetzt, wird die Nachricht auf der Konsole
// unabhängig vom Log-Level ausgegeben.
func InfoPrint(msg string, always bool) {
	logMe(msg, info, always)
}

// Gibt eine Nachricht auf dem Info-Level aus. Die Nachricht auf der Konsole wird unabhängig vom Log-Level ausgegeben.
func InfoPrintAlways(msg string) {
	logMe(msg, info, true)
}

// Gibt eine Nachricht auf dem Info-Level aus.
func Info(msg string) {
	logMe(msg, info, false)
}

// Gibt eine Nachricht auf dem Warn-Level aus.
func Warn(msg string) {
	logMe(msg, warn, false)
}

// Gibt eine Nachricht auf dem Error-Level aus.
func Err(msg string) {
	logMe(msg, erro, false)
}

// Gibt eine Nachricht auf dem Error-Level aus und beendet das Programm.
func ErrAndExit(msg string) {
	logMe(msg, erro, false)
	Exit()
}

// Gibt den übergebenen Error aus und beendet das Programm. Unterscheidet zwischen einem herkömmlichen Error
// und einem SyntaxError-Objekt.
func HandleError(err error) {
	if err != nil {
		exitString := ""
		Err(Seperate())

		switch err.(type) {
		case *SyntaxError:
			syntaxErr := err.(*SyntaxError)
			Err(syntaxErr.ErrorMsg)
			Err(createSyntaxErrorString(syntaxErr.Errorkeyword, syntaxErr.LineNo, syntaxErr.Line))

		default:
			Err(err.Error())
		}

		ErrAndExit(exitString)
	}
}

// Gibt eine Nachricht auf dem Debug-Level aus.
func Debug(msg string) {
	logMe(msg, debug, false)
}

// Löscht den Temp-Ordner und beendet das Programm.
func Exit() {
	// Temp-Ordner löschen
	if runtime.GOOS == "windows" && static.TempPath != "" {
		_ = os.RemoveAll(static.TempPath)
	}

	if l != (logger{}) {
		if l.wait {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("--- Enter drücken um fortzufahren ---")
			_, _, _ = reader.ReadLine()
		}
		os.Exit(1)
	} else {
		os.Exit(1)
	}
}

// Der übergebene String wird zwischen Linien gesetzt
func SeperateTitle(in string) string {
	return "—————————————————————————————————" + in + "—————————————————————————————————"
}

// Returned eine Linie
func Seperate() string {
	return "——————————————————————————————————————————————————————————————————"
}

// Übernimmt die Logik des Ausgeben der Nachrichten, sowie die Progressbar
func logMe(msg string, verbosity level, alwaysPrint bool) {
	if static.ProgressBar != nil {
		_ = static.ProgressBar.Clear()
	}

	if msg != "" {
		l.SetPrefix(verbosity.name)

		if l.cmdVerbosity >= verbosity.lv {
			fmt.Println(verbosity.name + time.Now().Format("15:04:05") + " - " + msg)
		} else if alwaysPrint {
			fmt.Print("                    ")
			fmt.Println(msg)
		}

		if l.logVerbosity >= verbosity.lv {
			l.Logger.Println("- " + msg)
		}
	}

	if static.ProgressBar != nil {
		_ = static.ProgressBar.RenderBlank()
	}
}

// Hält die Logik für die LogPanic-Methode
func panicExit(cs *ConfigStruct, preLog []LogMsg) (string, error) {
	// Ziel: Finden eines Pfads, an den wir einen Log schreiben können

	// 1. Versuch: Pfad der Executable
	r, w, isDir, err := permissions.Permission(".")
	if err == nil && r && w && isDir {
		// Hier kann kein err auftreten, das checkt permissions.Permission bereits
		path, _ := util.GetAbsolutePath(".")
		cs.OutputPath = path
		err = InitializeLogger(cs, preLog)

		// Wenn der err nicht nil ist, machen wir weiter
		if err == nil {
			return path, nil
		}
	}

	// 2. Versuch: Working-Directory
	pfad, err := filepath.Abs(".")
	if err == nil {
		r, w, isDir, err = permissions.Permission(pfad)
		if err == nil && r && w && isDir {
			cs.OutputPath = pfad
			err = InitializeLogger(cs, preLog)

			// Wenn der err nicht nil ist, machen wir weiter
			if err == nil {
				return pfad, nil
			}
		}
	}

	// 3. Versuch: System-spezifische Log-Pfade
	if runtime.GOOS == "windows" {
		cs.OutputPath = static.TempPath
		err = InitializeLogger(cs, preLog)

		// Wenn der err nicht nil ist, machen wir weiter
		if err == nil {
			return pfad, nil
		}
	} else {
		pfad = "/var/log/JungbuschAuditorium"
		err = os.Mkdir(pfad, static.CREATE_DIRECTORY_PERMISSIONS)

		cs.OutputPath = pfad
		err = InitializeLogger(cs, preLog)

		// Wenn der err nicht nil ist, machen wir weiter
		if err == nil {
			return pfad, nil
		}
	}

	return "", errors.New("Alle Log-Optionen fehlgeschlagen: " + err.Error())
}

// Diese Methode erstellt einen Error-String, der einen benutzerdefinierte Nachricht, sowie, wenn möglich, die genaue Position des Errors ausgiebt.
func createSyntaxErrorString(keyword string, lineNumber int, line string) string {
	var errorString string

	if lineNumber != -1 {
		index := s.Index(line, keyword)
		errorString = "Zeile " + strconv.Itoa(lineNumber) + ": " + line
		upIndex := (len(errorString) - len(line)) + index + (len(keyword) / 2) + 20

		errorString += "\n" + fmt.Sprintf("%"+strconv.Itoa(upIndex)+"v", "") + "^"
	} else {
		if len(line) > 0 {
			errorString = line
		} else {
			errorString = ""
		}
	}

	return errorString
}

// Verarbeitet die übergebenen LogMsg-Objekte des PreLoggers. Wenn die Nachricht vom Typ Error ist, wird sie ausgegeben
// und das Programm beendet, ansonsten wird sie wie gewohnt auf dem korrekten Level ausgegeben.
func processPrelogMsg(msg LogMsg) {
	if msg.Level == 1 {
		ErrAndExit(msg.Message)
	} else {
		logMe(msg.Message, getLevel(msg.Level), msg.AlwaysPrint)
	}
}

// Mappt einem Level vom Typ int ein level-Objekt zu.
func getLevel(lvl int) level {
	switch lvl {
	case 1:
		return erro
	case 2:
		return warn
	case 3:
		return info
	case 4:
		return debug
	default:
		return none
	}
}
