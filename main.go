package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

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

	out, err := exec.Command("pm2", "jlist").Output()
	if err != nil {
		log.Fatal(err)
	}
	var processList []struct {
		Name string `json:"name"`
		Env  struct {
			Status string `json:"status"`
		} `json:"pm2_env"`
	}
	json.Unmarshal(out, &processList)
	for _, processItem := range processList {
		tooltip := fmt.Sprintf("Enable / disable %s process", processItem.Name)
		process := Process{
			processItem.Name,
			systray.AddMenuItemCheckbox(processItem.Name, tooltip, processItem.Env.Status == "online"),
		}
		processes = append(processes, process)
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
