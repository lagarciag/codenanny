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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	"vet":         `go vet`,
	"vetshadow":   `go tool vet -shadow=true`,
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

var dirMultiLinters = []string{
	//"goimports",
	"goconst",
}

var dirRecurseiveLinters = []string{
	//"goimports",
	//"goconst",
	"vetshadow",
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

			out, errOut := cmd.CombinedOutput()
			if errOut != nil {
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
				out, errOut := cmd.CombinedOutput()
				if errOut != nil {
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
		//log.Info("Found patterns for tool:", tool)
		//log.Info("Found patterns for tool:", listOfPatterns)
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

//CheckMultiDirs runs linters and checkers on directories provided in listOfDirs
func CheckMultiDirs(listOfDirs []string) (err error) {
	var errCount int
	var tmpErr error

	//------------------------------
	//change dir into git root path
	//------------------------------
	//if err = chgDirToGitRootPath();err != nil {
	//	return err
	//}

	//------------------------------
	// Iterate through each passed
	// directory
	//------------------------------
	for _, aDir := range listOfDirs {
		log.Debug("Checking dir:", aDir)
		//---------------------------------
		// To each directory run a checker
		//---------------------------------
		for _, checker := range dirMultiLinters {
			if !installer.DisabledTool[checker] {
				var lCmd *exec.Cmd
				log.Debug("Running dir checker:", checker)
				cmdString := lintersFlag[checker]
				splitCmd := strings.Split(cmdString, " ")
				switch {
				case len(splitCmd) == 1:
					lCmd = exec.Command(splitCmd[0], aDir)
				case len(splitCmd) == 2:
					lCmd = exec.Command(splitCmd[0], splitCmd[1], aDir)
				case len(splitCmd) == 3:
					lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], aDir)
				case len(splitCmd) == 4:
					lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], splitCmd[3], aDir)
				case len(splitCmd) == 5:
					lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], splitCmd[3], splitCmd[4], aDir)
				case len(splitCmd) == 0 || len(splitCmd) > 5:
					log.Error("This command is not correctly handled", splitCmd[0])
					log.Error("LENTH:", len(splitCmd))
					log.Fatal("Handle this error")
				}

				//---------------------------------------------
				//                Execute command
				//---------------------------------------------
				out, errOut := lCmd.CombinedOutput()

				//---------------------------------
				// Handle errors, if there are any
				//---------------------------------
				if errOut != nil {
					//-------------------------------------
					// Iterate through erros and verify
					// if any has an exclusion
					//-------------------------------------

					errList, _ := readErrorsFromChecker(out, checker)
					if len(errList) > 0 {
						errCount++
						tmpErr = fmt.Errorf("%s found errors", checker)
						log.Error(tmpErr)
					} else {
						log.Warnf("%s found Errors but an error ingnore matched", checker)
						err = nil
					}
				}

			} else {
				log.Warn("Could not run disabled tool:", checker)
			}

		}

	}
	if tmpErr != nil {
		err = fmt.Errorf("Found %d package linter errors", errCount)
	}
	return nil
}

//CheckRecursiveDirs runs linters and checkers on directories provided in listOfDirs
func CheckRecursiveDirs(listOfDirs []string) (err error) {
	var errCount int
	var tmpErr error
	var theDir string

	log.Debug("List of Dirs:", listOfDirs)

	if len(listOfDirs) == 0 {
		log.Warn("CheckRecursiveDir List is empty")
		return nil
	}

	if len(listOfDirs) == 1 {
		theDir = listOfDirs[0]
	} else {
		theDir = "./"
	}

	//---------------------------------
	// To each directory run a checker
	//---------------------------------
	for _, checker := range dirRecurseiveLinters {
		if !installer.DisabledTool[checker] {
			var lCmd *exec.Cmd
			log.Debug("Running dir checker:", checker)
			cmdString := lintersFlag[checker]
			splitCmd := strings.Split(cmdString, " ")
			switch {
			case len(splitCmd) == 1:
				lCmd = exec.Command(splitCmd[0], theDir)
			case len(splitCmd) == 2:
				lCmd = exec.Command(splitCmd[0], splitCmd[1], theDir)
			case len(splitCmd) == 3:
				lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], theDir)
			case len(splitCmd) == 4:
				lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], splitCmd[3], theDir)
			case len(splitCmd) == 5:
				lCmd = exec.Command(splitCmd[0], splitCmd[1], splitCmd[2], splitCmd[3], splitCmd[4], theDir)
			case len(splitCmd) == 0 || len(splitCmd) > 5:
				log.Error("This command is not correctly handled", splitCmd[0])
				log.Error("LENTH:", len(splitCmd))
				log.Fatal("Handle this error")
			}

			//---------------------------------------------
			//                Execute command
			//---------------------------------------------
			out, errOut := lCmd.CombinedOutput()

			//---------------------------------
			// Handle errors, if there are any
			//---------------------------------
			if errOut != nil {
				//-------------------------------------
				// Iterate through erros and verify
				// if any has an exclusion
				//-------------------------------------
				//log.Info("Errors found for checker:",checker)
				errList, _ := readErrorsFromChecker(out, checker)
				if len(errList) > 0 {
					errCount++
					tmpErr = fmt.Errorf("%s found errors", checker)
					log.Error(tmpErr)
				} else {
					log.Warnf("%s found Errors but an error ingnore matched", checker)
					err = nil
				}
			}

		} else {
			log.Warn("Could not run disabled tool:", checker)
		}

	}
	if tmpErr != nil {
		err = fmt.Errorf("Found %d package linter errors", errCount)
	}
	return nil
}

//ChgDirToGitRootPath chages current working dir to gits repo root
func ChgDirToGitRootPath() (err error) {

	//-------------------------------------
	// Detect the git root path
	//-------------------------------------
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//-------------------------------------
	// cd into the git root path
	//-------------------------------------
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	return err
}
