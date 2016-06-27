package lint_test

import (
	"fmt"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/lagarciag/codenanny/installer"
	"github.com/lagarciag/codenanny/lint"
	"github.com/lagarciag/codenanny/parser"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	v := t.Run()
	os.Exit(v)

}

func TestLintBasic(t *testing.T) {
	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error(err)
	}
	var1 := "parser/parser.go"
	var2 := "parser/parser_test.go"
	var3 := "cmd/root.go"
	var4 := "lint/lint.go"
	argsVars := []string{var1, var2, var3, var4}
	dirList, pkag, err := parser.Parse(argsVars)

	if err != nil {
		t.Error("Error:", err)
	}

	if len(pkag) != 3 {
		t.Error("List of packages must be 3")
	}
	if len(dirList) != 3 {
		t.Error("List of dir must be 3")
	}

	t.Log("PKGs:", pkag)
	t.Log("DIRs:", dirList)

	CreateUnCheckedError()

	err = lint.CheckPackages(pkag)

	if err == nil {
		t.Error("At least 1 error should have been detected")
	}

	err = lint.CheckDirs(dirList)

	if err != nil {
		t.Error("Lint dirs failed")
	}

}

func CreateUnCheckedError() (err error) {
	return nil
}

func DeadCode() {
	fmt.Print("Hello World")
}
