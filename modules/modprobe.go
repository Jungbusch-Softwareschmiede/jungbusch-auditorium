// +build linux

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
)

func (mh *MethodHandler) ModprobeInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "Modprobe",
		ModuleDescription:          "Modprobe simuliert das Laden des angegebenen Moduls zur Laufzeit des Systems und speichert das Ergebnis ausführlich ab. Daraufhin wird überprüft, ob das angegebene Modul aktuell geladen ist.",
		ModuleCompatibility:        []string{"linux"},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"name": ParameterSyntax{
				ParamName:        "name",
				ParamAlias:       []string{},
				ParamDescription: "Name des Moduls",
			},
		},
	}
}

// Modprobe simuliert das Laden des angegebenen Moduls zur Laufzeit des Systems und speichert das Ergebnis ausführlich ab.
// Daraufhin wird überprüft, ob das angegebene Modul aktuell geladen ist.
func (mh *MethodHandler) Modprobe(params ParameterMap) (r ModuleResult) {
	res_modprobe, err := util.ExecCommand("modprobe -n -v " + params["name"])
	if err != nil {
		r.Err = err
		return
	}

	r.ResultRaw = "modprobe: " + res_modprobe
	r.Artifacts = append(r.Artifacts, Artifact{Name: "modprobe -n -v " + params["name"], Value: res_modprobe})

	if res_modprobe != "install /bin/true\n" {
		Debug("Ergebnis von 'modprobe' != 'install /bin/true'")
		r.Result = "false"
		return
	}

	res_lsmod, err := util.ExecCommand("lsmod " + params["name"])
	if err != nil {
		r.Err = err
		return
	}
	r.ResultRaw += "\nlsmod: " + res_lsmod
	r.Artifacts = append(r.Artifacts, Artifact{Name: "lsmod " + params["name"], Value: res_lsmod})

	r.Result = mh.Grep(ParameterMap{
		"input": res_lsmod,
		"grep":  params["name"],
	}).Result

	if r.Result != "" {
		Debug("Ergebnis von 'lsmod | grep " + params["name"] + "' != <no output>")
		r.Result = "false"
		return
	}

	r.Result = "true"
	return
}

func (mh *MethodHandler) ModprobeValidate(params ParameterMap) error {

	return nil
}
