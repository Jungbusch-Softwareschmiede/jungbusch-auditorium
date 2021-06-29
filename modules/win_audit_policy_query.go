// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/pkg/errors"
	s "strings"
)

func (mh *MethodHandler) AuditPolicyQueryInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "AuditPolicyQuery",
		ModuleAlias:                []string{"auditpolicyquery", "auditpol"},
		ModuleDescription:          "AuditPolicyQuery ermöglicht den Bereich 'Advanced Audit Policy Configuration' der Security Settings auszulesen.",
		ModuleCompatibility:        []string{"windows"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"guid": ParameterSyntax{
				ParamName:        "guid",
				ParamDescription: "GUID der jeweiligen Value. Bitte im Benutzerhandbuch oder Internet nachschlagen.",
			},
		},
	}
}

func (mh *MethodHandler) AuditPolicyQuery(params ParameterMap) (r ModuleResult) {
	guid := params["guid"]
	guid = s.ReplaceAll(guid, " ", "")

	// Fehlende Klammern um GUID setzen
	if !s.HasPrefix(guid, "{") && !s.HasSuffix(guid, "}") {
		guid = "{" + guid + "}"
	}

	if s.Count(guid, "{") != 1 && s.Count(guid, "}") != 1 {
		r.Err = errors.New("Die Klammern der GUID sind ungültig.")
		return
	}

	// Status der Erweiterten Überwachungsrichtlinienkonfiguration abfragen
	status, err := util.ExecCommand("cmd.exe /c \"auditpol /get /subcategory:" + params["guid"] + "\"")

	if err != nil {
		r.Err = err
		return
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:  "auditpol /get /subcategory:" + params["guid"],
		Value: status,
	})

	r.ResultRaw = status

	// Prüfen ob ein Fehler aufgetreten ist
	if s.Contains(status, "/?") || s.Contains(status, "0x00000057") {
		r.Err = errors.New("Die GUID ist ungültig.")
		return
	}

	// !!ACHTUNG!!
	// r.Result wird in der jeweiligen Sprache des OS zurückgegeben.
	// Bsp. Mögliche deutsche Wörter: Erfolg, und, Fehler, Keine Überwachung
	r.Result = s.TrimSpace(s.ReplaceAll(status[s.LastIndex(status, "  ")+2:], "\r", ""))
	return
}
