package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, func() {})
}

type Process struct {
	Name     string
	MenuItem *systray.MenuItem
}

var processes []Process

func onReady() {
	iconData, _ := base64.StdEncoding.DecodeString(getFiles()["pm2.ico"])
	systray.SetIcon(iconData)
	systray.SetTitle("Pm2 systray")
	mQuit := systray.AddMenuItem("Exit", "Quit pm2 systray")
	systray.AddSeparator()
	mQuit.SetIcon(iconData)

	out, err := exec.Command("pm2", "l").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(out), "\n")
	for x, line := range lines {
		if x > 2 && x < len(lines)-2 {
			status := strings.Contains(line, "online")
			name := strings.Split(line[12:], " ")[0]
			tooltip := fmt.Sprintf("Enable / disable %s process", name)
			process := Process{
				name,
				systray.AddMenuItemCheckbox(name, tooltip, status),
			}
			printProcess(process, "[I]")

			processes = append(processes, process)
		}
	}
	systray.AddSeparator()
	mSave := systray.AddMenuItem("save", "Save current configuration")

	go registerExit(*mQuit)
	go registerSave(*mSave)
	for _, process := range processes {
		go registerProcess(process)
	}
}

func registerExit(mQuit systray.MenuItem) {
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}
func registerSave(mSave systray.MenuItem) {
	for {
		select {
		case <-mSave.ClickedCh:
			exec.Command("pm2", "save").Run()
		}
	}
}
func registerProcess(process Process) {
	for {
		select {
		case <-process.MenuItem.ClickedCh:
			printProcess(process, "[B]")
			if process.MenuItem.Checked() {
				process.MenuItem.Uncheck()
			} else {
				process.MenuItem.Check()
			}
			toggleProcess(process, process.MenuItem.Checked())
			printProcess(process, "[A]")
		}
	}
}
func toggleProcess(process Process, enable bool) {
	var action string
	if enable {
		action = "start"
	} else {
		action = "stop"
	}
	exec.Command("pm2", action, process.Name).Run()
}

func printProcess(process Process, prefix string) {
	fmt.Printf("%s %v | %s", prefix, process.MenuItem.Checked(), process.Name)
}
