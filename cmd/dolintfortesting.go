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

//Package cmd holds multiple cobra commands
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/config"
)

//DoLintForTesting is a service function for implementing codenanny tests
func DoLintForTesting(conf config.CodeNannyConfig) (err error) {
	config.GlobalConfig = conf

	//Find out what the Root Path is
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("Could not find git root for this repo:%s", err.Error())
		log.Errorf("Is this a git repo???:%s", err.Error())
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		err = fmt.Errorf("Could not change to repos root dir:%s", err.Error())
		return err
	}
	log.Debug("Dolint for testing...")
	err = Lintdir("./")
	return err
}
