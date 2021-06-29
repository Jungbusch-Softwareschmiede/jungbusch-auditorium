// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
)

func (mh *MethodHandler) RegistryQueryInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "RegistryQuery",
		ModuleAlias:       []string{"registryquery", "regquery", "registry"},
		ModuleDescription:   "RegistryQuery gibt die Value eines Registry Keys aus.",
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"key": ParameterSyntax{
				ParamName:        "key",
				ParamDescription: "Vollständiger Key z.B.: HKEY_LOCAL_MACHINE\\SOFTWARE\\Policies...",
			},
			"value": ParameterSyntax{
				ParamName:        "value",
				ParamDescription: "Value des Registry Keys.",
			},
		},
	}
}

func (mh *MethodHandler) RegistryQuery(params ParameterMap) (r ModuleResult) {
	res, err := util.RegQuery(params["key"], params["value"])

	if err != nil {
		switch {
		case err == static.ERROR_VALUE_NOT_FOUND:
			r.ResultRaw = err.Error()
			r.Result = "valueNonExistent"
		case err == static.ERROR_KEY_NOT_FOUND:
			r.ResultRaw = err.Error()
			r.Result = "keyNonExistent"
		default:
			r.Err = err
		}
	} else {

		r.Artifacts = append(r.Artifacts, Artifact{
			Name:  "Key: " + params["key"] + " | Value: " + params["value"],
			Value: res,
		})

		r.ResultRaw = res
		r.Result = r.ResultRaw
	}
	return
}
