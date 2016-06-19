package parser_test

import (
	"os"
	"testing"

	"fmt"
	"github.com/lagarciag/gocomlinter/parser"
)

func TestMain(t *testing.M) {
	v := t.Run()
	os.Exit(v)

}

func TestParserBasic(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	var1 := gopath + "src/github.com/lagarciag/gocomlinter/parser/parser.go"
	var2 := gopath + "src/github.com/lagarciag/gocomlinter/parser/parser_test.go"
	var3 := gopath + "src/github.com/lagarciag/gocomlinter/cmd/root.go"

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
