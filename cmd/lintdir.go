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
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/dirlister"
	"github.com/spf13/cobra"
)

var pathFlag string

// lintdirCmd represents the lintdir command
var lintdirCmd = &cobra.Command{
	Use:   "lintdir",
	Short: "runs linters and code checkers on the provided dir using the -p flag",
	Long:  `runs linters and code checkers on the provided dir using the -p flag`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("lintdir called")
		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("verbose mode enabled")
		}

		if err := lintdir(pathFlag); err != nil {
			log.Fatal("Lint dir found errors")
		}

	},
}

func lintdir(path string) (err error) {
	fileSlice, _, err := dirlister.ListDir(path)

	if err != nil {
		return err
	}

	log.Debug("The list:", fileSlice)
	err = doLint(fileSlice)
	return err
}

func init() {
	RootCmd.AddCommand(lintdirCmd)
	//RootCmd.PersistentFlags().StringVar(&path, "path", "p", "path to lint")
	lintdirCmd.PersistentFlags().StringVarP(&pathFlag, "path", "p", "./", "path to lint")
}
