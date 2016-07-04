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

//Package config encapsulates configuration functionality for codenanny
package config

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"strconv"
)

//Version holds the current version
var Version = "0.0.2"

//GlobalConfig holds configuration functionality for Codennay
var GlobalConfig CodeNannyConfig

//CodeNannyConfig is the struct used to marshall in the configuration
type CodeNannyConfig struct {
	Version       string              `yaml:"required_version"`
	Disabled      map[string]bool     `yaml:"disabled"`
	IgnorePattern map[string][]string `yaml:"ignore_pattern"`
	IgnorePath    string              `yaml:"ignore_path_pattern"`
}

//LoadConfig loads and processes the configuration file
func LoadConfig() (err error) {
	var yamlFile []byte

	//Find out what the Root Path is
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		log.Errorf("Is this a git repo???:%s", err.Error())
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}

	configFile := rootPath + "/.codenanny"
	if _, err := os.Stat(configFile); err == nil {
		log.Debug("Found Config file")
		yamlFile, err = ioutil.ReadFile(configFile)
		if err != nil {
			panic("Could not load file")
		}
		//Unmarshal yaml file into allocated go struct
		if err := yaml.Unmarshal(yamlFile, &GlobalConfig); err != nil {
			return err
		}
		//Check Version
		if err := CheckVersion(); err != nil {
			return err
		}

		log.Debug(GlobalConfig)
	} else {
		log.Debug("No config file cound in:", configFile)
	}

	return nil
}

func versionToInt(ver string) (verInt int, err error) {
	slicedVersion := strings.Split(ver, ".")
	concatVersion := slicedVersion[0] + slicedVersion[1] + slicedVersion[2]
	verInt, err = strconv.Atoi(concatVersion)
	return verInt, err
}

//CheckVersion checks that .codenanny required version matches the actual version
func CheckVersion() (err error) {
	if GlobalConfig.Version != "" {
		requiredVersion, errV := versionToInt(GlobalConfig.Version)
		if err != nil {
			return errV
		}
		actualVersion, errV := versionToInt(Version)
		if err != nil {
			return errV
		}
		if requiredVersion > actualVersion {
			//log.Errorf("The .codenanny file (configuration) says that codenanny version %s is required, but this is %s ",GlobalConfig.Version,Version)
			log.Fatalf("The .codenanny file (configuration) says that codenanny version %s is required, but this is %s ", GlobalConfig.Version, Version)
		}
		log.Info("Required version:", requiredVersion)
		log.Info("Actual version:", actualVersion)
	}
	return err
}
