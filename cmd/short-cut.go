package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

var homePath string
var bashrcPath string
var shortcutsPath string
var shortcutName string
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

func setupBashrc(currentDir string) (error) {
	homePath, _ := os.UserHomeDir()
	bashrcPath := homePath + "/.bashrc"
	_, exists := os.LookupEnv("SC_LOADED")
	if !exists {
		if bashrc, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		} else {
			bashLines := "\n# loading short-cut environment variables. see: short-cut --help\nexport SC_LOADED=true\n. /var/shortcut/shortcuts"
			if _, err := bashrc.Write([]byte(bashLines)); err != nil {
				return err
			}
		}
	}
	return nil
}

func setupShortcuts(shortcutName string, currentDir string) (error) {
	shortcutsPath := "/var/shortcut/shortcuts"
	shortCut := "export " + shortcutName + "=" + currentDir
	if file, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		if _, err := file.Write([]byte(shortCut + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please enter a name for the shortcut you would like.")
		os.Exit(1)
	}
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

	currentDir, _ := os.Getwd()
	if err := setupBashrc(currentDir); err != nil {
		log.Fatal(err)
	}
	if err := setupShortcuts(shortcutName, currentDir); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please reload your bashrc file to access the new shortcut. (source ~/.bashrc)")
}