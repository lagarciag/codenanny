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

//installer does this blah blah
package dirlister_test

import (
	"os"
	"testing"
	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/dirlister"
	"strings"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
	v := t.Run()
	os.Exit(v)

}

func TestDirListerBasic(t *testing.T) {
	rawGopath := os.Getenv("GOPATH")
	splitGoPath := strings.Split(rawGopath,":")
	gopath := splitGoPath[len(splitGoPath)-1]
	if err := os.Chdir(gopath + "/src/github.com/lagarciag/codenanny/"); err != nil {
		t.Error("Could not change dir")
	}

	t.Log("TestDirLister run")


	dirListSlice, dirListString, err := dirlister.ListDir("./")

	if err != nil {
		t.Error("Dir list error:",err.Error())
	}
	if (len(dirListSlice) != 13) {
		t.Error("List size should be 12, but it is:",len(dirListSlice))
	}
	t.Log("dirList:",dirListSlice)
	t.Log("dirString:",dirListString)

}

