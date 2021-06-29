package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
)

func (mh *MethodHandler) TemplateInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:                 "Template",
		ModuleDescription:          "",
		ModuleAlias:                []string{},
		ModuleCompatibility:        []string{},
		RequiresElevatedPrivileges: true,
		InputParams: ParameterSyntaxMap{
			"param": ParameterSyntax{
				ParamName:        "param",
				ParamAlias:       []string{},
				ParamDescription: "",
				IsOptional:       false,
			},
		},
	}
}

func (mh *MethodHandler) TemplateValidate(params ParameterMap) error {

	return nil
}

func (mh *MethodHandler) Template(params ParameterMap) (r ModuleResult) {

	return
}
