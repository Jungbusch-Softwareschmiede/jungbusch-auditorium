package models

import "strconv"

// In diesem struct wird ein vollständiger Audit-Schritt, gemeinsam mit allen verschachtelten Modulen gespeichert.
type AuditModule struct {
	Condition                  string        // Die Bedingung, unter welcher das Modul ausgeführt wird.
	ModuleName                 string        // Der Name des Moduls, welches aufgerufen wird. Aliase werden vom Parser zum Name konvertiert.
	Description                string        // Beschreibung des Audit-Schritts.
	StepID                     string        // Eindeutige Identifikation des Audit-Schritts.
	Print                      string        // Variablen und Werte die nach Ausführung des Moduls ausgegeben werden.
	RequiresElevatedPrivileges bool          // Wert kommt aus dem Modul-Initializer. Kann aus der Audit-Konfiguraion überschrieben werden.
	PrivilegesOverwritten      bool          // True, wenn RequiresElevatedPrivileges aus der Audit-Konfiguration überschrieben wurde.
	Passed                     string        // Die Passed-Condition des Audit-Schritts.
	IsGlobal                   bool          // True, wenn der Schritt eine globale Variable enthält
	Variables                  VariableMap   // Eine Map aus allen models.Variable, auf die der Audit-Schritt Zugriff hat.
	ModuleParameters           ParameterMap  // Eine Map aus allen Parametern, die für das Modul in der Audit-Konfig angegeben wurde.
	NestedModules              []AuditModule // Alle verschachtelten Audit-Schritte
}

// Alias für [string]string Map
type ParameterMap map[string]string

// Alias für [string]Variable Map
type VariableMap map[string]Variable

// In diesem struct werden Variablen aus der Audit-Konfiguration gespeichert.
type Variable struct {
	Name     string
	Value    string
	IsEnv    bool
	IsGlobal bool
}

// In diesem struct werden Parameter gespeichert. Module bekommen Parameter aus der Audit-Konfiguration übergeben.
type Parameter struct {
	ParamName  string
	ParamValue string
}

// Syntax-Error-Objekte werden vom Audit-Konfigurations-Parser verwendet.
type SyntaxError struct {
	ErrorMsg     string // Die Error-Nachricht
	LineNo       int    // Die mögliche Zeilen-Nummern des Errors
	Line         string // Der Inhalt der Zeile des Errors
	Errorkeyword string // Das Schlüsselwort, an dem der Error aufgetreten ist, wenn vorhanden
	Err          error  // Ein herkömmlicher Error
}

// Diese Methode ist zuständig für das Parsen eines Syntax-Errors in einen String
func (e *SyntaxError) Error() string {
	if e.LineNo != -1 {
		return "Fehler-Nachricht: " + e.ErrorMsg + ", " + "Zeilen-Nummer: " + strconv.Itoa(e.LineNo) + ", Zeile: \"" + e.Line + "\""
	} else {
		return e.ErrorMsg
	}
}
