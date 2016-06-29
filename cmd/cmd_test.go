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

	log "github.com/Sirupsen/logrus"
	"github.com/lagarciag/codenanny/cmd"
	"github.com/lagarciag/codenanny/config"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	formatter := &log.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableTimestamp = true
	log.SetFormatter(&log.TextFormatter{})
	v := t.Run()
	os.Exit(v)

}

func TestNanny(t *testing.T) {
	log.Debug("TestNanny run")
	conf := config.CodeNannyConfig{}
	conf.IgnorePattern = make(map[string][]string)

	conf.IgnorePattern["errcheck"] = []string{
		"lint_test.go:.+:.+CreateUnCheckedError",
		"packagewitherrors.go:.+:.+returnError"}
	conf.IgnorePath = "dirtoexclude"

	if err := cmd.DoLintForTesting(conf); err != nil {
		t.Error(err)
	}

}
