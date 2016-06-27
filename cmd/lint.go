// Copyright © 2016 Luis Garcia <luis.a.garcia@hpe.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"container/list"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/installer"
	"github.com/lagarciag/codenanny/lint"
	"github.com/lagarciag/codenanny/parser"
	"github.com/spf13/cobra"
)

var argList string

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run the linters",
	Long:  `Runs the linters`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if argList == "" {
			log.Fatal("--list flag must be set and point to a list of files that need to be linted")
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("verbose mode enabled")
		}
		if err := doLint(parseListFromArgs()); err != nil {
			log.Fatal("Lint found errors")
		}
	},
}

func parseListFromArgs() (listSlice []string) {

	log.Debug("ARGS:", os.Args)
	cList := list.New()
	takeArg := false
	for _, fArg := range os.Args {
		if takeArg {
			log.Debug("ARG:", fArg)
			cList.PushBack(fArg)
		}

		if fArg == "--list" {
			takeArg = true
		}
	}
	listSize := cList.Len()
	listSlice = make([]string, listSize)

	eID := 0
	for e := cList.Front(); e != nil; e = e.Next() {
		log.Debug("ADDING ARG:", e.Value.(string))
		listSlice[eID] = e.Value.(string)
		eID++
	}

	log.Debug("List:", len(listSlice))
	log.Debug("Processing files:", listSlice)

	return listSlice
}

func doLint(listSlice []string) (err error) {

	if err := installer.CheckExternalDependencies(); err != nil {
		return err
	}

	dirList, pkag, err := parser.Parse(listSlice)

	log.Debug("Packages:", pkag)
	log.Debug("dirList:", dirList)

	if err != nil {
		log.Error(err)
		return err
	}

	err = lint.CheckPackages(pkag)

	if err != nil {
		log.Error("Lint packages failed:", err)
		return err
	}

	err = lint.CheckDirs(dirList)

	if err != nil {
		log.Error("Lint dirs failed")
		return err
	}

	return err
}

func init() {
	RootCmd.AddCommand(lintCmd)
	//RootCmd.PersistentFlags().StringVar(&list, "list", "", "list of files to process")
	lintCmd.Flags().StringVar(&argList, "list", "./", "list of files to process")

}
