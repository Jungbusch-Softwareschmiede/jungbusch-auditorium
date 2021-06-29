//go:generate goversioninfo -icon=icon.ico
// Dieses Package führt das Jungbusch-Auditorium aus. Es ist unterteilt in einige Hilfs-Methoden,
// die das Ziel haben, auf den ersten Blick irrelevante Funktionalität (Error-Behandlung, Log-Ausgaben, etc.)
// aus der Main-Funktion zu entfernen um diese möglichst übersichtlich und lesbar zu halten.
// So kann man sich beispielsweise beim Debuggen leicht an den Methoden entlanghangeln und der Ablauf des
// Programms ist auf den ersten Blick ersichtlich.
package main

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/parser"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/syntaxchecker"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/modulecontroller"
	output "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/outputgenerator"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/privilege"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	var err error
	start := time.Now()

	// Temp-Ordner erstellen
	initializeTemp()

	// Rechte des Prozesses bestimmen
	static.HasElevatedPrivileges = privilege.HasRootPrivileges()

	// Programm-Konfiguration laden
	cs, log := config_parser.LoadConfig()

	// Logger initialisieren
	initializeLogger(&cs, log)

	// Betriebssystem erkennen und setzen
	static.OperatingSystem = setOs(&cs)

	// Die restliche Programm-Konfiguration validieren und interpretieren
	interpretConfig(&cs)

	// Module basierend auf Betriebssystem initialisieren (kompatible Module laden)
	loadedModules := loadModules(&cs)

	// Einlesen der Audit-Konfigurations-Datei
	Info(SeperateTitle("Audit-Parser"))
	lines, err := acutil.ReadAuditConfiguration(cs.AuditConfig)
	HandleError(err)

	// Audit-Konfiguration Syntax überprüfen
	syntaxCheck(lines, &cs)

	// Audit-Konfiguration Parsen
	auditModules, auditRequiresPrivileges, numberOfModules, err := parser.Parse(lines, loadedModules)
	HandleError(err)

	// Die benötigten Privilegien mit denen des Prozesses abgleichen
	checkPrivileges(auditRequiresPrivileges, &cs)

	// Die Audit-Module aus der Konfiguration interpretieren
	report, err := interpreter.InterpretAudit(auditModules, numberOfModules, cs.AlwaysPrintProgress)
	HandleError(err)

	// Erstellen des Reports und Speichern der Artefakte
	err = output.GenerateOutput(report, cs.OutputPath, numberOfModules, cs.Zip || cs.ZipOnly, start, time.Since(start))
	HandleError(err)

	// Die Log-Datei schließen, da wir das Ergebnis nicht zippen können, solange es in dem Prozess geöffnet ist
	CloseLogger()

	if cs.ZipOnly {
		if err = os.RemoveAll(cs.OutputPath); err != nil {
			panic(err) // Der Logger ist schon geschlossen, also können wir hier maximal panicen
		}
	}

	Exit()
}

// Initialisiert einen Temp-Ordner (nur Windows) in dem über Windows-Umgebungsvariablen festgelegten Pfad.
// Hier werden möglicherweise während dem Audit-Prozess temporär zu verwendende Dateien abgelegt.
// Der Ordner wird wieder entfernt, auch wenn im Programm ein Error auftritt.
func initializeTemp() {
	if runtime.GOOS == "windows" {
		path, err := filepath.Abs(os.Getenv("temp"))
		HandleError(err)

		// Pfad zum Ordner im Temp-Pfad resolven
		path, err = util.GetAbsolutePath(path + static.PATH_SEPERATOR + "JungbuschAuditorium")
		static.TempPath = path

		// Am Pfad eine Directory erstellen
		err = os.Mkdir(path, static.CREATE_DIRECTORY_PERMISSIONS)

		// Wenn der Error vom Typ ErrExist ist (Ordner existiert schon), ignorieren wir ihn
		if !errors.Is(err, fs.ErrExist) {
			HandleError(err)
		}
	}
}

// Intitialisiert den Logger. Tritt beim Initialisieren ein Fehler auf, wird die Panic-Methode des Loggers aufgerufen,
// die versucht an unterschiedlichen Pfaden einen Log zu erstellen, sodass der Nutzer über den Fehler
// informiert werden kann.
func initializeLogger(cs *ConfigStruct, log []LogMsg) {

	// Output-Pfad validieren, hier soll der Log liegen
	err := config_interpreter.ValidateOutputPath(cs)
	if err != nil {
		LogPanic(cs, []LogMsg{
			{
				Message:     "Fehler beim Validieren des Output-Pfads: " + err.Error(),
				Level:       1,
				AlwaysPrint: false,
			},
			{
				Message:     "Es konnte kein Log erstellt werden.",
				Level:       1,
				AlwaysPrint: false,
			},
		})
	}

	// Output-Ordner mit den Namen aus static erstellen
	err = InitializeOutput(cs)
	if err != nil {
		LogPanic(cs, []LogMsg{
			{
				Message:     "Fehler beim Erstellen des Output-Pfads: " + err.Error(),
				Level:       1,
				AlwaysPrint: false,
			},
			{
				Message:     "Es konnte kein Log erstellt werden.",
				Level:       1,
				AlwaysPrint: false,
			},
		})
	}

	// Den Logger selbst initialisieren
	err = InitializeLogger(cs, log)
	if err != nil {
		LogPanic(cs, []LogMsg{
			{
				Message:     "Fehler beim Initialisieren des Loggers: " + err.Error(),
				Level:       1,
				AlwaysPrint: false,
			},
			{
				Message:     "Es konnte kein Log erstellt werden.",
				Level:       1,
				AlwaysPrint: false,
			},
		})
	}
}

// Die Programmkonfiguration interpretieren
func interpretConfig(cs *ConfigStruct) {
	// Programm-Konfiguration validieren, interpretieren
	continueExecution, err := config_interpreter.InterpretConfig(cs)
	HandleError(err)

	// ContinueExecution ist dann false, wenn Commandline-Parameter angegeben wurde, nach denen das
	// Auditorium beendet. Bspw help, version, etc
	if !continueExecution {
		Exit()
	}
}

// Setzt das Betriebssystem ins Config-Struct
func setOs(cs *ConfigStruct) string {
	// Betriebssystem aus Konfiguration
	var err error
	currOS := cs.ForceOS
	if currOS == "" {
		// Das Betriebssystem bestimmen, wenn keins gesetzt wurde
		currOS, err = modulecontroller.GetOS()
		HandleError(err)
	}
	return currOS
}

// Initialisiert den Modulecontroller.
// Validiert den Syntax der Module und lädt sie.
func loadModules(cs *ConfigStruct) []ModuleSyntax {
	loadedModules, err := modulecontroller.Initialize(cs.SkipModuleCompatibilityCheck)
	HandleError(err)
	if cs.ShowModule != "" {
		InfoPrintAlways(modulecontroller.GetModuleSyntax(cs.ShowModule))
		Exit()
	}
	return loadedModules
}

// Checkt den Syntax der Audit-Konfigurationsdatei
func syntaxCheck(lines []string, cs *ConfigStruct) {
	// Syntax-Check
	err := syntaxchecker.Syntax(lines)
	HandleError(err)

	if cs.CheckSyntax {
		InfoPrintAlways("Der Syntaxcheck der Audit-Konfigurationsdatei wurde fehlerfrei beendet. Das Programm beendet jetzt.")
		Exit()
	}
}

// Überprüft die von der Audit-Konfiguration benötigten Privilegien und vergleicht sie mit denen des Prozesses
func checkPrivileges(auditRequiresPrivileges bool, cs *ConfigStruct) {
	// Privilegien checken
	if auditRequiresPrivileges {
		if !static.HasElevatedPrivileges {
			// Module, die Admin benötigen, wir haben aber kein Admin -> Programm beenden wenn CLI-Parameter nicht gesetzt ist
			if !cs.IgnoreMissingPrivileges {
				Err(Seperate())
				ErrAndExit("In der Audit-Konfiguration sind Module definiert, die Administrator, bzw. Root-Privilegien benötigen, das Jungbusch Auditorium wurde aber ohne solche Privilegien gestartet. Das Programm beendet nun.")
			} else {
				Warn("In der Audit-Konfiguration sind Module definiert, die Administrator, bzw. Root-Privilegien benötigen, das Jungbusch Auditorium wurde aber ohne solche Privilegien gestartet. Es wurde der Parameter <IgnoreMissingPrivileges> angegeben, das JBA wird soweit ausgeführt wie möglich.")
			}
		}
	} else {
		if static.HasElevatedPrivileges {
			// Keine Module, die Admin benötigen, wir haben aber Admin -> Warnung
			Warn("In der Audit-Konfiguration sind keine Module definiert, die Administrator, bzw. Root-Privilegien benötigen, das Jungbusch-Auditorium wurde aber mit solchen gestartet.")
		}
	}

	// Wenn wir nur den Syntax der Konfiguration überprüfen sollen, hören wir hier auf
	if cs.CheckConfiguration {
		InfoPrintAlways("Die Konfigurations-Datei wurde fehlerfrei geparsed. Das Programm beendet jetzt.")
		Exit()
	}
}
