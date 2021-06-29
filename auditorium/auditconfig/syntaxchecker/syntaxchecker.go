// Dieses Package ist für das Überprüfen des Syntaxes der Audit-Konfigurationsdatei zuständig.
package syntaxchecker

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/auditconfig/acutil"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/pkg/errors"
	"regexp"
	s "strings"
)

// Diese Methode überprüft den Syntax der übergebenen Audit-Konfigurationsdatei.
// Sie gibt einen Error zurück, der in den Typ models.SyntaxError gecasted werden kann.
func Syntax(lines []string) error {
	n := new(int)
	var err error
	var l string
	SetLines(lines)

	// Überprüfen, ob die Datei leer ist
	if AreLinesEmpty() {
		return GenerateSyntaxError(static.AUDIT_CONFIG_EMPTY, -1, "", "")
	}

	// Alle Zeilen iterieren
	for ; *n < len(lines); *n++ {

		// Kommentare, leere Zeilen überspringen
		if err = SkipIrrelevantLines(lines, n); err != nil {
			return err
		}

		if LinesFinished(n) {
			break
		}

		l, err = RemoveInlineComment(Trim(lines[*n]))
		if err != nil {
			return GenerateSyntaxError(err.Error(), *n, l, "")
		}

		// An dieser Stelle muss jetzt ein Modul geöffnet werden
		if l == "{" {
			if err = syntaxCheckModule(lines, n); err != nil {
				return err
			}
		} else {
			return GenerateSyntaxError(static.MODULE_MISSING_OPENING_BRACKET, *n, l, "")
		}
	}

	return nil
}

// Diese Methode überprüft den Syntax eines einzelnen Moduls. Sie wird rekursiv für verschachtelte Module aufgerufen.
func syntaxCheckModule(lines []string, n *int) (err error) {

	// Überprüfen, ob die Datei zu Ende ist
	*n++
	if LinesFinished(n) {
		return GenerateSyntaxError(static.MODULE_EMPTY, *n, lines[len(lines)-1], "")
	}
	l := Trim(lines[*n])

	// Zeilen iterieren bis das aktuelle Modul endet
	for err == nil && l != "}," {

		// Überprüfen, ob die Datei zu Ende ist
		if LinesFinished(n) {
			return EvaluateClosingBracketError(n)
		}

		// Kommentare, leere Zeilen überspringen
		err = SkipIrrelevantLines(lines, n)
		l, err = PrepareLine(lines[*n])
		if err != nil {
			return GenerateSyntaxError(err.Error(), *n, l, "")
		}

		// Fehlendes Komma bei der schließenden Klammer abfangen
		if l == "}" {
			return EvaluateClosingBracketError(n)
		}

		// Verschachteltes Modul
		if l == "{" {
			if err = syntaxCheckModule(lines, n); err != nil {
				return err
			}
			*n++
			l, err = PrepareLine(lines[*n])
			continue
		}

		// Die Zeile enthält einen Parameter
		if IsParameter(n) {
			if err = validateParameter(lines, n); err != nil {
				return err
			}
			*n++

			// Multiline-Parameter können die Zeilennummer manipulieren
			if LinesFinished(n) {
				err = GenerateSyntaxError(static.MODULE_MISSING_CLOSING_BRACKET, *n-1, lines[len(lines)-1], "")
				break
			}

			// Die nächste Zeile vorbereiten, dann an den Anfang der Schleife gehen
			l, err = PrepareLine(lines[*n])
			continue
		} else

		// Die Zeile enthält eine Variable
		if IsVariableAssignment(n) {
			if err = validateVariableAssignment(lines, n); err != nil {
				return err
			}
			*n++
			l, err = PrepareLine(lines[*n])
			continue
		}

		// Ungültiger Ausdruck
		return GenerateSyntaxError(static.MODULE_INVALID_EXPRESSION, *n, l, "")
	}

	return err
}

// Diese Methode validiert den Syntax einer Variablen-Zuweisung.
func validateVariableAssignment(lines []string, n *int) (err error) {
	line, err := RemoveInlineComment(lines[*n])
	if err != nil {
		return err
	}

	// Name und Wert der Variablenzuweisung trennen
	name, value := SplitVariable(line)

	// Name der Variable validieren
	if err = validateVariableName(name); err != nil {
		return GenerateSyntaxError(err.Error(), *n, lines[*n], name)
	}

	// ist Wert der Variable ein Wert?
	if IsValue(value) {
		if err = validateVariableValue(value); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], value)
		}
	} else

	// ist Wert der Variable eine weitere Variable?
	if IsVariable(value) {
		if err = validateVariableName(name); err != nil {
			return GenerateSyntaxError(err.Error(), *n, lines[*n], value)
		}
	} else {
		// Weder Wert noch Variable, also Error
		return GenerateSyntaxError(static.VARIABLE_INVALID_VALUE, *n, lines[*n], value)
	}

	return nil
}

// Diese Methode validiert den Syntax eines Parameters.
func validateParameter(lines []string, n *int) (err error) {
	line, err := PrepareLine(lines[*n])
	if err != nil {
		return err
	}

	name, value := SplitParameter(line)

	// Parametername validieren
	if err = validateParameterName(name); err != nil {
		return GenerateSyntaxError(err.Error(), *n, lines[*n], name)
	}

	// Wert validieren
	if lineNo, err := validateParameterValue(value, lines, n); err != nil {
		return GenerateSyntaxError(err.Error(), lineNo, lines[lineNo], value)
	}

	return nil
}

// Diese Methode validiert den Syntax eines Variablennames.
func validateVariableName(name string) error {
	// Variablenname checken
	if !(s.HasPrefix(name, "%") && s.HasSuffix(name, "%")) {
		return errors.New(static.VARIABLE_INVALID_NAME + static.VARIABLE_MISSING_PERCENTAGE)
	}

	// Erlaubte Zeichen des Variablennames checken
	match, _ := regexp.MatchString(static.REG_VARIABLE_NAME, s.Trim(name, "%"))
	if !match {
		return errors.New(static.VARIABLE_INVALID_NAME + static.INVALID_CHARACTERS)
	}

	return nil
}

// Diese Methode validiert den Syntax des Werts einer Variable.
func validateVariableValue(_ string) error {
	// Aktuell gibt es nichts zu validieren, außer dass ein Wert in Anführungszeichen oder Backticks stehen muss
	// Methode ist hier für einheitliche Validierung oder spätere Additionen
	return nil
}

// Diese Methode validiert den Syntax eines Parameter-Namens.
func validateParameterName(in string) error {
	// Erlaubte Zeichen im Parametername checken
	match, _ := regexp.MatchString(static.REG_MODULE_NAME, in)
	if !match {
		return errors.New(static.MODULE_INVALID_MODULE_NAME + static.INVALID_CHARACTERS)
	}

	return nil
}

// Diese Methode validiert den Wert eines Parameters.
func validateParameterValue(value string, lines []string, n *int) (int, error) {
	var err error

	// Multiline-Parameter validieren
	firstLine := *n
	if s.HasPrefix(value, "`") {

		// Inline-Kommentar entfernen
		value, err = PrepareLine(lines[*n])
		if err != nil {
			return 0, err
		}

		// Wenn sich in der aktuellen Zeile nur ein einzelnes Backtick befindet, wissen wir, dass der
		// Multiline-Parameter über mehrere Zeilen geht. Diese Zeile muss einzeln abgefangen und
		// übersprungen werden, da die folgende for-Schleife nicht funktioniert, wenn in der ersten Zeile
		// ausschließlich ein Backtick und sonst nichts steht.
		if s.Index(value, "`") == s.LastIndex(value, "`") {
			*n++
			value, err = PrepareLine(lines[*n])
			if err != nil {
				return 0, err
			}
		}

		// Den Multiline-Parameter iterieren
		for ; !s.HasSuffix(value, "`"); {

			// Hinter dem Ende des Multiline-Parameters darf kein Text stehen.
			if *n != firstLine && s.Contains(value, "`") {
				return firstLine, errors.New(static.MODULE_NO_TEXT_BEHIND_MULTILINE)
			}

			// Das Ende der Datei wurde erreicht, ohne dass der Multiline-Parameter geschlossen wurde.
			if *n+1 >= len(lines) {
				return firstLine, errors.New(static.MODULE_VALUE_INVALID_MULTILINE + static.MISSING_TICKS)
			}
			*n++

			// Die nächste Zeile vorbereiten
			value, err = PrepareLine(lines[*n])
			if err != nil {
				return 0, err
			}
		}

		return *n, nil
	}

	// Standard-Parameter validieren
	// Parameter-Value muss entweder eine Condition sein oder mit Anführungszeichen beginnen/aufhören
	if !IsCondition(value) {
		if !s.HasPrefix(value, "\"") || !s.HasSuffix(value, "\"") {
			return *n, errors.New(static.MODULE_VALUE_INVALID + static.INVALID_SYNTAX_OR_MISSING_QUOTATIONMARKS)
		}

		// Die Parametervalue darf nicht leer sein
		if s.Trim(value, "\"") == "" {
			return *n, errors.New(static.PARAMETER_EMPTY)
		}
	}

	return 0, nil
}
