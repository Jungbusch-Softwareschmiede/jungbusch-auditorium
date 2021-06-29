// Dieses Package ist für das Einlesen der Auditkonfigurations-Datei in models.AuditModule -Objekte zuständig.
// Sie ist darauf angewiesen, dass vorher der Syntaxchecker erfolgreich durchgelaufen ist und somot alle
// möglichen Syntaxerror abgefangen wurden. Ist dies nicht der Fall und enthält die Auditkonfigurationsdatei
// Fehler, dann ist das Auftreten eines Panics in diesem Package wahrscheinlich.
package parser

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	s "strings"
)

var (
	processRequiresElevatedPrivileges bool            // Speichert ob der Prozess Administratorprivilegien benötigt.
	numberOfModules                   int             // Enthält die gesamte Anzahl der eingelesenen Module. Wird für die Progressbar benötigt.
	usedIds                           map[string]bool // Mit dieser Map wird gespeichert, welche IDs für die Module bereits verwendet wurde.
	loadedModules                     []ModuleSyntax  // In dieser Map werden die geladenen Modulen für den Zugriff aus dem gesamten Package gespeichert.
)

// Liest die übergebene Auditkonfigurations-Datei anhand der übergebenen geladenen Modulen ein.
// Gibt eine Slice aller Audit-Schritte, ob der Prozess Administratorberechtigungen benötigt, die gesamte
// Anzahl der Module und gegebenenfalls einen Error zurück.
func Parse(lines []string, moduleSyntax []ModuleSyntax) ([]AuditModule, bool, int, error) {
	modules := make([]AuditModule, 0)
	usedIds = make(map[string]bool, 0)
	n := new(int)
	loadedModules = moduleSyntax
	var err error
	var l string
	var currentModule AuditModule

	// Die zu verwendenden Zeilen für die Utility setzen
	SetLines(lines)

	// Alle Zeilen iterieren
	for ; *n < len(lines); *n++ {
		// Kommentare, leere Zeilen überspringen
		if err = SkipIrrelevantLines(lines, n); err != nil {
			return nil, false, 0, err
		}

		if LinesFinished(n) {
			break
		}

		// Entfernt den Inline-Kommentar aus der aktuellen Zeile, wenn vorhanden
		l, err = RemoveInlineComment(Trim(lines[*n]))
		if err != nil {
			return nil, false, 0, GenerateSyntaxError(err.Error(), *n, l, "")
		}

		// An dieser Stelle muss jetzt ein Modul geöffnet werden
		if l == "{" {
			currentModule, err = generateModule(lines, n)
			if err != nil {
				return nil, false, 0, err
			}
			modules = append(modules, currentModule)
			numberOfModules++
		} else {
			return nil, false, 0, GenerateSyntaxError(static.MODULE_MISSING_OPENING_BRACKET, *n, l, "")
		}
	}

	return modules, processRequiresElevatedPrivileges, numberOfModules, validateGlobals(modules)
}

// Generiert ein einzelnes Modul. Wird rekursiv aufgerufen.
func generateModule(lines []string, n *int) (AuditModule, error) {
	var err error
	*n++

	// Ein neues Audit-Modul mit einigen Default-Werten generieren, welches jetzt gefüllt werden soll.
	currentModule := initializeAuditModule()

	if LinesFinished(n) {
		return AuditModule{}, GenerateSyntaxError(static.MODULE_EMPTY, *n, lines[len(lines)-1], "")
	}
	l := Trim(lines[*n])

	// Zeilen iterieren bis das aktuelle Modul endet
	for err == nil && l != "}," {
		// Überprüfen, ob die Datei zu Ende ist
		if LinesFinished(n) {
			return AuditModule{}, EvaluateClosingBracketError(n)
		}

		// Kommentare, leere Zeilen überspringen
		err = SkipIrrelevantLines(lines, n)
		l, err = PrepareLine(lines[*n])
		if err != nil {
			return AuditModule{}, GenerateSyntaxError(err.Error(), *n, l, "")
		}

		// Fehlendes Komma bei der schließenden Klammer abfangen
		if l == "}" {
			return AuditModule{}, EvaluateClosingBracketError(n)
		}

		// Verschachteltes Modul
		if l == "{" {
			mod, err := generateModule(lines, n)
			if err != nil {
				return AuditModule{}, err
			}
			currentModule.NestedModules = append(currentModule.NestedModules, mod)
			numberOfModules++
			*n++
			l, err = PrepareLine(lines[*n])
			continue
		}

		// Die Zeile enthält einen Parameter
		if IsParameter(n) {
			if err = handleParameter(lines, n, &currentModule); err != nil {
				return AuditModule{}, err
			}
			*n++
			l, err = PrepareLine(lines[*n])
			continue
		} else

		// Die Zeile enthält eine Variable
		if IsVariableAssignment(n) {
			if err = handleVariable(lines, n, &currentModule); err != nil {
				return AuditModule{}, err
			}
			*n++
			l, err = PrepareLine(lines[*n])
			continue
		}

		// Ungültiger Ausdruck
		return AuditModule{}, GenerateSyntaxError(static.MODULE_INVALID_EXPRESSION, *n, l, "")
	}

	return currentModule, validateCompletenessOfModule(&currentModule, lines, n)
}

// Überprüft, ob alle nötigen Werte im übergebenen Schritt gesetzt sind.
func validateCompletenessOfModule(module *AuditModule, lines []string, n *int) error {

	// Checken ob zwingend zu setzende Parameter fehlen
	missingParameters := ""
	if len(s.Trim(s.TrimSpace(Trim(module.ModuleName)), "\t")) == 0 {
		missingParameters += "Modul-Name "
	}

	if len(s.Trim(s.TrimSpace(Trim(module.StepID)), "\t")) == 0 {
		missingParameters += "Step-ID "
	}

	if missingParameters != "" {
		return GenerateSyntaxError(static.MODULE_MISSING_PARAMETER+missingParameters+". "+static.ROW_IS_END_OF_MODULE, *n, lines[*n], "")
	}

	// Den Blueprint des aktuellen Moduls finden
	moduleSyntax := GetModuleSyntaxFromNameOrAlias(module.ModuleName, loadedModules)

	// Checken ob der Modulname/ das Modul existiert
	if moduleSyntax.ModuleName == "" {
		return GenerateSyntaxError(static.MODULE_NOT_FOUND+static.ROW_IS_END_OF_MODULE, *n, lines[*n], "")
	}

	// Falls ein Alias gesetzt war, mit dem Modulname austauschen
	module.ModuleName = moduleSyntax.ModuleName

	// Umgebungsvariable für das Modul setzen
	cm := module.Variables["%currentmodule%"]
	cm.Value = module.ModuleName
	module.Variables["%currentmodule%"] = cm

	// Die Keys aus der Modul-Parametermap extrahieren
	keys := reflect.ValueOf(module.ModuleParameters).MapKeys()

	// Alle Modul-Parameter-Aliase zu den korrekten Modul-Parameter-Namen konvertieren
	for _, v := range keys {
		paramAlias := v.String()
		param := GetModuleParameterSyntaxNameFromAlias(moduleSyntax, paramAlias)

		if param == "" {
			return GenerateSyntaxError(static.MODULE_INVALID_PARAMETER+static.ROW_IS_END_OF_MODULE, *n, lines[*n], paramAlias)
		}

		// Wenn in der Audit-Konfiguration für den Modulparameter ein Alias angegeben war,
		// setzen wir den Wert in der Map an die Stelle des Parameternamens
		if param != paramAlias {
			module.ModuleParameters[param] = module.ModuleParameters[paramAlias]
			module.ModuleParameters[paramAlias] = ""
		}
	}

	// Überprüfen ob alle zwingend benötigten Modul-Parameter gesetzt wurden und ob sie Multiline sind oder nicht
	for _, m := range moduleSyntax.InputParams {
		// Wenn Parameter nicht optional, checken ob gesetzt
		if !m.IsOptional {
			if module.ModuleParameters[m.ParamName] == "" {
				return GenerateSyntaxError(static.MODULE_REQUIRED_MODULE_PARAMETER_NOT_SET+m.ParamName+". "+static.ROW_IS_END_OF_MODULE, *n, lines[*n], m.ParamName)
			}
		}
	}

	// Wenn die Privilegien in der Audit-Konfig nicht überschrieben wurden, holen wir uns sie aus dem Modul-Syntax
	if !module.PrivilegesOverwritten {
		module.RequiresElevatedPrivileges = moduleSyntax.RequiresElevatedPrivileges
		processRequiresElevatedPrivileges = module.RequiresElevatedPrivileges || processRequiresElevatedPrivileges
	}

	// Wenn Passed nicht gesetzt wurde, ist es true by default
	if module.Passed == "" {
		module.Passed = "true"
	}

	return nil
}

// Übernimmt das extrahieren, konvertieren und validieren eines Parameters aus der übergebenen Zeile.
func handleParameter(lines []string, n *int, currentModule *AuditModule) error {
	// Error wurde im Syntaxcheck abgefangen
	line, _ := RemoveInlineComment(lines[*n])
	name, value := SplitParameter(line)

	// Multiline behandeln
	if s.HasPrefix(value, "`") {

		if s.Index(value, "`") == s.LastIndex(value, "`") {
			*n++
			line, _ = PrepareLine(lines[*n])
			value = value + "\n" + line
		}

		for ; !s.HasSuffix(line, "`"); {
			*n++
			line, _ = RemoveInlineComment(lines[*n])
			value = value + "\n" + line
		}
	}
	value = VariablesInStringToLower(value)

	p := Parameter{ParamName: name, ParamValue: value}
	pSyntax := GetParameterSyntaxFromKeyword(p.ParamName)

	// Durch alle allgemeinbekannten Syntaxe iterieren und den Wert validieren
	switch s.ToLower(pSyntax.ParamName) {

	case "condition":
		condition := getConditionStringFromString(p.ParamValue)
		if err := setValue(&currentModule.Condition, condition); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], "")
		}

	case "module":
		if err := setValue(&currentModule.ModuleName, p.ParamValue); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], "")
		}

	case "description":
		if err := setValue(&currentModule.Description, p.ParamValue); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], "")
		}

	case "passed":
		condition := getConditionStringFromString(p.ParamValue)
		if err := setValue(&currentModule.Passed, condition); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], "")
		}

	case "stepid":
		if !usedIds[p.ParamValue] {
			if err := setValue(&currentModule.StepID, p.ParamValue); err != nil {
				return GenerateSyntaxError(err.Error(), *n, lines[*n], "")
			}
			usedIds[p.ParamValue] = true
		} else {
			return GenerateSyntaxError(static.MODULE_VALUE_INVALID+static.ID_ALREADY_USED, *n, lines[*n], p.ParamValue)
		}

	case "requireselevatedprivileges":
		if !currentModule.PrivilegesOverwritten {
			currentModule.PrivilegesOverwritten = true
			processRequiresElevatedPrivileges = true

			boolean, err := util.ParseStringToBool(s.Trim(p.ParamValue, "\""))
			if err != nil {
				return GenerateSyntaxError(static.MODULE_VALUE_INVALID+static.INVALID_VALUE_FOR_BOOL, *n, lines[*n], p.ParamValue)
			}
			currentModule.RequiresElevatedPrivileges = boolean
		} else {
			return GenerateSyntaxError(static.MODULE_VALUE_ALREADY_SET, *n, lines[*n], "")
		}

	case "print":
		if err := setValue(&currentModule.Print, p.ParamValue); err != nil {
			return err
		}

	default:
		// Es liegt ein Modul-Parameter vor
		if s.HasPrefix(p.ParamValue, "`") && s.HasSuffix(p.ParamValue, "`") {
			p.ParamValue = s.Trim(p.ParamValue, "`")
		} else if s.HasPrefix(p.ParamValue, "\"") && s.HasSuffix(p.ParamValue, "\"") {
			p.ParamValue = s.Trim(p.ParamValue, "\"")
		}

		if currentModule.ModuleParameters[p.ParamName] != "" {
			return GenerateSyntaxError(static.MODULE_VALUE_ALREADY_SET, *n, lines[*n], "")
		}
		currentModule.ModuleParameters[p.ParamName] = p.ParamValue
	}

	return nil
}

// Übernimmt das extrahieren, konvertieren und validieren einer Variable aus der übergebenen Zeile.
func handleVariable(lines []string, n *int, currentModule *AuditModule) error {
	// Error wurde im Syntaxcheck abgefangen
	line, _ := RemoveInlineComment(lines[*n])
	name, value := SplitVariable(line)

	// Variablen sind nicht Case-Sensitiv, daher mache ich alle Variablen lowercase
	name = s.ToLower(name)
	if IsVariable(value) {
		value = s.ToLower(value)
	} else {
		value = s.Trim(value, "\"")
	}

	// Umgebungsvariablen dürfen nicht überschrieben werden
	if currentModule.Variables[name].IsEnv {
		return GenerateSyntaxError(static.VARIABLE_INVALID_NAME+static.ENVIRONMENT_VARIABLES_READONLY, *n, lines[*n], name)
	}

	// Variablen dürfen pro Modul nur einmal gesetzt werden
	if currentModule.Variables[name].Name != "" {
		return GenerateSyntaxError(static.VARIABLE_INVALID+static.VARIABLE_ALREADY_SET, *n, lines[*n], name)
	}

	isGlobal := s.HasPrefix(name, "%g_")
	if isGlobal {
		currentModule.IsGlobal = true
	}

	currentModule.Variables[name] = Variable{
		Name:     name,
		Value:    value,
		IsGlobal: isGlobal,
	}

	return nil
}

// Validiert globale Variablen, in allen übergebenen Auditschritten und überprüft ob sie an der richtigen Stelle
// deklariert sind und dem geforderten Syntax entsprechen.
func validateGlobals(mods []AuditModule) error {
	prev := mods[0]

	for _, m := range mods {
		if m.IsGlobal {

			// Globale Module dürfen keine Verschachtelung haben
			if len(m.NestedModules) != 0 {
				return errors.New(static.GLOBAL_MODULES_NOT_ALLOWED_NESTED + m.StepID)
			}

			if !prev.IsGlobal {
				// Das aktuelle Modul ist Global, das vorherige nicht
				return errors.New(static.GLOBAL_VARIABLE_READONLY + m.StepID)
			}
		} else {
			for _, nest := range m.NestedModules {
				if nest.IsGlobal {
					return errors.New(static.GLOBAL_MODULE_NESTED + m.StepID)
				}
			}
		}
		prev = m
	}

	return nil
}

// Setzt den übergebenen Wert in den übergebenen Pointer und überprüft ob in den Pointer bereits etwas gesetzt wurde.
func setValue(pointerToSetIn *string, valueToSet string) error {
	if *pointerToSetIn == "" {

		if s.HasPrefix(valueToSet, "\"") && s.HasSuffix(valueToSet, "\"") {
			*pointerToSetIn = Trim(s.Trim(valueToSet, "\""))
		} else if s.HasPrefix(valueToSet, "`") && s.HasSuffix(valueToSet, "`") {
			*pointerToSetIn = Trim(s.Trim(valueToSet, "`"))
		} else {
			*pointerToSetIn = valueToSet
		}
	} else {
		return errors.New(static.MODULE_VALUE_ALREADY_SET)
	}
	return nil
}

// Returned den Condition-String aus dem übergebenen String. Klammern und Anführungszeichen werden entfernt.
func getConditionStringFromString(condition string) string {
	// Errors hat der Syntaxchecker abgefangen
	ifreg, _ := regexp.Compile("(\\()( *)([\"`])(.*)([\"`])( *)(\\))")
	result := ifreg.FindString(condition)
	return result[1 : len(result)-1]
}

// Diese Methode initialisiert AuditModules mit ihren im Interpreter zu setzenden Environment-Variables
func initializeAuditModule() AuditModule {
	varSlice := make(VariableMap)
	varSlice["%result%"] = Variable{Name: "%result%", IsEnv: true}
	varSlice["%passed%"] = Variable{Name: "%passed%", IsEnv: true}
	varSlice["%unsuccessful%"] = Variable{Name: "%unsuccessful%", IsEnv: true}
	varSlice["%os%"] = Variable{Name: "%os%", Value: static.OperatingSystem, IsEnv: true}
	varSlice["%currentmodule%"] = Variable{Name: "%currentmodule%", IsEnv: true}
	return AuditModule{Variables: varSlice, ModuleParameters: ParameterMap{}}
}
