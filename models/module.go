package models

// Mithilfe dieses structs können die Module ihre eigenen Parameter setzen.
type ModuleSyntax struct {
	ModuleName                 string             // Der Name des Moduls
	ModuleAlias                []string           // Eine Liste aus allen Modul-Aliasen
	ModuleDescription          string             // Die Beschreibung des Moduls
	ModuleCompatibility        []string           // Eine Liste aus allen kompatiblen Betriebssystemen oder Wildcards
	RequiresElevatedPrivileges bool               // True, wenn das aktuelle Modul Administrator-/Root-Privilegien benötigt
	InputParams                ParameterSyntaxMap // Eine Map aus erwarteten Input-Parametern
}

// Alias für eine Map aus string und Parametersyntax
type ParameterSyntaxMap map[string]ParameterSyntax

// Mit diesem struct setzen die Module welche Parameter sie erwarten und in welcher Form
type ParameterSyntax struct {
	ParamName        string   // Der Name des Parameters
	ParamAlias       []string // Eine Liste mit allen Parameter-Aliasen
	ParamDescription string   // Die Beschreibung des Parameters
	IsOptional       bool     // True, wenn der Parameter nicht zwingend anzugeben ist
}

// In diesem struct sammelt ein Modul seine Ergebnisse
type ModuleResult struct {
	Artifacts []Artifact // Eine Liste aus allen Artefakten des Moduls
	Result    string     // Das Ergebnis des Moduls, beeinflusst vom Grep-Parameter
	ResultRaw string     // Das Ergebnis des Moduls, unverändert
	Err       error      // Ein Error, falls aufgetreten
}

// In diesem Datentyp werden von den Modulen generierte Artefakte gespeichert
type Artifact struct {
	Name   string // Der Name des Artefakts
	Value  string // Der Wert des Artefakts
	IsFile bool   `json:"-"` // True, wenn das Artefakt eine Datei ist
}
