package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

var homePath string
var bashrcPath string
var shortcutsPath string
var shortcutName string

const bashLines = "\n# loading short-cut environment variables. see: short-cut --help\nexport SC_LOADED=true\n. ~/.shortcuts"
const helpText = `
         __               __                   __ 
   _____/ /_  ____  _____/ /_      _______  __/ /_
  / ___/ __ \/ __ \/ ___/ __/_____/ ___/ / / / __/
 (__  ) / / / /_/ / /  / /_/_____/ /__/ /_/ / /_  
/____/_/ /_/\____/_/   \__/      \___/\__,_/\__/  

Usage:
	short-cut [NAME]
	Quickly create an env variable that equates to the current working directory.
	(e.g. "user@machine:PATH$ short-cut test" = "$test=PATH")

Options:
	-h, --help 	Display this help text
`

func main() {
	shortcutName := os.Args[1]
	if shortcutName == "--help" || shortcutName == "-h" {
		fmt.Println(helpText)
		os.Exit(0)
	}

	re := regexp.MustCompile(`^[a-zA-z0-9]*$`)
	if validName := re.MatchString(shortcutName); !validName {
		log.Fatal("Invalid shortcut name")
		os.Exit(1)
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	shortcutsPath := homePath + "/.shortcuts"
	bashrcPath := homePath + "/.bashrc"
	currentDir, err := os.Getwd()
	shortCut := "export " + shortcutName + "=" + currentDir

	_, exists := os.LookupEnv("SC_LOADED")
	if !exists {
		if bashrc, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			log.Fatal(err)
		} else {
			if _, err := bashrc.Write([]byte(bashLines)); err != nil {
				log.Fatal(err)
			}
		}
	}

	if file, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	} else {
		if _, err := file.Write([]byte(shortCut + "\n")); err != nil {
			log.Fatal(err)
		}
	}
	
	cmd := exec.Command("bash", "./reload.sh")
	cmdErr := cmd.Run()
	if cmdErr != nil {
		log.Fatal(cmdErr)	
	}
}