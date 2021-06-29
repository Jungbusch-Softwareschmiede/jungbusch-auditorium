package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/dop251/goja"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func (mh *MethodHandler) ScriptInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "Script",
		ModuleDescription:   "Mit Script lässt sich JavaScript-Code ausführen. Es lässt sich auf alle anderen Module zugreifen.",
		ModuleAlias:         []string{"javascript", "js"},
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"script": ParameterSyntax{
				ParamName:        "script",
				ParamAlias:       []string{"js"},
				ParamDescription: "Auszuführendes Script",
			},
		},
	}
}

// Script führt JavaScript-Code aus um zum Beispiel auf andere Module zugreifen zu können
func (mh *MethodHandler) Script(params ParameterMap, variables *VariableMap) (r ModuleResult) {

	vm := goja.New()
	vm.SetFieldNameMapper(
		goja.UncapFieldNameMapper(),
	)
	if err := vm.Set("params", vm.ToValue(ParameterMap{})); err != nil {
		r.Err = err
		return
	}

	if err := vm.Set("newResult", newResult); err != nil {
		r.Err = err
		return
	}

	if err := initModules(vm, mh); err != nil {
		r.Err = err
		return
	}

	if err := initLogging(vm); err != nil {
		r.Err = err
		return
	}

	if err := initVariables(vm, *variables); err != nil {
		r.Err = err
		return
	}

	res, err := vm.RunString(params["script"])
	if err != nil {
		if ierr, ok := err.(*goja.InterruptedError); ok {
			r.Err = ierr.Value().(ModuleResult).Err
			return
		} else {
			r.Err = err
			return
		}
	}

	r, ok := res.Export().(ModuleResult)
	if !ok {
		// Das Skript liefert kein Ergebnis vom Typ ModuleResult
		t := ""
		switch res.ExportType().(type) {
		case nil:
			t = "nil"
		default:
			t = res.ExportType().String()
		}
		r.Err = errors.New("Das letzte Statement im Script muss vom Typ ModuleResult sein, ist " + t)
		return
	}

	variables = exportVariables(vm, *variables)

	return
}

// Initialisiert den MethodHandler in JS mit allen Modulen
func initModules(vm *goja.Runtime, mh *MethodHandler) (err error) {
	artifacts := make([]Artifact, 0)

	// Iteriert über alle Module und definiert sie in der JS-VM
	t := reflect.TypeOf(mh)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !strings.HasSuffix(m.Name, "Init") && !strings.HasSuffix(m.Name, "Validate") {
			err = vm.Set(m.Name, func(params ParameterMap) ModuleResult {
				in := []reflect.Value{reflect.ValueOf(mh), reflect.ValueOf(params)}
				res := m.Func.Call(in)[0].Interface().(ModuleResult)
				if res.Err != nil {
					vm.Interrupt(res)
				}
				// Speichert Artefakte für alle ausgeführten Module

				artifacts = append(artifacts, res.Artifacts...)
				res.Artifacts = artifacts

				return res
			})
			if err != nil {
				return err
			}
		}
	}
	return
}

func (mh *MethodHandler) ScriptValidate(params ParameterMap) error {
	_, err := goja.Compile("", params["script"], false)
	if err != nil {
		return errors.New(params["script"] + "\n" + err.Error())
	} else {
		return nil
	}
}

func initLogging(vm *goja.Runtime) (err error) {
	if err = vm.Set("info", logger.Info); err != nil {
		return err
	}
	if err = vm.Set("warn", logger.Warn); err != nil {
		return err
	}
	if err = vm.Set("err", logger.Err); err != nil {
		return err
	}
	if err = vm.Set("debug", logger.Debug); err != nil {
		return err
	}
	return
}

func initVariables(vm *goja.Runtime, variables VariableMap) (err error) {
	for _, v := range variables {
		if !v.IsEnv {
			if err = vm.Set(strings.Trim(v.Name, "%"), v.Value); err != nil {
				return
			}
		}
	}
	return
}

func exportVariables(vm *goja.Runtime, variables VariableMap) *VariableMap {
	for _, v := range variables {
		if !v.IsEnv {
			variables[v.Name] = Variable{
				Name:  v.Name,
				Value: vm.Get(strings.Trim(v.Name, "%")).Export().(string),
			}
		}
	}
	return &variables
}

func newResult(result, resultRaw string, err error) ModuleResult {
	return ModuleResult{
		Result:    result,
		ResultRaw: resultRaw,
		Err:       err,
	}
}
