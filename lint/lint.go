/**
 * Copyright (C) 2015 Hewlett Packard Enterprise Development LP
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

//Package lint implements linting methods
package lint

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bytes"

	"bufio"
	"container/list"
	"io"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/config"
	"github.com/lagarciag/codenanny/installer"
)

var lintersFlag = map[string]string{
	"aligncheck":  `aligncheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`,
	"deadcode":    `deadcode  `,
	"dupl":        `dupl -plumbing -threshold {duplthreshold} ./*.go:^(?P<path>[^\s][^:]+?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)$`,
	"errcheck":    `errcheck `,
	"goconst":     `goconst`,
	"gocyclo":     `gocyclo -over {mincyclo} .:^(?P<cyclo>\d+)\s+\S+\s(?P<function>\S+)\s+(?P<path>[^:]+):(?P<line>\d+):(\d+)$`,
	"gofmt":       `gofmt -l -s ./*.go:^(?P<path>[^\n]+)$`,
	"goimports":   `goimports -w`,
	"golint":      "golint -set_exit_status ",
	"gotype":      "gotype -e -a ",
	"ineffassign": `ineffassign -n .:PATH:LINE:COL:MESSAGE`,
	"interfacer":  `interfacer `,
	"lll":         `lll -g -l {maxlinelength} ./*.go:PATH:LINE:MESSAGE`,
	"structcheck": `structcheck {tests=-t} .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`,
	"test":        `go test:^--- FAIL: .*$\s+(?P<path>[^:]+):(?P<line>\d+): (?P<message>.*)$`,
	"testify":     `go test:Location:\s+(?P<path>[^:]+):(?P<line>\d+)$\s+Error:\s+(?P<message>[^\n]+)`,
	"varcheck":    `varcheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):[\s\t]+(?P<message>.*)$`,
	"vet":         "go vet ",
	"vetshadow":   "go tool vet --shadow ./*.go:PATH:LINE:MESSAGE",
	"unconvert":   "unconvert -apply",
	"gosimple":    "gosimple ",
	"staticcheck": "staticcheck ",
	"misspell":    "misspell ./*.go:PATH:LINE:COL:MESSAGE",
}

var singlePackageLinters = []string{
	"golint",
	//	"unconvert",
	//	"staticcheck",
	//	"interfacer",
}

var multiPapackageLinter = []string{
	"errcheck",
	"vet",
	"gosimple",
}

var dirLinters = []string{
	"goimports",
	"goconst",
}

//CheckMultiPackages runs linters and code checkers in the passed list of packages
func CheckMultiPackages(listOfPackages []string) (err error) {
	//Find out what the Root Path is
	var tmpErr error
	var errCount int
	log.Debug("Checking pakages...", listOfPackages)
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}
	log.Debug("Root path is:", rootPath)

	for _, linter := range multiPapackageLinter {
		log.Debug("Running linter:", linter)
		if !installer.DisabledTool[linter] {
			log.Debug("Running package checker:", linter)
			cmdString := lintersFlag[linter]
			splitCmd := strings.Split(cmdString, " ")

			cmdList := list.New()

			//Create a linked list of the original command
			for _, aCmd := range splitCmd {
				cmdList.PushBack(aCmd)
			}
			//Add packages to the linked list
			for _, aPackage := range listOfPackages {

				cmdList.PushBack(aPackage)
			}

			//Convert the linked list into a string
			argsList := make([]string, cmdList.Len())
			count := 0
			for e := cmdList.Front(); e != nil; e = e.Next() {
				argsList[count] = e.Value.(string)
				count++

			}

			msg := fmt.Sprintf("LINTER CMD: %s %v", splitCmd[0], argsList)
			log.Debug(msg)
			cmd := exec.Command(splitCmd[0], splitCmd[1])
			cmd.Args = argsList

			out, err := cmd.CombinedOutput()
			if err != nil {
				//Check patterns here
				errList, _ := readErrorsFromChecker(out, linter)
				if len(errList) > 0 {
					errCount++
					tmpErr = fmt.Errorf("%s found errors", linter)
					log.Error(tmpErr)
				} else {
					log.Warnf("%s found Errors but an error ingnore matched", linter)
					err = nil
				}
			}
		} else {
			log.Warn("Could not run disabled tool:", linter)
		}

	}

	if tmpErr != nil {
		err = fmt.Errorf("Found %d package linter errors", errCount)
	}
	return err
}

//CheckSinglePackages runs linters and code checkers in the passed list of packages
func CheckSinglePackages(listOfPackages []string) (err error) {
	//Find out what the Root Path is
	var tmpErr error
	var errCount int
	log.Debug("Checking pakages...", listOfPackages)
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}
	log.Debug("Root path is:", rootPath)
	for _, aPackage := range listOfPackages {
		log.Debug("Checking package:", aPackage)
		for _, linter := range singlePackageLinters {
			if !installer.DisabledTool[linter] {
				log.Debug("Running package checker:", linter)
				cmdString := lintersFlag[linter]
				splitCmd := strings.Split(cmdString, " ")
				msg := fmt.Sprintf("CMD: %s %s %s", splitCmd[0], splitCmd[1], aPackage)
				log.Debug(msg)
				cmd := exec.Command(splitCmd[0], splitCmd[1], aPackage)
				out, err := cmd.CombinedOutput()
				if err != nil {
					//Check patterns here
					errList, _ := readErrorsFromChecker(out, linter)
					if len(errList) > 0 {
						errCount++
						tmpErr = fmt.Errorf("%s found errors", linter)
						log.Error(tmpErr)
					} else {
						log.Warnf("%s found Errors but an error ingnore matched", linter)
						err = nil
					}
				}
			} else {
				log.Warn("Could not run disabled tool:", linter)
			}
		}
	}
	if tmpErr != nil {
		err = fmt.Errorf("Found %d package linter errors", errCount)
	}
	return err
}

func readErrorsFromChecker(cherrs []byte, tool string) (retList []string, err error) {
	var errEOF error
	var aLine string
	var pattern string
	var match = false
	patterns := config.GlobalConfig.IgnorePattern

	listOfPatterns, foundPattern := patterns[tool]
	log.Debug("PATTERNS", listOfPatterns)
	if foundPattern {
		log.Debug("Found patterns for tool:", tool)
		log.Debug("Found patterns for tool:", listOfPatterns)
		for patID, pat := range listOfPatterns {
			if patID == 0 {
				pattern = pat
			} else {
				pattern = fmt.Sprintf("%s|%s", pattern, pat)
			}
		}
		log.Debug("Pattern to exclude:", pattern)
	}

	errList := list.New()
	aReader := bytes.NewReader(cherrs)
	r1 := bufio.NewReader(aReader)
	for errEOF != io.EOF {
		aLine, errEOF = r1.ReadString(10) //  line was defined before
		aLine = strings.Trim(aLine, "\n")
		if foundPattern {
			match, _ = regexp.MatchString(pattern, aLine)
			if match {
				log.Debug("------>>>> MATCH:", aLine)
				log.Debug("------>>>> PATTERN:", pattern)
			}
		}
		if !match && aLine != "" {
			errList.PushBack(aLine)
			log.Errorf("%s:%s", tool, aLine)
		} else {
			if aLine != "" {
				log.Warnf("%s:%s", tool, aLine)
			}
		}
		match = false
	}
	retList = make([]string, errList.Len())
	count := 0
	for e := errList.Front(); e != nil; e = e.Next() {
		retList[count] = e.Value.(string)
		count++
	}
	return retList, err
}

//CheckDirs runs linters and checkers on directories provided in listOfDirs
func CheckDirs(listOfDirs []string) (err error) {
	var errCount int
	var tmpErr error
	//Find out what the Root Path is
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}

	for _, aDir := range listOfDirs {
		log.Debug("Checking dir:", aDir)
		for _, linter := range dirLinters {
			if !installer.DisabledTool[linter] {
				log.Debug("Running dir checker:", linter)
				cmdString := lintersFlag[linter]
				splitCmd := strings.Split(cmdString, " ")
				var msg string
				var cmd *exec.Cmd
				if len(splitCmd) == 2 {
					msg = fmt.Sprintf("DIR CMD: %s %s %s", splitCmd[0], splitCmd[1], aDir)
					log.Debug(msg)
					cmd = exec.Command(splitCmd[0], splitCmd[1], aDir)
				} else if len(splitCmd) == 1 {
					msg = fmt.Sprintf("DIR CMD: %s %s ", splitCmd[0], aDir)
					log.Debug(msg)
					cmd = exec.Command(splitCmd[0], aDir)
				} else if len(splitCmd) == 3 {
					msg = fmt.Sprintf("DIR CMD: %s %s %s %s", splitCmd[0], splitCmd[1], splitCmd[2], aDir)
					log.Debug(msg)
					cmd = exec.Command(splitCmd[0], aDir)
				} else if len(splitCmd) == 4 {
					msg = fmt.Sprintf("DIR CMD: %s %s %s %s %s", splitCmd[0], splitCmd[1], splitCmd[2], splitCmd[3], aDir)
					log.Debug(msg)
					cmd = exec.Command(splitCmd[0], aDir)
				} else {
					log.Error("CMD:", splitCmd[0])
					log.Error("LENTH:", len(splitCmd))
					log.Fatal("Handle this error")
				}

				out, err := cmd.CombinedOutput()
				if err != nil {
					//Check patterns here
					errList, _ := readErrorsFromChecker(out, linter)
					if len(errList) > 0 {
						errCount++
						tmpErr = fmt.Errorf("%s found errors", linter)
						log.Error(tmpErr)
					} else {
						log.Warnf("%s found Errors but an error ingnore matched", linter)
						err = nil
					}
				}
				/*
					if err != nil {
						log.Debug("args:", cmd.Args)
						err = fmt.Errorf("%s found error in:%s", linter, string(out))
						log.Error(err)
						return err
					}
				*/
			} else {
				log.Warn("Could not run disabled tool:", linter)
			}

		}

	}
	if tmpErr != nil {
		err = fmt.Errorf("Found %d package linter errors", errCount)
	}
	return nil
}
