package parser_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/lagarciag/codenanny/installer"
	"github.com/lagarciag/codenanny/parser"

	"github.com/lagarciag/codenanny/lint"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableTimestamp = true
	log.SetFormatter(formatter)
	v := t.Run()
	os.Exit(v)

}

func TestParserBasic(t *testing.T) {
	//gopath := os.Getenv("GOPATH")
	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error("Error found checking dependencies")
	}

	if err := lint.ChgDirToGitRootPath(); err != nil {
		t.Error(err)
	}

	var1 := "parser/parser.go"
	var2 := "parser/parser_test.go"
	var3 := "cmd/root.go"
	argsSlice := []string{var1, var2, var3}
	log.Info("string to parse:", argsSlice)
	dirList, pkag, err := parser.Parse(argsSlice)

	if err != nil {
		t.Error("Error:", err)
	}

	if len(pkag) != 2 {
		t.Error("List of packages must be 2 but got ", len(pkag))
	}
	if len(dirList) != 2 {
		t.Error("List of dir must be 2 but got ", len(dirList))
	}

	t.Log("PKGs:", pkag)
	t.Log("DIRs:", dirList)

	t.Log("Pass")

}
