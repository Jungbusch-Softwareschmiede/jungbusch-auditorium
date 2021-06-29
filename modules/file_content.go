package modules

import (
	"bytes"
	"errors"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

func (mh *MethodHandler) FileContentInit() ModuleSyntax {
	return ModuleSyntax{
		ModuleName:          "FileContent",
		ModuleDescription:   "FileContent gibt den Inhalt der angegebenen Datei. Ist der grep-parameter gesetzt, werden nur die Zeilen zurückgegeben, in denen das Pattern gefunden wurde.",
		ModuleAlias:         []string{"file_content", "fileContent"},
		ModuleCompatibility: []string{"all"},
		InputParams: ParameterSyntaxMap{
			"file": ParameterSyntax{
				ParamName:        "file",
				ParamAlias:       []string{"datei"},
				ParamDescription: "Pfad zur Datei",
			},
			"grep": ParameterSyntax{
				ParamName:        "grep",
				IsOptional:       true,
				ParamDescription: "Optionaler Suchbegriff, entspricht Pipen des Outputs in grep",
			},
		},
	}
}

// FileContent gibt den Inhalt der angegebenen Datei als String zurück.
// Ist der grep-Parameter gesetzt, werden auch Zeilen zurückgegeben, in denen
// das Pattern gefunden wurde.
func (mh *MethodHandler) FileContent(params ParameterMap) (r ModuleResult) {
	// Finden von allen passenden Dateien, falls Wildcards angegeben wurden
	glob, err := filepath.Glob(params["file"])
	if err != nil {
		r.Err = err
		return
	}

	if glob == nil {
		r.Err = errors.New("Datei wurde nicht gefunden: " + params["file"])
		return
	}

	// Iteriert über jede gefundene Datei
	for _, file := range glob {
		file, err = util.GetAbsolutePath(file)
		if err != nil {
			r.Err = err
			return
		}

		r.Artifacts = append(r.Artifacts, Artifact{Name: "file", Value: file, IsFile: true})

		content, err := os.ReadFile(file)
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				r.Err = errors.New("Zur Ausführung des Moduls werden Administrator, bzw. Root-Privilegien benötigt.")
			} else {
				r.Err = err
			}
			return
		}

		// Entfernen des CarriageReturn-Characters aus CRLF-formatierten Dateien
		content = bytes.ReplaceAll(content, []byte("\r"), []byte(""))
		r.ResultRaw += "\n" + string(content)

		// Falls der Grep-Parameter angegeben wurde, wird hier grep ausgeführt
		if params["grep"] != "" {
			r.Result += mh.Grep(ParameterMap{
				"input": r.ResultRaw,
				"grep":  params["grep"],
			}).Result
		}
	}

	return
}

func (mh *MethodHandler) FileContentValidate(params ParameterMap) error {
	if _, err := util.GetAbsolutePath(params["file"]); err != nil {
		return errors.New("Modul: FileContent - " + err.Error())
	}

	_, err := regexp.Compile(params["grep"])
	if err != nil {
		return errors.New("Modul: FileContent - " + err.Error())
	}

	return nil
}
