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
package installer_test

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/lagarciag/codenanny/installer"
)

func TestMain(t *testing.M) {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
	v := t.Run()
	os.Exit(v)

}

func TestInstallerBasic(t *testing.T) {

	if err := installer.CheckExternalDependencies(); err != nil {
		t.Error("Could not install package:", err)
	}

}
