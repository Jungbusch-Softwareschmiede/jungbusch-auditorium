package config_test

import (
	"fmt"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-interpreter"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/auditorium/config/config-parser"
	. "github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/models"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util"
	"github.com/Jungbusch-Softwareschmiede/jungbusch-auditorium/util/logger"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	s "strings"
	"testing"
)

var (
	gopath = os.Getenv("gopath") + `\src\github.com\Jungbusch-Softwareschmiede\jungbusch-auditorium\`
)

func InterpretConfig(t *testing.T, cs *ConfigStruct) {
	logger.InitializeLogger(cs, []LogMsg{})
	con, err := config_interpreter.InterpretConfig(cs)
	if err != nil {
		t.Errorf("Fehlgeschlagen: Error aufgetreten: " + err.Error())
	}
	if !con {
		t.Errorf("Continue ist false, obwohl keine One-And-Done Parameter gesetzt wurde.")
	}
	config_parser.ResetFlags()
}

func EvaluateLog(t *testing.T, log []LogMsg, cs ConfigStruct) {
	err := GetLogErr(log)
	if err != "" {
		t.Errorf("Fehlgeschlagen: %v, %v", err, cs)
	}
	config_parser.ResetFlags()
}

func printLog(log []LogMsg) {
	for _, msg := range log {
		fmt.Println(msg.Message)
	}
}

func logCleaner(log []LogMsg, iscontinue bool) (result string) {
	for _, msg := range log {
		if (!iscontinue && msg.AlwaysPrint) || iscontinue {
			result = msg.Message + "\n"
		}
	}
	return
}

// Gibt den ersten Error aus dem Log zurück, wenn vorhanden
func GetLogErr(msg []LogMsg) string {
	for n := range msg {
		if msg[n].Level == 1 {
			return msg[n].Message
		}
	}
	return ""
}

func deletLocalConfigs(dirPath string) {
	// Im ganzen Verzeichniss nach '_local' Datein suchen und diese löschen
	dirPath, err := util.GetAbsolutePath(dirPath)
	if err != nil {
		fmt.Println(err)
	}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if s.Contains(file.Name(), "_local") {
			if err := os.Remove(dirPath + `\` + file.Name()); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getStaticConf(filePath string) string {
	// Aktuelle config einlesen
	fileSlice, err := util.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	// Neues leeres FileSlice
	var newFileSlice []string
	for _, line := range fileSlice {
		// Zeilen mit '%'
		if s.Contains(line, "%") {
			// Pfad zwischen '%' extrahieren und in Absoluten umwandeln
			path := util.GetStringInBetween(line, "%")
			absPath, err := util.GetAbsolutePath(path)
			if err != nil {
				fmt.Println(err)
			}
			// Pfad der config mit absolutem ersetzen
			newLine := s.ReplaceAll(line, "%"+path+"%", absPath)
			newFileSlice = append(newFileSlice, newLine)
		} else {
			newFileSlice = append(newFileSlice, line)
		}
	}

	// Neue Config mit _local im Dateinamen anlegen
	newFileName := filePath[:s.LastIndex(filePath, ".")] + "_local" + filePath[s.LastIndex(filePath, "."):]
	err = util.CreateFile(newFileSlice, newFileName)
	if err != nil {
		fmt.Println(err)
	}
	return newFileName
}

func TestShowFlag(t *testing.T) {
	os.Args = []string{gopath, `-showModules`}
	config, log := config_parser.LoadConfig()
	fmt.Println(log, config)
}

func TestCLIDynamic(t *testing.T) {
	os.Args = []string{gopath, `-config=./test/testdata/cli_testdata/basic_config.ini`}
	config, log := config_parser.LoadConfig()
	EvaluateLog(t, log, config)
}

func TestCLIParameter(t *testing.T) {
	match := ConfigStruct{
		AuditConfig:                  filepath.Clean(gopath + "/test/testdata/cli_testdata/auditDummy2.jba"),
		OutputPath:                   filepath.Clean(gopath + "/test/testdata/cli_testdata/"),
		Config:                       filepath.Clean(gopath + "/test/testdata/cli_testdata/basic_config.ini"),
		VerbosityLog:                 0,
		VerbosityConsole:             0,
		SkipModuleCompatibilityCheck: false,
	}
	os.Args = []string{gopath, "-config=./test/testdata/cli_testdata/basic_config.ini", "-auditConfig=./test/testdata/cli_testdata/auditDummy2.jba", "-outputPath=./test/testdata/cli_testdata/", "1", "1"}

	cs, log := config_parser.LoadConfig()
	EvaluateLog(t, log, cs)
	_, err := config_interpreter.InterpretConfig(&cs)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = logger.InitializeOutput(&cs)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = logger.InitializeLogger(&cs, log)
	if err != nil {
		t.Errorf(err.Error())
	}

	cs.OutputPath = filepath.Dir(cs.OutputPath)

	if cs != match {
		t.Errorf("Ungleiche Config!\nSoll: %v\nIst:  %v", match, cs)
	}
}
