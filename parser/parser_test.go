package parser_test

import (
	"os"
	"testing"

	"fmt"

	"github.com/lagarciag/codenanny/installer"
	"github.com/lagarciag/codenanny/parser"
)

func TestMain(t *testing.M) {
	v := t.Run()
	os.Exit(v)

}

func TestParserBasic(t *testing.T) {
	//gopath := os.Getenv("GOPATH")
	installer.CheckExternalDependencies()
	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error(err)
	}
	var1 := "parser/parser.go"
	var2 := "parser/parser_test.go"
	var3 := "cmd/root.go"

	dirList, pkag, err := parser.Parse(fmt.Sprintf("%s %s %s", var1, var2, var3))

	if err != nil {
		t.Error("Error:", err)
	}

	if len(pkag) != 2 {
		t.Error("List of packages must be 2")
	}
	if len(dirList) != 2 {
		t.Error("List of dir must be 2")
	}

	t.Log("PKGs:", pkag)
	t.Log("DIRs:", dirList)

	t.Log("Pass")

}
