// Dieses Package übernimmt die Kommunikation zwischen dem Rest des Frameworks und der Module.
// Das Framework ist von den Modulen vollständig getrennt, insofern dass die Methoden der Module nie explizit aus dem
// Code aufgerufen wird. Die einzige Referenz auf die Module befinden sich in der Slice static.Modules.
// Stattdessen wird Reflection verwendet um die Methoden der Module aufzurufen. Dafür ist es nötig, dass alle
// Modul-Methoden einheitliche Namen haben. Welche Methodennamen erwartet werden, wird von
// static.MODULE_INITIALIZER_SUFFIX, static.MODULE_EXECUTOR_SUFFIX und static.MODULE_VALIDATOR_SUFFIX.
// Die Methodennamen setzen sich durch den Modulname und den Suffix zusammen.
package modulecontroller

import (
	"errors"
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/modules"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"reflect"
	"regexp"
	s "strings"
)

var (
	loadedModules []ModuleSyntax // Speichert alle geladenen Module
)

// Initialisiert alle Module und überprüft währenddessen, ob die Module mit dem aktuellen OS kompatibel sind,
// und alle benötigten Methoden implementieren. Wenn ja, wird der Syntax des Moduls (models.ModuleSyntax)
// in die Slice loadedModules geschrieben. Mit `skipCompatibilityCheck` kann die Kompatibilitätsüberprüfung
// übersprungen werden.
func Initialize(skipCompatibilityCheck bool) ([]ModuleSyntax, error) {
	// Hier basierend auf dem os die ModuleSyntax laden
	loadedModules = make([]ModuleSyntax, 0)
	Info(SeperateTitle("Module-Intializer"))

	// Durch alle Module iterieren
	for n := range static.Modules {
		initExists := doesMethodExist(static.Modules[n] + static.MODULE_INITIALIZER_SUFFIX)
		executeExists := doesMethodExist(static.Modules[n] + static.MODULE_EXECUTOR_SUFFIX)

		switch {
		// Wenn nur eine der Methoden existiert, ist das Modul ungültig.
		// Wenn keine existiert, kann es sein dass das Modul aufgrund von GO-Buildtags nicht verfügbar ist,
		// daher wird dafür kein Error geworfen
		case (initExists && !executeExists) || (!initExists && executeExists):
			return []ModuleSyntax{}, errors.New("Dem Modul " + static.Modules[n] + " fehlt mindestens eine der benötigten Methoden. Es benötigt mindestens folgende Methoden: " + static.Modules[n] + static.MODULE_INITIALIZER_SUFFIX + ", " + static.Modules[n] + static.MODULE_EXECUTOR_SUFFIX)

		case initExists && executeExists:
			currModule := parseReflectToModuleSyntax(callFunctionFromString(static.Modules[n]+static.MODULE_INITIALIZER_SUFFIX, []reflect.Value{}))
			if currModule.ModuleName != static.Modules[n] {
				return nil, errors.New("Der in der Initialize-Methode eingetragene Modulname (\"" + currModule.ModuleName + "\") stimmt nicht mit dem Methodenname und dem in static.Modules eingetragenem Wert (\"" + static.Modules[n] + "\") überein.")
			}

			// Überprüfen, ob das Modul mit dem aktuellen Betriebssystem kompatibel ist.
			markedAsCompatible := util.ArrayContainsString(currModule.ModuleCompatibility, static.OperatingSystem)
			wildcardCompatibility := wildcardCompatibilityCheck(currModule.ModuleCompatibility)
			if skipCompatibilityCheck || markedAsCompatible || wildcardCompatibility {

				// Wenn es nicht kompatibel ist, aber trotzdem geladen werden soll, geben wir eine Warnung aus.
				if (!markedAsCompatible || !wildcardCompatibility) && skipCompatibilityCheck {
					Warn("Das Modul " + static.Modules[n] + " wurde nicht als kompatibel mit " + static.OperatingSystem + " markiert. Es wurde geladen, weil der Module-Kompatibilitäts-Check über einen Parameter deaktiviert wurde.")
				} else {
					Debug("Das Modul " + static.Modules[n] + " wurde erfolgreich geladen.")
				}

				// In die zurückzugebende Liste aller geladenen Module einfügen
				loadedModules = append(loadedModules, currModule)
			}
		}
	}

	return loadedModules, sanityCheck(loadedModules)
}

// Überprüft für alle übergeben models.AuditModule-Objekte ob das referenzierte Modul eine Validate-Methode hat.
// Wenn ja, wird diese mit den in der AuditKonfiguration angegebenen Parametern aufgerufen.
func ValidateAuditModules(mod []AuditModule) error {
	var err error
	for n := range mod {
		err = validateModule(mod[n])
		if err != nil {
			return err
		}
	}
	return nil
}

// Wird zum Aufrufen der Hauptmethode (Execute) eines Moduls verwendet werden. Diese Methode parsed
// die zu übergebenden Werte, sowie die Rückgabewerte in den korrekten Datentyp.
func Call(functionName string, m AuditModule) ModuleResult {
	// Das Skript-Modul ist das einzige Modul, welches eine custom Einbindung in den Rest des Frameworks hat.
	if s.ToLower(functionName) == "script" {
		return parseReflectToModuleResult(callFunctionFromString(functionName, parseVariableMapToReflect(m.ModuleParameters, &m.Variables)))
	}
	return parseReflectToModuleResult(callFunctionFromString(functionName, parseParameterSliceToReflect(m.ModuleParameters)))
}

// Wird zum Aufrufen der Validate-Methode eines Moduls verwendet werden. Diese Methode parsed
// die zu übergebenden Werte, sowie die Rückgabewerte in den korrekten Datentyp.
// Ist der zurückgegebene Error nil, sind alle Parameter valide.
func CallValidate(functionName string, params ParameterMap) error {
	return parseReflectToError(callFunctionFromString(functionName, parseParameterSliceToReflect(params)))
}

// Sucht den Modul-Syntax für den Übergebenen Modul-Name oder Modul-Alias und gibt diesen als formatierten String zurück.
func GetModuleSyntax(module string) string {
	module = s.ToLower(module)
	if module == "all" {
		var out string
		for n := range loadedModules {
			out += getModuleDescription(loadedModules[n])
		}
		return "\n" + out
	}

	for n := range loadedModules {
		if s.ToLower(loadedModules[n].ModuleName) == module || util.ArrayContainsString(loadedModules[n].ModuleAlias, module) {
			return "\n" + getModuleDescription(loadedModules[n])
		}
	}
	return "Modul konnte nicht gefunden werden."
}

//
//
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
// Utility
//
// *=*=*=*=*=*=*=*=*=*=*=*=
//
//
//

// Überprüft die Modul-Kompatibilität mit Wildcards gegeben ist
func wildcardCompatibilityCheck(moduleCompatibility []string) bool {
	for n := range moduleCompatibility {
		switch s.ToLower(moduleCompatibility[n]) {
		case "all":
			return true
		case "windows":
			return util.ArrayContainsString(static.Windows, static.OperatingSystem)
		case "linux":
			return util.ArrayContainsString(static.Linux, static.OperatingSystem)
		case "darwin":
			return util.ArrayContainsString(static.Darwin, static.OperatingSystem)
		}
	}
	return false
}

// Überprüft ob das im models.AuditModule-Objekt (also ein Audit-Schritt) aufgerufene Modul die Validate-Methode
// implementiert. Wenn ja, wird sie mit den gegebenen Parametern aufgerufen.
// Des Weiteren werden die Validate-Module von Kindern des übergebenen Audit-Schritts ausgeführt.
// Sind alle Parameter laut den Validate-Methoden valide wird nil zurückgegeben, ansonsten ein Error.
func validateModule(mod AuditModule) error {
	var err error

	if doesMethodExist(mod.ModuleName + static.MODULE_VALIDATOR_SUFFIX) {
		err = CallValidate(mod.ModuleName+static.MODULE_VALIDATOR_SUFFIX, mod.ModuleParameters)
		if err != nil {
			return err
		} else {
			for n := range mod.NestedModules {
				err = validateModule(mod.NestedModules[n])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Formatiert den models.ModuleSyntax eines Moduls in einen string.
func getModuleDescription(module ModuleSyntax) string {
	var out string

	out = fmt.Sprintf("%-25v", "Modul: ") + "\n"
	out += fmt.Sprintf("%-25v", "   Name, Alias: ") + module.ModuleName + ", " + util.PrintStrArray(module.ModuleAlias) + "\n"
	if s.HasSuffix(out, ", \n") {
		out = out[:len(out)-3] + "\n"
	}
	out += fmt.Sprintf("%-25v", "   Beschreibung: ") + fmt.Sprint(module.ModuleDescription) + "\n"
	out += fmt.Sprintf("%-25v", "   Kompatibilität: ") + util.PrintStrArray(module.ModuleCompatibility) + "\n"
	out += fmt.Sprintf("%-25v", "   Parameter: ") + "\n"

	for n := range module.InputParams {
		out += fmt.Sprintf("%-25v", "      Name, Alias: ") + module.InputParams[n].ParamName + ", " + util.PrintStrArray(module.InputParams[n].ParamAlias) + "\n"
		if s.HasSuffix(out, ", \n") {
			out = out[:len(out)-3] + "\n"
		}
		out += fmt.Sprintf("%-25v", "      Beschreibung: ") + module.InputParams[n].ParamDescription + "\n"
		out += fmt.Sprintf("%-25v", "      Optional: ")
		if module.InputParams[n].IsOptional {
			out += "Ja\n\n"
		} else {
			out += "Nein\n\n"
		}
	}
	return out
}

// Parsed den Rückgabewert eines Calls zurück in models.ModuleSyntax.
func parseReflectToModuleSyntax(in []reflect.Value) ModuleSyntax {
	return in[0].Interface().(ModuleSyntax)
}

// Parsed den Rückgabewert eines Calls in ein Error-Struct.
func parseReflectToError(in []reflect.Value) error {
	if in[0].Interface() != nil {
		return in[0].Interface().(error)

	} else {
		return nil
	}
}

// Parsed den Rückgabewert eines Calls in ein models.ModuleResult-Struct.
func parseReflectToModuleResult(in []reflect.Value) ModuleResult {
	return in[0].Interface().(ModuleResult)
}

// Parsed eine models.ParameterMap in mit Reflect zu verwendenden Reflect-Values.
func parseParameterSliceToReflect(in ParameterMap) []reflect.Value {
	return []reflect.Value{reflect.ValueOf(in)}
}

// Parsed eine models.ParameterMap, sowie models.VariableMap in mit Reflect zu verwendenden Reflect-Values.
func parseVariableMapToReflect(pm ParameterMap, vm *VariableMap) []reflect.Value {
	return []reflect.Value{parseParameterSliceToReflect(pm)[0], reflect.ValueOf(vm)}
}

// Überprüft ob eine Methode mit dem übergebenen Name in Sichtweite des modules.MethodHandler-Objekts existiert.
func doesMethodExist(function string) bool {
	st := reflect.TypeOf((*modules.MethodHandler)(nil))
	_, exists := st.MethodByName(function)
	return exists
}

// Ruft die Methode mit dem übergebenen Name und den übergebenen Übergabewerten auf. Gibt Rückgabewerte in einer Slice
// im Format reflect.Value zurück, die zurück in den richtigen Datentyp gecasted werden muss.
func callFunctionFromString(function string, values []reflect.Value) []reflect.Value {
	mh := new(modules.MethodHandler)
	method := reflect.ValueOf(mh).MethodByName(function)
	response := method.Call(values)
	return response
}

// Diese Methode überprüft, ob sich die in den Initialize-Methoden von Modulen gesetzten Namen und Aliase überschneiden.
// Modulnamen müssen einzigarig sein. Ist das nicht der Fall, wird ein passender Error geworfen.
func sanityCheck(mod []ModuleSyntax) error {
	// Modulnamen checken
	usedNames := make(map[string]bool, 0)
	for _, m := range mod {
		if hasDisallowedCharacters(m.ModuleName) {
			return errors.New("Der Modul-Name enthält ungültige Zeichen: " + m.ModuleName)
		}

		if s.TrimSpace(s.Trim(m.ModuleName, "\t")) == "" {
			return errors.New("Der Modul-Name darf nicht ein leerer String oder ein String bestehend aus ausschließlich Leerzeichen sein: " + m.ModuleName)
		}

		syn := acutil.GetParameterSyntaxFromKeyword(m.ModuleName)
		if syn.ParamName != "" {
			return errors.New("Der Modul-Name " + m.ModuleName + " überschneidet sich mit dem Jungbusch-Auditorium-Parameter " + syn.ParamName + " (" + syn.ParamDescription + " ). Bitte den Modulnamen ändern.")
		}

		if usedNames[m.ModuleName] {
			return errors.New("Das Modul " + m.ModuleName + " kommt mehrmals vor. Modulnamen müssen einzigartig sein und dürfen sich nicht mit den Namen/Aliasen anderer Module überschneiden.")
		} else {
			usedNames[m.ModuleName] = true
		}

		for _, sm := range m.ModuleAlias {
			if hasDisallowedCharacters(sm) {
				return errors.New("Der Modul-Alias enthält ungültige Zeichen: " + sm)
			}

			if s.TrimSpace(s.Trim(sm, "\t")) == "" {
				return errors.New("Modul-Aliase dürfen nicht ein leerer String oder ein String bestehend aus ausschließlich Leerzeichen sein: " + sm)
			}

			syn = acutil.GetParameterSyntaxFromKeyword(sm)
			if syn.ParamName != "" {
				return errors.New("Der Modul-Alias " + sm + " des Moduls " + m.ModuleName + " überschneidet sich mit dem Jungbusch-Auditorium-Parameter " + syn.ParamName + " (" + syn.ParamDescription + " ). Bitte den Alias ändern oder entfernen.")
			}

			if usedNames[sm] {
				return errors.New("Das Modul " + sm + " kommt mehrmals vor. Modulnamen müssen einzigartig sein und dürfen sich nicht mit den Namen/Aliasen anderer Module überschneiden.")
			} else {
				usedNames[sm] = true
			}
		}
	}

	// Modulparameter checken
	for _, m := range mod {
		usedNames = make(map[string]bool, 0)
		for _, mp := range m.InputParams {
			if hasDisallowedCharacters(mp.ParamName) {
				return errors.New("Der Parameter-Name enthält ungültige Zeichen: " + mp.ParamName)
			}

			if s.TrimSpace(s.Trim(mp.ParamName, "\t")) == "" {
				return errors.New("Parameter-Namen dürfen nicht ein leerer String oder ein String bestehend aus ausschließlich Leerzeichen sein: " + mp.ParamName)
			}

			syn := acutil.GetParameterSyntaxFromKeyword(mp.ParamName)
			if syn.ParamName != "" {
				return errors.New("Der Parameter-Name " + mp.ParamName + " des Moduls " + m.ModuleName + " überschneidet sich mit dem Jungbusch-Auditorium-Parameter " + syn.ParamName + " (" + syn.ParamDescription + " ). Bitte den Parametername ändern.")
			}

			usedNames[mp.ParamName] = true
			for _, mpa := range mp.ParamAlias {
				if hasDisallowedCharacters(mpa) {
					return errors.New("Der Parameter-Alias enthält ungültige Zeichen: " + mpa)
				}

				if s.TrimSpace(s.Trim(mpa, "\t")) == "" {
					return errors.New("Parameter-Aliase dürfen nicht ein leerer String oder ein String bestehend aus ausschließlich Leerzeichen sein: " + mpa)
				}

				syn = acutil.GetParameterSyntaxFromKeyword(mpa)
				if syn.ParamName != "" {
					return errors.New("Der Parameter-Alias " + mpa + " des Parameters " + mp.ParamName + " des Moduls " + m.ModuleName + " überschneidet sich mit dem Jungbusch-Auditorium-Parameter " + syn.ParamName + " (" + syn.ParamDescription + " ). Bitte den Parametername ändern.")
				}

				if usedNames[mpa] {
					return errors.New("Der Parametername " + mpa + " kommt im Modul " + m.ModuleName + " mehrmals vor. Parameternamen/Aliase sollten sich nicht überschneiden.")
				} else {
					usedNames[mpa] = true
				}
			}
		}
	}

	return nil
}

// True, wenn im übergebenen String für einen Modulname ungültige Zeichen enthalten sind.
func hasDisallowedCharacters(in string) bool {
	match, _ := regexp.MatchString(static.REG_MODULE_NAME, in)
	return !match
}
