// Dieses Package ist ein Utility-Package für das Validieren und Parsen der Audit-Konfiguration.
// Alle Methoden des Packages greifen ausschließlich auf die über die Methode SetLines gesetzte, rohe Audit-Konfigurationsdatei zu und sind anderweitig nicht zu verwenden.
package acutil

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/pkg/errors"
	"regexp"
	s "strings"
)

var (
	lines []string // Die Audit-Konfigurationsdatei mit welcher gearbeitet wird
)

// Setzt die Zeilen der Audit-Konfiguration, mit der die Methoden dieses Packages arbeiten.
func SetLines(l []string) {
	lines = l
}

// Konvertiert die Namen aller Variablen im übergebenen string zu lowercase und gibt das Ergebnis zurück.
func VariablesInStringToLower(in string) string {
	variablesInParameter := GetVariablesInString(in)
	for _, v := range variablesInParameter {
		in = s.ReplaceAll(in, v, s.ToLower(v))
	}
	return in
}

// Sucht und returned alle Vorkommnisse von Variablen im übergebenen string
func GetVariablesInString(s string) []string {
	reg, _ := regexp.Compile(static.REG_VARIABLE)
	return reg.FindAllString(s, -1)
}

// Diese Methode entfernt Inlinekommentare und Leerzeichen oder Tabs am Start des übergebenen strings und
// gibt das Ergebnis zurück.
func PrepareLine(line string) (string, error) {
	line, err := RemoveInlineComment(Trim(line))
	if err != nil {
		return "", err
	}

	return Trim(line), nil
}

// Liest die übergebene Auditkonfigurationsdatei ein.
func ReadAuditConfiguration(file string) (lines []string, err error) {
	if file != "" {
		Info("Die Datei " + file + " wird eingelesen.")
		lines, err = util.ReadFile(file)

		if err != nil {
			err = errors.New("Die Audit-Datei " + file + " konnte nicht eingelesen werden: " + err.Error())
		}
	} else {
		err = errors.New("Es wurde keine Audit-Konfigurationsdatei angegeben.")
	}
	return
}

// Sucht anhand eines Keywords das passende ParameterSyntax-Objekt (allgemeingültig) und returned es oder ein leeres Objekt,
// falls nichts gefunden wurde.
func GetParameterSyntaxFromKeyword(keyword string) ParameterSyntax {
	syntax := parameterSyntax()

	for n := range syntax {
		if s.ToLower(syntax[n].ParamName) == s.ToLower(keyword) || util.ArrayContainsString(syntax[n].ParamAlias, keyword) {
			return syntax[n]
		}
	}
	return ParameterSyntax{}
}

// Konvertiert einen Modulparameter-Alias in dessen Name anhand vom übergebenen Modul.
func GetModuleParameterSyntaxNameFromAlias(moduleSyntax ModuleSyntax, alias string) string {
	for _, ms := range moduleSyntax.InputParams {
		if s.ToLower(ms.ParamName) == s.ToLower(alias) || util.ArrayContainsString(ms.ParamAlias, alias) {
			return ms.ParamName
		}
	}
	return ""
}

// Sucht in den übergebenen Modulen anhand des übergebenen Aliases nach dem Name des Moduls.
func GetModuleSyntaxFromNameOrAlias(name string, modules []ModuleSyntax) (ms ModuleSyntax) {
	for n := range modules {
		if s.ToLower(modules[n].ModuleName) == s.ToLower(name) || util.ArrayContainsString(modules[n].ModuleAlias, name) {
			return modules[n]
		}
	}
	return ms
}

// True, wenn der übergebene String im Format einer Condition ist.
func IsCondition(in string) bool {
	ifSyn := ifSyntax()

	condition := "^( *)(" + ifSyn[0].ParamName + ")( *)(\\()( *)([\"`])(.*)([\"`])( *)(\\))"
	match, err := regexp.MatchString(condition, in)

	if err == nil && match {
		return true
	}

	for _, n := range ifSyn[0].ParamAlias {
		// Genereller Syntax
		condition = "^( *)(" + n + ")( *)(\\()( *)([\"`])(.*)([\"`])( *)(\\))"
		match, err = regexp.MatchString(condition, in)
		if err != nil {
			return false
		}

		// Der Syntax, der verwendet wird um die Condition zu finden
		ifreg, err := regexp.Compile("(\\()( *)([\"`])(.*)([\"`])( *)(\\))")
		if err != nil {
			return false
		}

		result := ifreg.FindString(in)
		if match && len(result) > 4 {
			return true
		}
	}
	return false
}

// Splittet eine Variable in Name und Wert.
func SplitVariable(in string) (string, string) {
	return s.TrimSpace(Trim(in[:s.Index(in, "=")])), s.TrimSpace(Trim(in[s.Index(in, "=")+1:]))
}

// Splittet einen Parameter in Name und Wert. Bei einem Multilineparameter wird nur die erste Zeile returned.
func SplitParameter(in string) (string, string) {
	return s.TrimSpace(Trim(in[:s.Index(in, ":")])), s.TrimSpace(Trim(in[s.Index(in, ":")+1:]))
}

// True, wenn der übergebene String dem Format eines Werts entspricht. (Umgeben von Anführungszeichen)
func IsValue(in string) bool {
	in = s.TrimSpace(Trim(in))

	// Jede Value muss entweder von Anführungszeichen oder Backticks umgeben sein
	return (s.HasPrefix(in, "\"") && s.HasSuffix(in, "\"")) ||
		(s.HasPrefix(in, "`") && s.HasSuffix(in, "`"))
}

// True, wenn der übergebene String dem Format einer Variable entspricht. (Also von Prozentzeichen umgeben)
func IsVariable(in string) bool {
	in = s.TrimSpace(Trim(in))

	// könnte mit static Variable Regex und Match besser gemacht werden

	// Jede Variable muss von Prozentzeichen umgeben sein
	return s.HasPrefix(in, "%") && s.HasSuffix(in, "%")
}

// True, wenn die Zeile an der übergebenen Zeilennummer dem Syntax eines Parameters entspricht.
func IsParameter(n *int) bool {
	l := Trim(lines[*n])
	indexDoppelpunkt := s.Index(l, ":")
	indexIstGleich := s.Index(l, "=")

	// Wenn kein =, aber ein :
	return (indexIstGleich == -1 && indexDoppelpunkt != -1) ||
		// Oder = und :, aber : vor =
		(indexIstGleich != -1 && indexDoppelpunkt != -1 && indexDoppelpunkt < indexIstGleich)
}

// True, wenn die Zeile an der übergebenen Zeilennummer dem Syntax einer Variablen-Zuweisung entspricht.
func IsVariableAssignment(n *int) bool {
	l := Trim(lines[*n])
	indexDoppelpunkt := s.Index(l, ":")
	indexIstGleich := s.Index(l, "=")

	// Wenn kein :, aber ein =
	return (indexIstGleich != -1 && indexDoppelpunkt == -1) ||
		// Oder = und :, aber : vor =
		(indexIstGleich != -1 && indexDoppelpunkt != -1 && indexIstGleich < indexDoppelpunkt)
}

// Überspringt für den Parser irrelevante Zeilen.
// Dies inkludiert leere Zeilen, Kommentare sowie Kommentarblöcke.
func SkipIrrelevantLines(lines []string, n *int) error {
	var l string

	for ; *n < len(lines); *n++ {
		// Leere Zeilen überspringen
		SkipEmpty(n)

		if LinesFinished(n) {
			return nil
		}

		l = Trim(lines[*n])

		// Kommentarblöcke überspringen
		if IsCommentBlock(n) {
			GoToEndOfCommentBlock(n)
			if *n >= len(lines) {
				return GenerateSyntaxError(static.COMMENTBLOCK_NEVER_CLOSED, *n+1, l, "")
			}
			continue
		} else

		// Kommentare überspringen
		if IsComment(n) {
			continue
		} else
		// Kein Kommentar oder so mehr gefunden, die Schleife verlassen
		{
			break
		}
	}

	return nil
}

// True, wenn die Konfigurationsdatei leer ist.
func AreLinesEmpty() bool {
	return len(lines) == 0 || (len(lines) == 1 && Trim(lines[0]) == "")
}

// True, wenn das Ende der Konfigurationsdatei erreicht wurde.
func LinesFinished(n *int) bool {
	return *n >= len(lines)
}

// Entfernt Inline-Kommentare aus dem übergebenen String.
func RemoveInlineComment(in string) (string, error) {
	result := in

	// Wenn möglicherweise ein Kommentar vorkommt
	if s.Contains(in, "//") {

		// Das erste Anführungszeichen finden
		indexQuota := s.Index(in, "\"")

		// Den ersten Backtick finden
		indexBacktick := s.Index(in, "`")

		// Der Wert steht in Anführungszeichen
		if indexQuota != -1 && ((indexQuota < indexBacktick) || indexBacktick == -1) {

			// Die ersten zwei Anführungszeichen überspringen. -> Den Wert überspringen
			in = in[s.Index(in, "\"")+1:]
			in = in[s.Index(in, "\"")+1:]

			// Wenn immernoch ein // vorkommt, ist wahrscheinlich irgendwo ein Kommentar
			if s.Contains(in, "//") {

				// Ich überprüfe nun, ob vor dem Kommentar noch ein Anführungszeichen ist
				indexQuota = s.Index(in, "\"")
				indexComment := s.Index(in, "//")

				// Wenn ja, ist in dieser Zeile irgendwas komisch
				if indexQuota != -1 && indexQuota < indexComment {
					return result, errors.New(static.INVALID_LINE + static.TOO_MANY_QUOTATIONMARKS)
				}

				// Wenn nein, wird der Kommentar entfernt
				result = result[:s.LastIndex(result, "//")]
			}
		} else if indexBacktick != -1 {
			// Wenn wir ein Backtick gefunden haben, suchen wir nach dem zweiten Tick
			in = in[s.Index(in, "`")+1:]

			// Wenn kein zweites gefunden wurde, haben wir einen Multiline-Parameter, da ignorieren wir inline-Kommentare
			// Wenn doch, suchen wir den zweiten Backtick
			if s.Contains(in, "`") {
				in = in[s.Index(in, "`")+1:]

				// An dieser Stelle sollten wir den aktuellen Wert übersprungen haben
				indexBacktick = s.Index(in, "`")
				indexComment := s.Index(in, "//")

				// Wenn ein weiteres Tick vor dem Kommentar steht, ist etwas an der Zeile seltsam
				if indexBacktick != -1 && indexBacktick < indexComment {
					return result, errors.New(static.INVALID_LINE + static.TOO_MANY_TICKS)
				}

				// Wenn das nicht der Fall ist und hinter den Backticks noch ein Kommentar steht, wird dieser entfernt
				if s.Contains(in, "//") {
					result = result[:s.LastIndex(result, "//")]
				}
			}
		} else if s.Index(in, "\"") == -1 && s.Index(in, "`") == -1 {
			result = result[:s.LastIndex(result, "//")]
		}
	}

	return s.Trim(Trim(result), " "), nil
}

// True, wenn die Zeile an der übergebenen Zeilennummer ein Kommentar ist.
func IsComment(n *int) bool {
	return s.HasPrefix(Trim(lines[*n]), "//")
}

// Überspringt leere Zeilen.
func SkipEmpty(n *int) {
	for len(lines) > *n && Trim(lines[*n]) == "" {
		*n++
	}
}

// True, wenn die Zeile an der übergebenen Zeilennummer ein Kommentarblock ist.
func IsCommentBlock(n *int) bool {
	return s.HasPrefix(Trim(lines[*n]), "/*")
}

// Überspringt einen Kommentarblock.
func GoToEndOfCommentBlock(n *int) {
	l := Trim(lines[*n])

	// Kommentarblock
	if s.HasPrefix(l, "/*") {

		// Ende des Blocks suchen
		for ; !s.HasSuffix(l, "*/"); {
			if *n+1 >= len(lines) {
				break
			}
			*n++
			l = Trim(lines[*n])
		}
	}
}

// Entfernt alle Leerzeichen und Tabs die vor dem übergebenen String stehen.
func Trim(in string) string {
	for s.HasPrefix(in, " ") || s.HasPrefix(in, "\t") {
		in = in[1:]
	}
	return in
}

// Generiert einen Error des Typs models.SyntaxError aus den übergebenen Informationen.
func GenerateSyntaxError(errorMsg string, lineNo int, line string, keyword string) *SyntaxError {
	return &SyntaxError{
		ErrorMsg:     errorMsg,
		LineNo:       lineNo + 1,
		Line:         Trim(line),
		Errorkeyword: keyword,
	}
}

// Versucht zu evaluieren, ob sich der Fehler an der übergebenen Zeile auf ein fehlendes Komma zurückführen lässt.
func EvaluateClosingBracketError(n *int) error {
	l := Trim(lines[*n])
	if l == "}" {
		return GenerateSyntaxError(static.MODULE_MISSING_CLOSING_BRACKET_COMMA, *n, l, "")
	}
	return GenerateSyntaxError(static.MODULE_MISSING_CLOSING_BRACKET, *n, l, "")
}
