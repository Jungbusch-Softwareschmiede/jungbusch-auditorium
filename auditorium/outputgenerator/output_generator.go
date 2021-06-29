package outputgenerator

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/static"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ReportEntry struct {
	ID        string            `json:"id"`
	Desc      string            `json:"desc"`
	Print     string            `json:"print,omitempty"`
	Artifacts []models.Artifact `json:"artifacts"`
	Expected  string            `json:"expected,omitempty"`
	Actual    string            `json:"actual,omitempty"`
	Error     string            `json:"error,omitempty"`
	Result    string            `json:"result"`
	Nested    []ReportEntry     `json:"nested,omitempty"`
}

type Metadata struct {
	Os           string        `json:"os"`
	Root         bool          `json:"root"`
	Total        int           `json:"total"`
	Passed       int           `json:"passed"`
	Notpassed    int           `json:"not_passed"`
	Unsuccessful int           `json:"unsuccessful"`
	Notexecuted  int           `json:"not_executed"`
	Started      string        `json:"started"`
	Elapsed      string        `json:"elapsed"`
	Report       []ReportEntry `json:"report"`
}

type stats struct {
	total        int
	passed       int
	notpassed    int
	unsuccessful int
	notexecuted  int
}

// GenerateOutput erstellt den Report und speichert die Artefakte am spezifizierten Pfad.
func GenerateOutput(report []ReportEntry, outputPath string, numberOfModules int, zip bool, start time.Time, elapsed time.Duration) error {
	Info(SeperateTitle("Output-Generator"))

	if err := saveArtifacts(report, outputPath+static.PATH_SEPERATOR+"artifacts"); err != nil {
		return err
	}
	if err := generateReport(report, outputPath, numberOfModules, start, elapsed); err != nil {
		return err
	}

	if zip {
		if err := zipFolder(outputPath); err != nil {
			return err
		}
	}

	return nil
}

func generateReport(report []ReportEntry, outputPath string, numberOfModules int, start time.Time, elapsed time.Duration) (err error) {
	s := getStats(report, numberOfModules)

	meta := Metadata{
		Os:           static.OperatingSystem,
		Root:         static.HasElevatedPrivileges,
		Total:        s.total,
		Passed:       s.passed,
		Notpassed:    s.notpassed,
		Unsuccessful: s.unsuccessful,
		Notexecuted:  s.notexecuted,
		Started:      start.Format(time.RFC1123),
		Elapsed:      elapsed.String(),
		Report:       cleanArtifacts(report),
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(meta)

	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(outputPath+"/report.json", buf.Bytes(), static.CREATE_FILE_PERMISSIONS); err != nil {
		return err
	}

	return nil
}

func getStats(report []ReportEntry, numberOfModules int) (s stats) {
	s.total = numberOfModules

	for _, entry := range report {
		switch entry.Result {
		case "PASSED":
			s.passed++
		case "NOTPASSED":
			s.notpassed++
		case "UNSUCCESSFUL":
			s.unsuccessful++
		}
		snested := getStats(entry.Nested, numberOfModules)
		s.passed += snested.passed
		s.notpassed += snested.notpassed
		s.unsuccessful += snested.unsuccessful
	}

	s.notexecuted = s.total - (s.passed + s.notpassed + s.unsuccessful)

	return
}

func saveArtifacts(report []ReportEntry, outputPath string) error {
	for _, entry := range report {
		commands := make(map[string]string)
		hasCommand := false
		outputPath += static.PATH_SEPERATOR + entry.ID

		for _, artifact := range entry.Artifacts {
			if err := os.MkdirAll(outputPath, static.CREATE_DIRECTORY_PERMISSIONS); err != nil {
				return err
			}

			if artifact.IsFile {
				filename := path.Base(filepath.ToSlash(artifact.Value))
				artifactPath := outputPath + static.PATH_SEPERATOR + filename

				source, err := os.Open(artifact.Value)
				if err != nil {
					return err
				}

				dest, err := os.Create(artifactPath)
				if err != nil {
					return err
				}

				_, err = io.Copy(dest, source)
				if err != nil {
					return err
				}

				Debug(fmt.Sprintf("Artefakt gespeichert - Pfad: %v", artifactPath))

				if err = source.Close(); err != nil {
					return err
				}
				if err = dest.Close(); err != nil {
					return err
				}

			} else {
				commands[artifact.Name] = artifact.Value
				hasCommand = true
			}
		}
		if hasCommand {
			c, err := json.MarshalIndent(commands, "", "	")
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(outputPath+static.PATH_SEPERATOR+"commands.json", c, static.CREATE_FILE_PERMISSIONS); err != nil {
				return err
			}
		}

		if err := saveArtifacts(entry.Nested, outputPath); err != nil {
			return err
		}

		outputPath = filepath.Dir(outputPath)
	}
	return nil
}

func cleanArtifacts(report []ReportEntry) []ReportEntry {
	for i, entry := range report {
		for j, artifact := range entry.Artifacts {
			if !artifact.IsFile {
				report[i].Artifacts[j].Value = artifact.Name
				report[i].Artifacts[j].Name = "command"
			}
		}
		entry.Nested = cleanArtifacts(entry.Nested)
	}
	return report
}

func zipFolder(folder string) error {
	zipfolder, err := os.Create(folder + ".zip")
	if err != nil {
		return err
	}
	defer zipfolder.Close()

	archive := zip.NewWriter(zipfolder)
	defer archive.Close()

	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filepath.Join(filepath.Base(folder), strings.TrimPrefix(path, folder))

		if d.IsDir() {
			header.Name += static.PATH_SEPERATOR
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
	return err
}
