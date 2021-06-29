// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	"os"
	s "strings"
)

func (mh *MethodHandler) IsGPTemplatePresentInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "IsGPTemplatePresent",
		ModuleDescription:   "IsGPTemplatePresent prüft ob das angegebene Group Policy Administrative Template vorhanden ist.",
		ModuleAlias:         []string{"isgptemplatetresent"},
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"templateName": ParameterSyntax{
				ParamName:        "templateName",
				ParamAlias:       []string{},
				ParamDescription: "Vollständiger Dateinamen des Templates z.B.: AdmPwd.admx/adml",
			},
		},
	}
}

// Prüfen ob angegebene .admx/.adml File vorhanden ist
func (mh *MethodHandler) IsGPTemplatePresent(params ParameterMap) (r ModuleResult) {

	// Format von templateName prüfen + admx und adml Dateipfad generieren
	admx := ""
	adml := ""
	if s.Contains(params["templateName"], "/") {
		filePathSlice := s.Split(params["templateName"], "/")
		admx = filePathSlice[0]
		adml = `\` + s.ReplaceAll(admx, ".admx", ".adml")
	} else {
		r.Err = errors.New("Der angegebene Templatename ist nicht im Format 'Name.admx/adml'.")
		return
	}

	// Ordnernamen der adml Dateien in allen Sprachen
	languages := []string{
		"de-DE",
		"en-US",
		"cs-CZ",
		"da-DK",
		"el-GR",
		"es-ES",
		"fi-FI",
		"fr-FR",
		"hu-HU",
		"it-IT",
		"ja-JP",
		"ko-KR",
		"nb-NO",
		"nl-NL",
		"pl-PL",
		"pt-BR",
		"pt-PT",
		"ru-RU",
		"sv-SE",
		"tr-TR",
		"zh-CN",
		"zh-TW",
	}

	// !!!Achtung!!! Abhäning von verwendeter Sprache.
	// Prüfen ob die admx und adml Dateien vorhanden sind
	if util.IsFile(os.Getenv("windir") + `\policyDefinitions\` + admx) {
		for n := range languages {
			if util.IsFile(os.Getenv("windir") + `\policyDefinitions\` + languages[n] + adml) {
				r.Result = "true"
				return
			}
		}
	}

	r.Result = "false"
	return
}
