// +build windows

package modules

import (
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
	s "strings"
)

func (mh *MethodHandler) ExportInstalledSoftwareInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "ExportInstalledSoftware",
		ModuleDescription:   "ExportInstalledSoftware speichert alle Namen und Versionsnummern von installierter Software in einer CSV Datei ab.",
		ModuleAlias:         []string{"exportinstalledsoftware"},
		ModuleCompatibility: []string{"windows"},
		InputParams: ParameterSyntaxMap{
			"path": ParameterSyntax{
				ParamName:        "path",
				IsOptional:       true,
				ParamDescription: "Hier kann ein Pfad angegeben werden, an welchem die CSV Datei abgelegt werden soll. Wird kein Pfad angegeben, wird sie in einem Temporären Ordner abgelegt. Unabhängig davon, wird sie so oder so in die Modul-Artefakte und somit den Output aufgenommen.",
			},
		},
	}
}

func (mh *MethodHandler) ExportInstalledSoftware(params ParameterMap) (r ModuleResult) {
	keySlice := []string{
		// x64
		`HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		// x32
		`HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`,
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:   "x64-registry-path",
		Value:  keySlice[0],
	})

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:   "x32-registry-path",
		Value:  keySlice[1],
	})

	separator := ";"
	var nameVersionSlice []string

	for n := range keySlice {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, s.Replace(keySlice[n], `HKEY_LOCAL_MACHINE\`, "", 1), registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			r.Err = err
			return
		}

		// Alle Name der SubKeys auslesen
		subKeys, err := key.ReadSubKeyNames(-1)
		if err != nil {
			r.Err = err
			return
		}

		// Durch alle SubKeys iterieren
		for _, subKeyName := range subKeys {
			// Wert der DisplayName Value auslesen
			displayName, err := util.RegQuery(keySlice[n]+"\\"+subKeyName, "DisplayName")
			if err != nil {
				if err != static.ERROR_VALUE_NOT_FOUND {
					r.Err = err
					return
				}
			}

			// Wert der DisplayVersion Value auslesen
			displayVersion, err := util.RegQuery(keySlice[n]+"\\"+subKeyName, "DisplayVersion")
			if err != nil {
				if err != static.ERROR_VALUE_NOT_FOUND {
					r.Err = err
					return
				}
			}

			/*
				// Leere Einträge nicht berücksichtigen
				if displayName != "" && displayVersion != "" {
					nameVersionSlice = append(nameVersionSlice, displayName + separator + displayVersion)
				}
			*/
			nameVersionSlice = append(nameVersionSlice, displayName+separator+displayVersion)
		}
		err = key.Close()
		if err != nil {
			logger.Err(errors.Wrap(err, "Fehler in ExportInstalledSoftware").Error())
		}
	}

	// nameVersionSlice als CSV exportieren
	var path string
	var err error

	if params["path"] != "" {
		path, err = util.GetAbsolutePath(params["path"] + "InstalledSoftware.csv")
	} else {
		path = static.TempPath + static.PATH_SEPERATOR + "InstalledSoftware.csv"
	}

	if err != nil {
		r.Err = err
	}

	r.Artifacts = append(r.Artifacts, Artifact{
		Name:   "file",
		Value:  path,
		IsFile: true,
	})

	err = util.CreateFile(nameVersionSlice, path)
	if err != nil {
		r.Err = err
		return
	}

	r.Result = "Die CSV Datei wurde erfolgreich gespeichert!"
	return
}
