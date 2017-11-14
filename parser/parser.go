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

//Package parser returns two lists:
// 1. A list of unique directories based on the passed list of modified files.
// 2. A list of unique packages based on the passed list of modified files.
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

//Parse parses the provided list of modified files
func Parse(stringList []string) (dir []string, pkag []string, err error) {

	/*
		//Find out what the Root Path is
		cmd := exec.Command("git", "rev-parse", "--show-toplevel")
		tmpRootPath, err := cmd.Output()
		if err != nil {
			return dir, pkag, err
		}

		//Trim return character
		rootPath := strings.TrimSpace(string(tmpRootPath))
		err = os.Chdir(rootPath)
		if err != nil {
			return dir, pkag, err
		}
	*/
	//log.Debug("Parser:",stringList)
	dir, err = getUniqueDirs(stringList)

	log.Info("DIR to parse", dir)

	pkag, err = getUniquePkgs(dir)

	return dir, pkag, err
}

func getUniquePkgs(dirList []string) (pkgList []string, err error) {
	pkgHash := make(map[string]bool)
	rawPkagList, err := readListOfPackages()
	fuulRootPackage := rawPkagList[0]
	r, err := regexp.Compile(`\w+$`)
	rootPackage := r.FindString(fuulRootPackage)
	for _, key := range dirList {
		var fullKey string
		if key == "." {
			fullKey = rootPackage
		} else {
			fullKey = rootPackage + "/" + key
		}
		for _, aPackage := range rawPkagList {
			match, _ := regexp.MatchString(fmt.Sprintf("%s(/.)?$", fullKey), aPackage)
			if match {
				pkgHash[aPackage] = true
			}
		}

	}

	pkgList = make([]string, len(pkgHash))
	count := 0
	for key := range pkgHash {
		pkgList[count] = key
		count++
	}
	return pkgList, nil
}

func getUniqueDirs(stringList []string) (dir []string, err error) {
	dirsHash := make(map[string]bool)
	splitString := stringList //strings.Split(stringList, " ")
	log.Debug("Stringlist:", splitString)
	for _, file := range splitString {
		path := filepath.Dir(file)
		dirsHash[path] = true

	}
	dirList := make([]string, len(dirsHash))
	count := 0
	for key := range dirsHash {
		dirList[count] = key
		count++
	}
	return dirList, nil

}

func readListOfPackages() (pkag []string, err error) {

	//-------------------------------------------
	//          Read list of packages
	//--------------------------------------------
	golistCmd := exec.Command("go", "list", "./...")
	tmpGoList, err := golistCmd.Output()
	if err != nil {
		log.Error("Parser failed in go list")
		return pkag, err
	}
	//fmt.Println(string(tmpGoList))
	aReader := bytes.NewReader(tmpGoList)
	bReader := bytes.NewReader(tmpGoList)
	r1 := bufio.NewReader(aReader)
	r2 := bufio.NewReader(bReader)
	count := 0

	var errEOF error
	_, err = r1.ReadString(10) // line defined once
	for errEOF != io.EOF {
		_, errEOF = r1.ReadString(10) //  line was defined before
		count++
	}

	packagesList := make([]string, count)

	count = 0
	line2, errEOF := r2.ReadString(10) // line defined once
	//Trim return character
	aLine := strings.TrimSpace(line2)
	for errEOF != io.EOF {
		packagesList[count] = aLine
		aLine, errEOF = r2.ReadString(10) //  line was defined before
		aLine = strings.TrimSpace(aLine)
		count++
	}
	//fmt.Println(packagesList)
	pkag = packagesList
	return packagesList, nil
}
