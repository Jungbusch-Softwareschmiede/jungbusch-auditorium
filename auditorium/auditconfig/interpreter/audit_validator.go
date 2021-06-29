package interpreter

import (
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/modulecontroller"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/pkg/errors"
	"strings"
)

// Dient als "Start-Funktion". Ruft validateParametersInModule für alle Level-0 Module auf
func ValidateParameters(mods []AuditModule) error {

	// Variablen validieren
	if err := validateVariables(mods); err != nil {
		return err
	}

	// Parameter validieren
	if err := modulecontroller.ValidateAuditModules(mods); err != nil {
		return err
	}

	return nil
}

// Überprüft, ob alle referenzierten Variablen vor Verwendung deklariert wurden
func validateVariables(modules []AuditModule) error {
	for _, m := range modules {
		// Alle Variablen in Print, Passed, Condition und Modul-Parametern finden
		vars := acutil.GetVariablesInString(m.Print)
		vars = append(vars, acutil.GetVariablesInString(m.Condition)...)
		vars = append(vars, acutil.GetVariablesInString(m.Passed)...)

		for _, param := range m.ModuleParameters {
			vars = append(vars, acutil.GetVariablesInString(param)...)
		}

		// Überprüfen, ob alle gefundenen Variablen deklariert wurden
		for _, v := range vars {
			if _, ok := m.Variables[strings.ToLower(v)]; !strings.HasPrefix(v, "%g_") && !ok {
				return errors.New("Variable " + strings.ToLower(v) + " wurde nicht deklariert!")
			}
		}

		for _, v := range m.Variables {
			// Umgebungsvariablen haben immer einen Wert
			if !v.IsEnv {

				// Wenn die aktuelle Variable auf eine unbekannte Variable verweist, ist das ungültig
				v.Value = strings.ToLower(v.Value)
				if _, ok := m.Variables[v.Value]; acutil.IsVariable(v.Value) && !ok {
					return errors.New("Die Variable " + v.Name + " verweist auf eine unbekannte Variable: " + v.Value)
				}

				// Setzen der Variable in Kind-Modulen
				for _, nm := range m.NestedModules {
					nm.Variables[v.Name] = m.Variables[v.Name]
				}
			}
		}
		if err := validateVariables(m.NestedModules); err != nil {
			return err
		}
	}
	return nil
}
