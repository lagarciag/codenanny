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

//Package dirlister contains structs and methods for finding and listing go files
package dirlister

import (
	"container/list"
	"path/filepath"

	log "github.com/Sirupsen/logrus"

	"os"
	"regexp"
)

//ListDir returns in multiple formats a list of go files found recursively in the provided directory.
func ListDir(path string) (dirListSlice []string, dirListString string, err error) {
	log.Debug("Dir lister:", path)
	//------------------------------
	//Create a list
	//-------------------------------
	dirList := list.New()

	//-----------------------------------
	// Helper closure for filepath.Walk
	//-----------------------------------
	visit := func(path string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString((".go$"), path)
		if match {
			dirList.PushBack(path)
		}
		return err
	}

	//---------------------------------
	// Do the walk
	//----------------------------------
	err = filepath.Walk(path, visit)
	if err != nil {
		return dirListSlice, dirListString, err
	}
	log.Debug("Elements:", dirList.Len())
	size := dirList.Len()

	//-----------------------------------
	//Iterate the list
	//-----------------------------------
	dirListSlice = make([]string, size)
	dirListString = ""
	count := 0
	for e := dirList.Front(); e != nil; e = e.Next() {
		log.Debug("File:", e.Value)
		dirListSlice[count] = e.Value.(string)
		dirListString = dirListString + e.Value.(string) + " "
	}

	return dirListSlice, dirListString, err
}
