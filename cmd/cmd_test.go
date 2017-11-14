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

//cmd does this blah blah
package cmd_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/lagarciag/codenanny/cmd"
	"github.com/lagarciag/codenanny/config"
	"github.com/lagarciag/codenanny/lint"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableTimestamp = true
	log.SetFormatter(formatter)
	v := t.Run()
	os.Exit(v)

}

func TestNannyNoErrors(t *testing.T) {
	log.Debug("TestNanny run")
	if err := config.LoadConfig(); err != nil {
		t.Error("Error loading configuration")
	}
	if err := lint.ChgDirToGitRootPath(); err != nil {
		t.Error("Chage dir to git root returned error")
	}

	if err := cmd.DoLintForTesting(config.GlobalConfig); err != nil {
		t.Error(err)
	}

}
