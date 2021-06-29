// In diesem Package werden im Jungbusch-Auditorium an unterschiedlichen Stellen verwendete Datentypen gesammelt.
// Dies hat den Vorteil, dass so Import-Loops vermieden werden. Hat das Package `Parser` beispielsweise das struct
// `AuditModule` und eine `Utility-Methode` bekommt ein solches Objekt zur Verwendung der Methode aus dem `Parser` übergeben,
// dann wurde ein Loop erschaffen, da Util den Parser wegen des structs importieren muss und der Parser wiederum Util importiert
// um die Methode verwenden zu können. Das Programm wird so nicht kompilieren.
package models

type ConfigStruct struct {
	// Programm-Parameter
	AuditConfig                  string // Der Pfad zur Audit-Konfiguration
	Config                       string // Der Pfad zur config.ini-Datei
	OutputPath                   string // Der Pfad, in welcher der Output-Ordner erstellt wird
	VerbosityLog                 int    // Das Verbositylevel der Log-Datei
	VerbosityConsole             int    // Das Verbositylevel der Konsolenausgaben
	SkipModuleCompatibilityCheck bool   // True, wenn alle Module unabhängig von Kompatibilität geladen werden sollen
	KeepConsoleOpen              bool   // True, wenn beim Doppelclicken auf die Executable das Konsolenfenster nach vollständigem Durchlauf offen bleiben soll
	ForceOS                      string // Mit diesem Parameter kann das Ergebnis des OS-Detectors überschrieben werden
	IgnoreMissingPrivileges      bool   // True, wenn das Programm nicht abbrechen soll, wenn es nicht mit den für die Audit-Konfiguration nötigen Privilegien gestartet wurde
	AlwaysPrintProgress          bool   // True, wenn der Fortschritt in den Modulen unabhängig vom Log-Level ausgegeben werden soll
	Zip                          bool   // True, wenn eine Zip-Datei, zusätzlich zum Output-Ordner erstellt werden soll
	ZipOnly                      bool   // True, wenn eine Zip Datei erstellt, der Output-Ordner aber entfernt werden soll

	// One-And-Done-Parameter (Programm führt keine Audits aus)
	Version             bool   // True, wenn nur der Version-String ausgegeben werden soll
	ShowModule          string // Hier wird ein Modulname übergeben, dessen Informationen ausgegeben werden sollen ("all" für alle Module)
	CheckConfiguration  bool   // True, wenn die Audit-Konfiguration geparsed, aber nicht ausgeführt werden soll
	CheckSyntax         bool   // True, wenn die Audit-Konfiguration ausschließlich auf ihre Syntax geprüft werden soll
	SaveConfiguration   bool   // True, wenn die aktuelle Konfiguration nach dem Parsen in die config.ini geschrieben werden soll
	CreateDefaultConfig bool   // True, wenn die Default-Configuration angelegt werden soll
}

type LogMsg struct {
	Message     string // Die Log-Nachricht
	Level       int    // Das Log-Level (0=None, Error, Warn, Info, 4=Debug)
	AlwaysPrint bool   // True, wenn die Nachricht unabhängig vom Log-Level ausgegeben werden soll
}
