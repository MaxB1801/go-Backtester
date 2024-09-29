package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rivo/tview"
)

// NO GPT, EXCEPT FOR HELP WITH RSI MATHS LOGIC
// aims of this code. Use tview for configurables of how to do back test. For now set up with the constants for back testing
// improve code logic for entrys so entrys are neat little functions that must all become true for entry same with exits
// sort lev, accs,dis and mids all into seperate folders for backtests, will be a tview configurable
// add starting pot function for testing that way, might need to be a seperacode

// Create a channel for messages
var channel chan string = make(chan string)
var errorChannel chan string = make(chan string)

func main() {

	for {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		selected := ""

		menu := tview.NewApplication()

		modal := tview.NewModal().
			SetText("Main Menu").
			AddButtons([]string{"Backtest", "Trade Tracker", "Trade Graphs", "Quit"})
		modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Backtest" {
				selected = "Backtest"
				menu.Stop()

			}
			if buttonLabel == "Trade Tracker" {
				// Set the environment variable
				err = os.Setenv("PERIOD", "1")
				if err != nil {
					log.Fatalf("Failed to set env var: %v", err)
				}
				err = os.Setenv("PYTEMP", "1")
				if err != nil {
					log.Fatalf("Failed to set env var: %v", err)
				}

				selected = "Trade Tracker"

				menu.Stop()

			}
			if buttonLabel == "Trade Graphs" {
				selected = "Trade Graphs"
				menu.Stop()
			}
			if buttonLabel == "Quit" {
				selected = "Quit"
				menu.Stop()
			}
		})
		if err := menu.SetRoot(modal, false).EnableMouse(true).SetFocus(modal).Run(); err != nil {
			panic(err)
		}

		if selected == "Quit" {
			break
		}
		switch selected {
		case "Backtest":
			backtestDashboard()
		case "Trade Tracker":
			list, err := os.ReadDir(fmt.Sprintf(dir + "\\tmp"))
			if err != nil {
				log.Fatal(err, list)
			}
			formatEntrys(dir, list)
		case "Trade Graphs":
			graphDashboard(dir)
		}
	}

}
