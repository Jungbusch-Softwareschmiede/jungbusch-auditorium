package config_parser

import (
	"flag"
	"fmt"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	s "strings"
)

// In diesem Struct werden die Beschreibungen für die Commandline-Flags gespeichert, die für das Anzeigen der
// Usage-Page benötigt werden
type cliFlag struct {
	name     string
	desc     string
	datatype string
}

// Beschreibung der einzelnen Parameter: Siehe models.ConfigStruct
var (
	// Programm-Parameter
	auditConfig_Flag                  string
	config_Flag                       string
	outputPath_Flag                   string
	verbosityLog_Flag                 int
	verbosityConsole_Flag             int
	skipModuleCompatibilityCheck_Flag bool
	keepConsoleOpen_Flag              bool
	forceOS_Flag                      string
	ignoreMissingPrivileges_Flag      bool
	zip_Flag                          bool
	zipOnly_Flag                      bool

	// One-And-Done-Parameter (Programm führt keine Audits aus)
	version_Flag             bool
	createDefaultConfig_Flag bool
	showModules_Flag         bool   // Alle Module (bool)
	showModuleInfo_Flag      string // Ein einzelnes Modul (string)
	checkConfiguration_Flag  bool
	checkSyntax_Flag         bool
	alwaysPrintProgress_Flag bool

	// Misc
	saveConfiguration_Flag bool // Speichern der Commandline-Parameter in Config
)

// Diese Methode initialisiert die Flags der Commandline, enthält die Beschreibungen der Flags, sowie die Help-Page
func initializeFlag() {
	if len(flags) == 0 {
		flag.StringVar(&auditConfig_Flag, "auditConfig", "", "")
		flag.StringVar(&auditConfig_Flag, "a", "", "")

		flag.StringVar(&config_Flag, "config", "", "")
		flag.StringVar(&config_Flag, "c", "", "")

		flag.StringVar(&outputPath_Flag, "outputPath", "", "")
		flag.StringVar(&outputPath_Flag, "o", "", "")

		flag.IntVar(&verbosityLog_Flag, "verbosityLog", -1, "")
		flag.IntVar(&verbosityLog_Flag, "vl", -1, "")

		flag.IntVar(&verbosityConsole_Flag, "verbosityConsole", -1, "")
		flag.IntVar(&verbosityConsole_Flag, "vc", -1, "")

		flag.BoolVar(&skipModuleCompatibilityCheck_Flag, "skipModuleCompatibilityCheck", false, "")

		flag.BoolVar(&keepConsoleOpen_Flag, "keepConsoleOpen", false, "")
		flag.StringVar(&forceOS_Flag, "forceOS", "", "")
		flag.BoolVar(&ignoreMissingPrivileges_Flag, "ignoreMissingPrivileges", false, "")
		flag.BoolVar(&zip_Flag, "zip", false, "")
		flag.BoolVar(&zipOnly_Flag, "zipOnly", false, "")

		flag.BoolVar(&version_Flag, "version", false, "")
		flag.BoolVar(&createDefaultConfig_Flag, "createDefault", false, "")
		flag.BoolVar(&showModules_Flag, "showModules", false, "")
		flag.StringVar(&showModuleInfo_Flag, "showModuleInfo", "", "")
		flag.BoolVar(&checkConfiguration_Flag, "checkConfiguration", false, "")
		flag.BoolVar(&checkSyntax_Flag, "checkSyntax", false, "")
		flag.BoolVar(&checkSyntax_Flag, "syntax", false, "")
		flag.BoolVar(&alwaysPrintProgress_Flag, "alwaysPrintProgress", false, "")

		flag.BoolVar(&saveConfiguration_Flag, "saveConfiguration", false, "")
		flag.BoolVar(&saveConfiguration_Flag, "s", false, "")

		// Initialisieren aller Flags
		flags = []cliFlag{
			//
			// Version
			//
			{
				name:     "-version",
				desc:     "Version & Datum der letzten Änderung",
				datatype: "-",
			},

			//
			// Help
			//
			{
				name:     "-help",
				desc:     "Hilfe-Seite",
				datatype: "-",
			},

			//
			// ShowModules
			//
			{
				name:     "-showModules",
				desc:     "Liste mit allen Modulnamen",
				datatype: "-",
			},

			//
			// Show specific module
			//
			{
				name:     "-showModuleInfo",
				desc:     "Gibt verfügbare Informationen zu dem Modul aus",
				datatype: "<Modulname>",
			},

			//
			// JustCheckConfiguration
			//
			{
				name:     "-checkConfiguration",
				desc:     "Überprüft Syntax und Werte der Audit-Konfigurationsdatei",
				datatype: "-",
			},

			//
			// syntaxCheck
			//
			{
				name:     "-checkSyntax, -syntax",
				desc:     "Überprüft nur die Syntax der Audit-Konfigurationsdatei",
				datatype: "-",
			},

			//
			// Create Default-Flag
			//
			{
				name:     "-createDefault",
				desc:     "Erzeugt config.ini-Datei mit Default-Werten am Pfad der Executable",
				datatype: "-",
			},

			//
			// Audit-Config
			//
			{
				name:     "-auditConfig, -a",
				desc:     "Pfad zur Audit-Konfigurations-Datei",
				datatype: "<Pfad>",
			},

			//
			// Config
			//
			{
				name:     "-config, -c",
				desc:     "Pfad zur Konfigurations-Datei",
				datatype: "<Pfad>",
			},

			//
			// Output-Path
			//
			{
				name:     "-outputPath, -o",
				desc:     "Pfad zur Output-Directory",
				datatype: "<Pfad>",
			},

			//
			// Log-Verbosity
			//
			{
				name:     "-verbosityLog, -vl",
				desc:     "Menge der Log-Einträge",
				datatype: "<0-4>",
			},

			//
			// Console-Verbosity
			//
			{
				name:     "-verbosityConsole, -vc",
				desc:     "Menge der Konsolen-Einträge",
				datatype: "<0-4>",
			},

			//
			// SkipModuleCompatibility
			//
			{
				name:     "-skipModuleCompatibilityCheck",
				desc:     "Überspringen der internen Modul-Kompatibilitätsüberprüfung",
				datatype: "-",
			},

			//
			// KeepConsoleOpen
			//
			{
				name:     "-keepConsoleOpen",
				desc:     "Verhindert das sofortige schließen der Konsole nach dem Ausführen des Programms per Doppelclick",
				datatype: "-",
			},

			//
			// OS
			//
			{
				name:     "-forceOS",
				desc:     "Überschreibt das Ergebnis des OS-Detectors",
				datatype: "<OS Name>",
			},

			//
			// IgnoreMissingPrivileges
			//
			{
				name:     "-ignoreMissingPrivileges",
				desc:     "Das JBA wird trotz fehlenden Privilegien soweit möglich ausgeführt und nicht frühzeitig beendet",
				datatype: "-",
			},

			//
			// zip
			//
			{
				name:     "-zip",
				desc:     "Erstellt eine Zip-Datei des Output-Ordners",
				datatype: "-",
			},

			//
			// zipOnly
			//
			{
				name:     "-zipOnly",
				desc:     "Erstellt eine Zip-Datei des Output-Ordners und löscht den Output-Ordner",
				datatype: "-",
			},

			//
			// AlwaysPrintProgress
			//
			{
				name:     "-alwaysPrintProgress",
				desc:     "Der Fortschritt der Audit-Schritte wird unabhängig vom Log-Level ausgegeben",
				datatype: "-",
			},

			//
			// Save cli-params to config
			//
			{
				name:     "-saveConfiguration, -s",
				desc:     "Mit diesem Befehl werden die aktuellen Commandline-Parameter in die config.ini-Datei geschrieben",
				datatype: "-",
			},
		}
		add(LogMsg{Message: "Die Commandline-Flags wurden initialisiert.", Level: 4})
	}

	// Was beim Ausführen von -help oder der inkorrekten Nutzung von Flags ausgegeben wird
	flag.Usage = func() {
		// Ausgabe der Usage-Page
		fmt.Println("Usage:")
		fmt.Println("./jungbusch-auditorium | jungbusch-auditorium.exe [NonExec] [Exec] <var>")
		fmt.Println("\nCommandline-Parameter bei denen keine Audit-Schritte ausgeführt werden:")
		for n := 0; n < 6; n++ {
			fmt.Println(fmt.Sprintf("%-30v", flags[n].name) + fmt.Sprintf("%-17v", flags[n].datatype) + flags[n].desc)
		}
		fmt.Println("\nCommandline-Parameter bei denen Audit-Schritte ausgeführt werden:")
		for n := 6; n < len(flags); n++ {
			fmt.Println(fmt.Sprintf("%-35v", flags[n].name) + fmt.Sprintf("%-17v", flags[n].datatype) + flags[n].desc)
		}
		fmt.Println("\nWeitere Informationen finden Sie im Benutzerhandbuch.")
	}

	flag.Parse()
	auditConfig_Flag = s.TrimSpace(s.Trim(auditConfig_Flag, "\""))
	config_Flag = s.TrimSpace(s.Trim(config_Flag, "\""))
	outputPath_Flag = s.TrimSpace(s.Trim(outputPath_Flag, "\""))
	forceOS_Flag = s.TrimSpace(s.Trim(forceOS_Flag, "\""))
	showModuleInfo_Flag = s.TrimSpace(s.Trim(showModuleInfo_Flag, "\""))
}
