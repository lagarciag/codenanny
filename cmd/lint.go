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

	"regexp"

	"fmt"

	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/config"
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
		if err := config.LoadConfig(); err != nil {
			log.Fatal("error loading config:", err)
		}
		dirSlice := parseListFromArgs()

		log.Debug("DIR SLICE:", dirSlice)

		if err := doLint(dirSlice); err != nil {
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
			match, _ := regexp.MatchString(".go$", fArg)
			if match {
				cList.PushBack(fArg)
			}

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
	var singPkgErr error
	var multiPkgErr error
	var dirCheckErr error
	var dirCheckRecErr error

	listSlice, _ = filterList(listSlice)

	if err = installer.CheckExternalDependencies(); err != nil {
		return err
	}

	// check if there are disabled tools in the configuration
	for key := range config.GlobalConfig.Disabled {
		log.Debug("CONFIG DISABLED:", key)
		if config.GlobalConfig.Disabled[key] {
			installer.DisabledTool[key] = true
		}
	}

	dirList, pkag, err := parser.Parse(listSlice)

	log.Debug("Packages:", pkag)
	log.Debug("dirList:", dirList)

	if err != nil {
		log.Error(err)
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(4)

	multiPackages := func() {
		err = lint.CheckMultiPackages(pkag)
		if err != nil {
			multiPkgErr = fmt.Errorf("Multi packages checker failed:%s", err.Error())
			log.Error(multiPkgErr.Error())
		}
		wg.Done()
	}

	singlePackages := func() {
		err = lint.CheckSinglePackages(pkag)
		if err != nil {
			singPkgErr = fmt.Errorf("Single packages checker failed:%s", err.Error())
			log.Error(singPkgErr.Error())
		}
		wg.Done()
	}

	checkRecDirs := func() {
		err = lint.CheckRecursiveDirs(dirList)
		if err != nil {
			dirCheckRecErr = fmt.Errorf("Directory recursive checker failed:,%s", err.Error())
			log.Error(dirCheckRecErr.Error())
		}
		wg.Done()
	}

	checkDirs := func() {
		err = lint.CheckMultiDirs(dirList)
		if err != nil {
			dirCheckErr = fmt.Errorf("Single dir checker failed:,%s", err.Error())
			log.Error(dirCheckErr.Error())
		}
		wg.Done()
	}

	go multiPackages()
	go singlePackages()
	go checkRecDirs()
	go checkDirs()

	wg.Wait()

	if multiPkgErr != nil || singPkgErr != nil || dirCheckErr != nil || dirCheckRecErr != nil {
		err = fmt.Errorf("Linters failed:%s", "")
	}

	return err
}

func filterList(listSlice []string) (newListSlice []string, err error) {
	tmpList := list.New()
	for _, file := range listSlice {
		match := false
		if config.GlobalConfig.IgnorePath != "" {
			match, _ = regexp.MatchString(config.GlobalConfig.IgnorePath, file)
		}
		if !match {
			tmpList.PushBack(file)
		} else {
			log.Debug("Ignore path matched:", file)
		}
	}
	newListSlice = make([]string, tmpList.Len())
	count := 0
	for e := tmpList.Front(); e != nil; e = e.Next() {
		newListSlice[count] = e.Value.(string)
		count++
	}
	return newListSlice, err
}

func init() {
	RootCmd.AddCommand(lintCmd)
	//RootCmd.PersistentFlags().StringVar(&list, "list", "", "list of files to process")
	lintCmd.Flags().StringVar(&argList, "list", "./", "list of files to process")

}
