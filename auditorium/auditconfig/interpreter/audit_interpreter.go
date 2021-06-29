package interpreter

import (
	"errors"
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/modulecontroller"
	output "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/outputgenerator"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/dop251/goja"
	"github.com/schollz/progressbar/v3"
	"os"
	"strconv"
	"strings"
)

// In diesem Struct werden die in Audit-Schritten verwendeten Variablen gespeichert. Dafür muss der originale Name
// (also mit %-Zeichen etc.) für später gespeichert werden und ein neuer Name, der dem ECMAScript-Syntax entspricht,
// erzeugt werden.
type moduleVariable struct {
	name     string
	origName string
	value    string
}

// Diese Funktion dient als eine Art Wrapper für das Interpretieren der Audit-Konfiguration.
// Es wird die Gültigkeit aller Variablen und Parameter überprüft, auch ob jeder Variable früher oder später
// ein Wert zugewiesen wird. Daraufhin werden alle Module ausgeführt und die Ergebnisse in einer Slice
// aus output.ReportEntry-Objekten zurückgegeben.
func InterpretAudit(audit []AuditModule, numberOfModules int, alwaysPrintProgress bool) (report []output.ReportEntry, err error) {

	// Parameter validieren, wenn irgendwas nicht stimmt, returnen
	err = ValidateParameters(audit)
	if err != nil {
		return
	} else {
		Info("Alle Variablen und Modulparameter sind gültig.")
	}

	Info(SeperateTitle("Audit-Interpreter"))
	static.ProgressBar = progressbar.NewOptions(0,
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionClearOnFinish())
	static.ProgressBar.ChangeMax(numberOfModules)

	report = executeModules(audit, numberOfModules, alwaysPrintProgress)
	return
}

// executeModules führt die in der Audit-Datei festegelegten Module aus und gibt einen Report zurück.
// Zuächst wird überprüft, ob die zum Ausführen des Moduls benötigten Berechtigungen vorliegen.
// Anschließend wird die Ausführ-Bedingung überprüft.
// Ist diese erfüllt, so wird das Modul ausgeführt und das Ergebnis in ein output.ReportEntry-Objekt geschrieben.
func executeModules(modules []AuditModule, numberOfModules int, alwaysPrintProgress bool) []output.ReportEntry {
	report := make([]output.ReportEntry, len(modules))
	globalVars := make([]Variable, 0)

	for i, m := range modules {
		if m.Description == "" {
			m.Description = fmt.Sprintf("%v-%v", m.StepID, m.ModuleName)
		}

		// Anlegen des Report-Eintrags
		report[i] = output.ReportEntry{
			ID:   m.StepID,
			Desc: m.Description,
		}

		// Das Modul kann aufgrund von fehlender Berechtigungen nicht ausgeführt werden
		if m.RequiresElevatedPrivileges && !static.HasElevatedPrivileges {
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, errors.New("Zur Ausführung des Moduls werden Administrator, bzw. Root-Privilegien benötigt."))
			continue
		}

		// Setzen der globalen Variablen
		for _, v := range globalVars {
			m.Variables[v.Name] = v
		}

		// Ausführbedingung überprüfen
		ok, err := checkExecuteCondition(m)
		if err != nil {
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, err)
			continue
		}

		// Wenn die Ausführbedingung nicht erfüllt ist, skippen
		if !ok {
			report[i].Result = "NOTEXECUTED"
			Info(fmt.Sprintf(static.NOTEXECUTED_OUTPUT, m.Description))
			Debug(Seperate())
			_ = static.ProgressBar.Add(1)
			addNotExecutedToProgressBar(m)
			continue
		}

		// Ersetzen von möglichen Variablen in den Parametern
		for k, v := range m.ModuleParameters {
			m.ModuleParameters[k] = replaceVariablesInString(v, m.Variables)
		}

		// Ausführ-Bedingung ist erfüllt
		Debug(fmt.Sprintf("Parameter: %v", m.ModuleParameters))

		// Ausführen des Moduls und Erhalten des Ergebnisses
		res, err := execute(m)
		if err != nil {
			// Wenn beim Ausführen des Moduls ein Fehler auftritt
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, err)
			continue
		}
		Debug(fmt.Sprintf("Erhaltenes Ergebnis: %v", res.Result))

		// Setzen der Ergebnis-Variable
		m.Variables["%result%"] = Variable{
			Name:  "%result%",
			Value: res.Result,
			IsEnv: true,
		}

		// Weitergeben der Artefakte
		report[i].Artifacts = res.Artifacts

		// Überprüfen der Passed-Bedingung
		if err = setPassed(&m); err != nil {
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, err)
			continue
		}

		// Die Variablen für nachfolgende Module setzen
		if err = setVariables(&m, &globalVars); err != nil {
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, err)
			continue
		}

		// Wenn der Print-Parameter nicht leer ist, ersetzen wir alle enthaltenen Variablen und geben ihn aus
		if m.Print != "" {
			m.Print = replaceVariablesInString(m.Print, m.Variables)
			InfoPrintAlways(m.Print)
			report[i].Print = m.Print
		}

		// Ergebnis der Passed-Bedingung überprüfen, in den Report eintragen und Konsolenausgaben
		passed, err := strconv.ParseBool(m.Variables["%passed%"].Value)
		if err != nil {
			handleUnsuccessful(m, &report[i], alwaysPrintProgress, err)
			continue
		}

		if passed {
			// Bedingung ist erfüllt
			InfoPrint(fmt.Sprintf(static.PASSED_OUTPUT, m.Description), alwaysPrintProgress)
			Debug(Seperate())
			_ = static.ProgressBar.Add(1)
			report[i].Result = "PASSED"
		} else {
			// Bedingung ist nicht erfüllt
			InfoPrint(fmt.Sprintf(static.NOTPASSED_OUTPUT, m.Description), alwaysPrintProgress)
			Debug(Seperate())
			_ = static.ProgressBar.Add(1)
			report[i].Result = "NOTPASSED"
			report[i].Expected = m.Passed
			report[i].Actual = m.Variables["%result%"].Value
			if report[i].Actual == "" {
				report[i].Actual = "null"
			}
		}

		// Rekursives Ausführen der verschachtelten Module
		report[i].Nested = executeModules(m.NestedModules, numberOfModules, alwaysPrintProgress)
	}

	return report
}

// Übernimmt das Zusammenbauen eines Report-Entries für einen Audit-Schritt, in dem ein Fehler aufgetreten ist.
func handleUnsuccessful(m AuditModule, reportEntry *output.ReportEntry, alwaysPrintProgress bool, err error) {
	InfoPrint(fmt.Sprintf(static.UNSUCCESSFUL_OUTPUT, m.Description), alwaysPrintProgress)
	reportEntry.Result = "UNSUCCESSFUL"
	reportEntry.Expected = m.Passed
	reportEntry.Actual = m.Variables["%result%"].Value
	if reportEntry.Actual == "" {
		reportEntry.Actual = "null"
	}
	reportEntry.Error = err.Error()
	m.Variables["%unsuccessful%"] = Variable{
		Name:  "%unsuccessful%",
		Value: "true",
		IsEnv: true,
	}

	Debug("Modul fehlgeschlagen: " + err.Error())
	Debug(Seperate())
	_ = static.ProgressBar.Add(1)
	addNotExecutedToProgressBar(m)
}

// Überprüft die Ausführ-Bedingung, des Moduls, wenn angegeben
func checkExecuteCondition(m AuditModule) (bool, error) {
	if m.Condition != "" {
		ok, err := interpretCondition(m.Condition, m.Variables)
		if err != nil {
			return false, err
		}
		Debug(fmt.Sprintf(`Ergebnis der Ausführ-Bedingung "%v": %v`, m.Condition, ok))
		return ok, nil
	} else {
		return true, nil
	}
}

// Führt die Execute-Methode des übergebenen Moduls aus.
func execute(m AuditModule) (res ModuleResult, err error) {
	res = modulecontroller.Call(m.ModuleName, m)
	if res.Err != nil {
		return ModuleResult{Artifacts: res.Artifacts}, res.Err
	}
	res.Result = strings.TrimSuffix(res.Result, "\n")
	return res, nil
}

// Interpretiert die Passed-Condition, macht Konsolenausgaben und setzt Umgebungsvariablen
func setPassed(m *AuditModule) error {
	passed, err := interpretCondition(m.Passed, m.Variables)
	if err != nil {
		return err
	}
	Debug(fmt.Sprintf(`Erhaltenes Ergebnis für "%v": %v`, m.Passed, passed))
	m.Variables["%passed%"] = Variable{
		Name:  "%passed%",
		Value: strconv.FormatBool(passed),
		IsEnv: true,
	}
	return nil
}

// Iteriert über die Variablen des übergebenen Moduls. Die Werte aller Variablen, die keine Umgebungsvariablen sind,
// werden an die Kinder des Moduls weitergegeben, da Variablen des Elternteils dort ebenfalls zugreifbar sind.
func setVariables(m *AuditModule, globalVars *[]Variable) error {
	for _, v := range m.Variables {
		if !v.IsEnv {
			// Bezieht sich die Variable auf eine andere wird hier ihr Wert korrekt gesetzt
			if _, ok := m.Variables[v.Value]; ok {
				m.Variables[v.Name] = Variable{
					Name:  v.Name,
					Value: m.Variables[v.Value].Value,
				}
			}
			// Setzen der Variable in Kind-Modulen
			for _, nm := range m.NestedModules {
				nm.Variables[v.Name] = m.Variables[v.Name]
			}
		}

		// Speichern der globalen Variablen
		if v.IsGlobal && m.Variables["%passed%"].Value == "true" {
			*globalVars = append(*globalVars, m.Variables[v.Name])
		}
	}

	return nil
}

// Interpretiert die übergebene Bedingung mithilfe von GOJA unter Verwendung aller übergebenen Variablen.
func interpretCondition(cond string, varMap VariableMap) (bool, error) {
	vm := goja.New()
	vars := extractVars(cond)

	// Entfernen der '%'-Zeichen von Variablen in der Condition, damit sie JS-Syntax entspricht
	cond = strings.ReplaceAll(cond, "%", "")
	if err := setVarValues(vars, varMap); err != nil {
		return false, errors.New("GOJA-Error: " + err.Error())
	}

	// Exportiert alle gefundenen Variablen in die JS-VM
	for _, v := range vars {
		if err := vm.Set(v.name, util.CastToAppropriateType(v.value)); err != nil {
			return false, err
		}
	}

	v, err := vm.RunString(cond)
	if err != nil {
		return false, errors.New("GOJA-Error: " + err.Error())
	}

	return v.ToBoolean(), nil
}

// extractVars extrahiert alle Variablen aus der Condition
// und gibt sie als Array zurück
func extractVars(cond string) (v []*moduleVariable) {
	vars := acutil.GetVariablesInString(cond)

	for i := range vars {
		vars[i] = strings.Trim(vars[i], "%")
		v = append(v, &moduleVariable{
			name:     vars[i],
			origName: strings.ToLower("%" + vars[i] + "%"),
		})
	}
	return
}

// setVarValues weist den extrahierten Variablen den passenden Wert aus dem Modul zu
func setVarValues(vars []*moduleVariable, varMap VariableMap) error {
	for _, v := range vars {
		if mvar, ok := varMap[v.origName]; ok {
			v.value = mvar.Value
		} else {
			return errors.New(fmt.Sprintf("variable '%v' does not exist", v.origName))
		}
	}
	return nil
}

// Fügt Schritte, die, warum auch immer nicht ausgeführt, z.B. weil das Elternmodul fehlgeschlafen ist, werden der Progressbar als "fertig" zu.
func addNotExecutedToProgressBar(m AuditModule) {
	for _, nested := range m.NestedModules {
		_ = static.ProgressBar.Add(1)
		addNotExecutedToProgressBar(nested)
	}
	return
}

// Ersetzt alle Variablen im übergebenenen String mit dem Wert der passenden Variable, wenn in der übergebenenen Map vorhanden. Löst außerdem Pfade auf.
func replaceVariablesInString(s string, variables VariableMap) string {
	vars := acutil.GetVariablesInString(s)

	for i, v := range vars {
		vars[i] = strings.ToLower(v)
		s = strings.ReplaceAll(s, v, vars[i])
	}

	for _, v := range vars {
		if path, err := util.GetAbsolutePath(variables[v].Value); err == nil {
			if _, err := os.Stat(path); err == nil {
				s = strings.ReplaceAll(s, v, path)
			} else {
				s = strings.ReplaceAll(s, v, variables[v].Value)
			}
		}
	}
	return s
}
