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
	formatter := &log.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableTimestamp = true
	log.SetFormatter(formatter)
	v := t.Run()
	os.Exit(v)

}

func TestLintBasic(t *testing.T) {
	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error(err)
	}

	var1 := "packagewitherrors/packagewitherrors.go"
	var2 := "parser/parser_test.go"
	var3 := "cmd/root.go"
	var4 := "lint/lint.go"
	argsVars := []string{var1, var2, var3, var4}

	if err := lint.ChgDirToGitRootPath(); err != nil {
		t.Error("Error changing to root dir")
	}

	dirList, pkag, err := parser.Parse(argsVars)

	if err != nil {
		t.Error("Error:", err)
	}

	if len(pkag) != 4 {
		t.Error("List of packages must be 4", pkag)
	}
	if len(dirList) != 4 {
		t.Error("List of dir must be 4", pkag)
	}

	t.Log("PKGs:", pkag)
	t.Log("DIRs:", dirList)

	CreateUnCheckedError()

	errCount := 0

	err = lint.CheckSinglePackages(pkag)

	if err != nil {
		errCount++
	}

	err = lint.CheckMultiPackages(pkag)

	if err != nil {
		errCount++
	}

	err = lint.CheckMultiDirs(dirList)

	if err != nil {
		errCount++
	}

	if errCount == 0 {
		t.Error("Should have detected errors")
	}

}

func TestVet(t *testing.T) {
	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error(err)
	}
	if err := lint.ChgDirToGitRootPath(); err != nil {
		t.Error("Error changing to root dir:", err)
	}

	var4 := "lint/lint.go"
	argsVars := []string{var4}
	dirList, pkag, err := parser.Parse(argsVars)

	if err != nil {
		t.Error("Error:", err)
	}

	if len(pkag) != 1 {
		t.Error("List of packages must be 1", pkag)
	}
	if len(dirList) != 1 {
		t.Error("List of dir must be 1", pkag)
	}

	t.Log("PKGs:", pkag)
	t.Log("DIRs:", dirList)

	errCount := 0

	err = lint.CheckSinglePackages(pkag)

	if err != nil {
		errCount++
	}

	err = lint.CheckMultiPackages(pkag)

	if err != nil {
		errCount++
	}

	err = lint.CheckRecursiveDirs(dirList)

	if err != nil {
		errCount++
	}

	err = lint.CheckMultiPackages(dirList)

	if err != nil {
		errCount++
	}

	err = lint.CheckMultiDirs(dirList)

	if err != nil {
		errCount++
	}

	if errCount == 0 {
		t.Error("No erros detected")
	}

}

func CreateUnCheckedError() (err error) {
	return nil
}

func DeadCode() {
	fmt.Print("Hello World")
}

func vetError() (err error) {
	if err = CreateUnCheckedError(); err != nil {
		return err
	}
	return err
}

func vetShadowError() (err error) {
	if true {

		if err := CreateUnCheckedError(); err != nil {
			return err
		}
	}
	return err
}
