package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/getlantern/systray"
	"log"
	"os/exec"
	"sort"
)

//go:embed pm2.ico
var icon []byte

type WidgetProcess struct {
	Name     string
	MenuItem *systray.MenuItem
}

type ModelProcess struct {
	Name string `json:"name"`
	Env  struct {
		Status string `json:"status"`
	} `json:"pm2_env"`
}

var processes []WidgetProcess

func main() {
	systray.Run(onReady, func() {})
}

func onReady() {
	processList := getOrderedProcesses()
	render(processList)

	for _, process := range processes {
		go registerProcess(process)
	}
}

func getOrderedProcesses() []ModelProcess {
	out, err := exec.Command("pm2", "jlist").Output()
	if err != nil {
		log.Fatal(err)
	}
	var processList []ModelProcess
	err = json.Unmarshal(out, &processList)
	if err != nil {
		log.Fatal("unable to read pm2 information")
	}
	sort.Slice(processList, func(i, j int) bool {
		return processList[i].Name < processList[j].Name
	})
	return processList
}

func render(processList []ModelProcess) {
	systray.SetIcon(icon)
	systray.SetTitle("Pm2 systray")
	mQuit := systray.AddMenuItem("Exit", "Quit pm2 systray")
	systray.AddSeparator()
	mQuit.SetIcon(icon)
	for _, processItem := range processList {
		tooltip := fmt.Sprintf("Enable / disable %s process", processItem.Name)
		process := WidgetProcess{
			processItem.Name,
			systray.AddMenuItemCheckbox(processItem.Name, tooltip, processItem.Env.Status == "online"),
		}
		processes = append(processes, process)
	}
	systray.AddSeparator()
	mSave := systray.AddMenuItem("save", "Save current configuration")
	go registerExit(*mQuit)
	go registerSave(*mSave)

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
			err := exec.Command("pm2", "save").Run()
			if err != nil {
				return
			}
		}
	}
}

func registerProcess(process WidgetProcess) {
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

func toggleProcess(process WidgetProcess, enable bool) {
	var action string
	if enable {
		action = "start"
	} else {
		action = "stop"
	}
	exec.Command("pm2", action, process.Name).Run()
}

func printProcess(process WidgetProcess, prefix string) {
	fmt.Printf("%s %v | %s", prefix, process.MenuItem.Checked(), process.Name)
}
