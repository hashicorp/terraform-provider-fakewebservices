// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fws

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// Credentials -
type Credentials struct {
	Credentials map[string]Credential `json:"credentials"`
}

// Credential -
type Credential struct {
	Token string `json:"token"`
}

// ======================================
// TODO: MAKE WORK ON WINDOWS TOO (START)
// Cribbed from github.com/hashicorp/terraform-svchost
// ======================================

func configFile() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "credentials.tfrc.json"), nil
}

func configDir() (string, error) {
	dir, err := homeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, ".terraform.d"), nil
}

func homeDir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		// FIXME: homeDir gets called from globalPluginDirs during init, before
		// the logging is setup.  We should move meta initializtion outside of
		// init, but in the meantime we just need to silence this output.
		//log.Printf("[DEBUG] Detected home directory from env var: %s", home)

		return home, nil
	}

	// If that fails, try build-in module
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	if user.HomeDir == "" {
		return "", errors.New("blank output")
	}

	return user.HomeDir, nil
}

// ==================================
// TODO: END MAKE WORK ON WINDOWS TOO
// ==================================

func cliCredentials() *Credentials {
	creds := new(Credentials)

	// Detect the CLI config file path.
	configFilePath := os.Getenv("TERRAFORM_CONFIG")
	if configFilePath == "" {
		filePath, err := configFile()
		if err != nil {
			log.Printf("[ERROR] Error detecting default CLI config file path: %s", err)
			return creds
		}
		configFilePath = filePath
	}

	// Read the CLI config file content.
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Printf("[ERROR] Error reading the CLI config file %s: %v", configFilePath, err)
		return creds
	}

	jsonErr := json.Unmarshal(content, creds)
	if jsonErr != nil {
		log.Printf("[ERROR] Error unmarshalling the CLI config file %s: %v", configFilePath, jsonErr)
		return creds
	}

	return creds
}
