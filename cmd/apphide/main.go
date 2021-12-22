package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rkoesters/xdg/desktop"
)

const (
	appPath     = "/usr/share/applications"
	hiddenEntry = "[Desktop Entry]\nHidden=true"
)

var home = os.Getenv("HOME")

var (
	idMode = flag.Bool("id", false, "Match the id instead of the name of the application")
	unhide = flag.Bool("uh", false, "Unhide the application")
)

type App struct {
	ID   string
	Name string
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Please specify a regular expression of the name of the application(s) to hide")
		return
	} else if flag.NArg() > 1 {
		fmt.Println("You can only specify one regular expression")
		return
	}

	nameRegexp, err := regexp.Compile(flag.Arg(0))
	if err != nil {
		fmt.Println("The regular expression you provided is invalid:", err)
		return
	}

	apps, err := os.ReadDir(appPath)
	if err != nil {
		fmt.Println("Could not read applications directory:", err.Error())
		return
	}

	var toHide []App

	for _, app := range apps {
		if !strings.HasSuffix(app.Name(), ".desktop") {
			continue
		}

		rawEntry, err := os.Open(appPath + "/" + app.Name())
		if err != nil {
			fmt.Println("Could not open application file, skipping:", err)
			continue
		}

		entry, err := desktop.New(rawEntry)
		_ = rawEntry.Close()
		if err != nil {
			continue
		}

		id := app.Name()[:len(app.Name())-len(".desktop")]

		if *idMode {
			if nameRegexp.MatchString(id) {
				toHide = append(toHide, App{ID: id, Name: entry.Name})
			}
		} else {
			if nameRegexp.MatchString(entry.Name) {
				toHide = append(toHide, App{ID: id, Name: entry.Name})
			}
		}
	}

	if len(toHide) == 0 {
		fmt.Println("No applications matched the regular expression")
		return
	}

	fmt.Printf("%d applications matched the regular expression:\n", len(toHide))
	for _, app := range toHide {
		fmt.Printf("• %s (%s)\n", app.Name, app.ID)
	}

	fmt.Println("Remove them? (Y/n)")

	r := bufio.NewReader(os.Stdin)

	answer, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Could not read input:", err)
		return
	}

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "" && answer != "y" {
		fmt.Println("Aborting")
		return
	}

	for _, app := range toHide {
		if *unhide {
			err = os.Remove("~/.local/share/applications/" + app.ID + ".desktop")
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					fmt.Printf("%s is not hidden, skipping\n", app.Name)
				} else {
					fmt.Println("Could not remove hidden application:", err)
				}

				continue
			}
		} else {
			file, err := os.Create(home + "/.local/share/applications/" + app.ID + ".desktop")
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					fmt.Printf("Could not hide %s as there is already an entry at "+
						"~/.local/share/applications/%s.desktop, skipping\n", app.Name, app.ID)
				} else {
					fmt.Println("Could not hide application:", err)
				}

				continue
			}

			_, err = io.WriteString(file, hiddenEntry)
			_ = file.Close()
			if err != nil {
				fmt.Printf("Could not create desktop entry to hide %s: %s\n", app.Name, err)
				continue
			}
		}
	}
}
