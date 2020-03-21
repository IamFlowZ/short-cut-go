package main

import (
	"fmt"
	"log"
	"strings"
	"os"
	"regexp"
	"flag"
	"text/tabwriter"
)

var homePath string
var bashrcPath string
const shortcutsPath string = "/var/shortcut/shortcuts"
var shortcutName string
const helpText = `
         __               __                   __ 
   _____/ /_  ____  _____/ /_      _______  __/ /_
  / ___/ __ \/ __ \/ ___/ __/_____/ ___/ / / / __/
 (__  ) / / / /_/ / /  / /_/_____/ /__/ /_/ / /_  
/____/_/ /_/\____/_/   \__/      \___/\__,_/\__/  

Usage:
	short-cut [NAME]

	Quickly create an environment variable that equals the current working directory.

	Example:
		[user@machine:/home/user]$ short-cut test
			output: Please reload your bashrc file to access the new shortcut. (source ~/.bashrc)
		[user@machine:/home/user]$ source ~/.bashrc
		[user@machine:/home/user]$ echo $test
			output: /home/user
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

func writeShortcut(shortcutName string, currentDir string) (error) {
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

func readShortcuts() ([]string, error) {
	if file, err := os.Open(shortcutsPath); err != nil {
		file.Close()
		return nil, err
	} else {
		data := make([]byte, 1000)
		if count, err := file.Read(data); err != nil {
			file.Close()
			return nil, err
		} else {
			lines := string(data[:count])
			file.Close()
			return strings.Split(lines, "\n"), nil
		}
	}
	
}

func listShortcuts() {
	w := tabwriter.NewWriter(os.Stdout, 10, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Shortcut:\tPath:\t")
	fmt.Fprintln(w, "------------------------------")
	splitLines, err := readShortcuts()
	if err != nil {
		log.Fatal("Couldn't read shortcuts, ", err)
	}
	for _, line := range splitLines {
		if len(line) > 0 {
			split := strings.Split(line, "=")
			split[0] = split[0][7:]
			s := fmt.Sprintf("%s\t%s\t", split[0], split[1])
			fmt.Fprintln(w, s)
			fmt.Fprintln(w, "------------------------------")
			w.Flush()
		}
	}
}

func removeShortcut(shortcutName string) (error) {
	splitLines, err := readShortcuts()
	newLines := make([]string, len(splitLines))
	for i, line := range splitLines {
		if strings.Contains(line, shortcutName) {
			fmt.Println("fuck yea")
		} else {
			newLines[i] = line
		}
	}
	if file, err := os.OpenFile(shortcutsPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	} else {
		for _, line := range newLines {
			if len(line) > 0 {
				if _, err := file.Write([]byte(line + "\n")); err != nil {
					return err
				}
			}
		}
	}
	
	return err
}

var list bool
var help bool
var remove bool
var removeAll bool
func init() {
	flag.BoolVar(&help, "help", false, "Display the help text")
	flag.BoolVar(&list, "list", false, "Display a list of the available shortcuts")
	flag.BoolVar(&remove, "remove", false, "Remove a shortcut")
	flag.BoolVar(&removeAll, "removeAll", false, "Remvoe all shortcuts. \n**WARNING** This flag does not work.")
}

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		fmt.Println("Please enter a name for the shortcut you would like.")
		os.Exit(1)
	}

	if help {
		fmt.Println(helpText)
		os.Exit(0)
	}

	if list {
		listShortcuts()
		os.Exit(0)
	}

	if removeAll {
		if err := os.Remove(shortcutsPath); err != nil {
			fmt.Println(err)
		}
	}
	
	if strings.Contains(os.Args[1], "-") {
		shortcutName = os.Args[2]
	} else {
		shortcutName = os.Args[1]
	}

	if remove {
		removeShortcut(shortcutName)
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
	if err := writeShortcut(shortcutName, currentDir); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Please reload your bashrc file to access the new shortcut. (source ~/.bashrc)")
}