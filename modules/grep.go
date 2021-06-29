package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func (mh *MethodHandler) GrepInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Grep",
		ModuleDescription:   "Grep dient als Suchfunktion und kann in verschiedenen Modulen aufgerufen werden.",
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"input": ParameterSyntax{
				ParamName:        "input",
				ParamDescription: "Übergebener String, in dem gesucht werden soll",
			},
			"grep": ParameterSyntax{
				ParamName:        "grep",
				ParamDescription: "Suchbegriff/Regex-Ausdruck",
			},
		},
	}
}

// Grep liefert die Zeile, wo der Suchbegriff im Input übereinstimmt als String zurück.
func (mh *MethodHandler) Grep(params ParameterMap) (r ModuleResult) {
	invert := false
	split := strings.SplitN(params["grep"], " ", 2)
	grep := ""

	if len(split) == 1 {
		grep = split[0]
	} else {
		if split[0][0] == '-' {
			if strings.Contains(split[0], "i") {
				grep = "(?i)" + split[1]
			}
			if strings.Contains(split[0], "v") {
				invert = true
			}
		} else {
			grep = params["grep"]
		}
	}

	re, err := regexp.Compile(grep)
	if err != nil {
		r.Err = err
		return
	}

	for _, line := range strings.SplitAfter(params["input"], "\n") {
		if !invert && re.MatchString(line) || invert && !re.MatchString(line) {
			r.Result += line
		}
	}
	return
}

func (mh *MethodHandler) GrepValidate(params ParameterMap) error {
	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: Auditctl - " + err.Error())
	}
	return nil
}
