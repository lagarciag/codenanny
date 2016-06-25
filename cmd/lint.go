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
	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/gocomlinter/lint"
	"github.com/lagarciag/gocomlinter/parser"
	"github.com/spf13/cobra"
)

var list string

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run the linters",
	Long:  `Runs the linters`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if list == "" {
			log.Fatal("--list flag must be set and point to a list of files that need to be linted")
		}
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("verbose mode enabled")
		}
		doLint()
	},
}

func doLint() {

	log.Debug("List:", len(list))
	log.Debug("Processing files:", list)

	dirList, pkag, err := parser.Parse(list)

	log.Debug("Packages:", pkag)

	if err != nil {
		log.Error(err)
	}

	err = lint.LintPackages(pkag)

	if err != nil {
		log.Error("Lint packages failed:", err)
	}

	err = lint.LintDirs(dirList)

	if err != nil {
		log.Error("Lint dirs failed")
	}

}

func init() {
	RootCmd.AddCommand(lintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lintCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().StringVar(&list, "list", "", "list of files to process")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
