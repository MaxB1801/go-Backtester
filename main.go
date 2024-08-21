package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
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

	var exit = true
	var downloadX = false
	var download = false
	var err error
	var stop_loss bool
	var rsi_low int
	var rsi_high int
	var rsi_exit_low int
	var rsi_exit_high int
	var folder string
	var rsi_increment int
	var buttonText string

	// creating loading screen while downloading apps, can press button to skip!!!
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	app1 := tview.NewApplication()

	modal := tview.NewModal().
		SetText("Would you like to download the test data to ensure it is current? The data will be retrieved from Yahoo Finance.").
		AddButtons([]string{"Download", "Download X Years", "Skip"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Download" {
			download = true
			downloadX = false

			app1.Stop()

		}
		if buttonLabel == "Download X Years" {
			download = true
			downloadX = true

			app1.Stop()

		}
		if buttonLabel == "Skip" {
			download = true
			download = false
			downloadX = false
			app1.Stop()
		}
	})
	if err := app1.SetRoot(modal, false).EnableMouse(true).SetFocus(modal).Run(); err != nil {
		panic(err)
	}

	if downloadX {
		app := tview.NewApplication()

		inputField := tview.NewInputField()
		inputField.SetLabel("Retrieve Data for the Past X Years: ").
			SetFieldWidth(10).
			SetAcceptanceFunc(tview.InputFieldInteger).
			SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					// Get the input value
					input := inputField.GetText()
					// Convert input to an integer
					years, err := strconv.Atoi(input)
					if err != nil {
						inputField.SetLabel("Retrieve Data for the Past X Years (Enter an Integer): ")
					} else {
						// Set the environment variable
						err = os.Setenv("PERIOD", strconv.Itoa(years))
						if err != nil {
							log.Fatalf("Failed to set env var: %v", err)
						}

						app.Stop()
					}
				}
			})

		if err := app.SetRoot(inputField, true).SetFocus(inputField).Run(); err != nil {
			panic(err)
		}
	}

	if download {
		wait := tview.NewApplication()
		go Downloader(wait, dir)

		waiter := tview.NewModal().
			SetText("Downloading: Skipping now may cause issues").
			AddButtons([]string{"Skip"})
		waiter.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Skip" {
				download = false
				wait.Stop()
			}
		})
		if err := wait.SetRoot(waiter, false).EnableMouse(true).SetFocus(modal).Run(); err != nil {
			panic(err)
		}
	}

	app := tview.NewApplication()

	// create flex to center
	message := tview.NewList()
	// error channels texts
	errors := tview.NewList()
	flexMess := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(message, 0, 2, false).
		AddItem(tview.NewBox(), 0, 1, false)
	flexMess.SetBorder(true).SetTitle("Results").SetBorderColor(tcell.ColorDarkCyan.TrueColor())
	flexErrors := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(errors, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)
	flexErrors.SetBorder(true).SetTitle("Error Messages").SetBorderColor(tcell.ColorRed.TrueColor())
	// making flex
	flextText := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexMess, 0, 2, false).
		AddItem(flexErrors, 0, 2, false)

	// define button here
	buttonSort := tview.NewButton("Start Backtests")
	// consume messages
	go recieve(message, errors)

	// groups
	options := []string{"Core Acc", "Core Dis", "Mid"}
	tradeGroup := tview.NewForm()
	tradeGroup.AddDropDown("Trading Group:", options, 0, func(option string, optionIndex int) {
		switch option {
		case "Core Acc":
			folder = "\\acc"
		case "Core Dis":
			folder = "\\dis"
		case "Mid":
			folder = "\\mid"
		}
		buttonText = fmt.Sprintf("Press To Start %s Backtests", option)
		buttonSort.SetLabel(buttonText)
	})
	// increments
	increments := []string{"  1  ", "  2  ", "  3  ", "  4  ", "  5  "}
	increment := tview.NewForm()
	increment.AddDropDown("Increase RSI in Increments of:", increments, 0, func(option1 string, optionIndex int) {
		// Trim spaces
		option1 = strings.TrimSpace(option1)

		rsi_increment, err = strconv.Atoi(option1)
		if err != nil {
			panic(err)
		}
	})
	// increments
	stops := []string{"Yes", "No"}
	stop := tview.NewForm()
	stop.AddDropDown("Include Stop Loss:", stops, 0, func(option2 string, optionIndex int) {
		switch option2 {
		case "Yes":
			stop_loss = true
		case "No":
			stop_loss = false
		}
	})

	// flex Checklists
	flexChecks := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false). // Empty item for left padding
		AddItem(tradeGroup, 0, 2, false).     // tradeGroup with weight 2
		AddItem(increment, 0, 2, false).      // increment with weight 2
		AddItem(stop, 0, 2, false).           // stop with weight 2
		AddItem(tview.NewBox(), 0, 1, false)  // Empty item for right padding
	flexChecks.SetBorder(true).SetTitle("Choose Trading Group| Choose Increment Increase in RSI| Choose to Incluse a Stop Loss")

	// rsi forms
	low := tview.NewForm()
	low.AddInputField("Entry RSI Minumum", "", 20, nil, func(text string) {
		rsi_low, err = strconv.Atoi(text)
		if err != nil {
			rsi_low = -1
		}

	})
	high := tview.NewForm()
	high.AddInputField("Entry RSI Maximum", "", 20, nil, func(text string) {
		rsi_high, err = strconv.Atoi(text)
		if err != nil {
			rsi_low = -1
		}

	})
	low_enter := tview.NewForm()
	low_enter.AddInputField("Exit RSI Minimum", "", 20, nil, func(text string) {
		rsi_exit_low, err = strconv.Atoi(text)
		if err != nil {
			rsi_low = -1
		}

	})
	high_enter := tview.NewForm()
	high_enter.AddInputField("Exit RSI Maximum", "", 20, nil, func(text string) {
		rsi_exit_high, err = strconv.Atoi(text)
		if err != nil {
			rsi_low = -1
		}

	})
	flexForms := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(low, 0, 2, false).       // Empty item for left padding
		AddItem(high, 0, 2, false).      // tradeGroup with weight 2
		AddItem(low_enter, 0, 2, false). // increment with weight 2
		AddItem(high_enter, 0, 2, false) // Empty item for right padding
	flexForms.SetBorder(true).SetTitle("Enter RSI Parameters")

	buttonSort.SetSelectedFunc(func() {
		message.Clear()

		list, err := os.ReadDir(fmt.Sprintf(dir + folder + "\\data"))
		if err != nil {
			log.Fatal(err, list)
		}

		if (rsi_low < 0) || (rsi_high < 0) || (rsi_exit_high < 0) || (rsi_exit_low < 0) {
			errorChannel <- "Err: Not a Number"
		}

		startBacktests(fmt.Sprintf(dir+folder), list, stop_loss, rsi_low, rsi_high, rsi_exit_low, rsi_exit_high, rsi_increment)

	})
	buttonSort.SetStyle(tcell.StyleDefault.Background(tcell.ColorTurquoise))      //.Foreground(tcell.ColorWhite))
	buttonSort.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGreen)) //.Foreground(tcell.ColorBlack))
	buttonSort.SetLabelColor(tcell.ColorBlack)
	buttonSort.SetLabelColorActivated(tcell.ColorBlack)
	buttonExit := tview.NewButton("Exit")
	buttonExit.SetSelectedFunc(func() {
		exit = true
		app.Stop()
	})
	buttonExit.SetStyle(tcell.StyleDefault.Background(tcell.ColorRed.TrueColor()))              //.Foreground(tcell.ColorWhite))
	buttonExit.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed.TrueColor())) //.Foreground(tcell.ColorBlack))
	buttonExit.SetLabelColor(tcell.ColorBlack.TrueColor())
	buttonExit.SetLabelColorActivated(tcell.ColorBlack.TrueColor())

	buttonList := tview.NewButton("View Upcoming Entrys")
	buttonList.SetSelectedFunc(func() {
		exit = false
		// Set the environment variable
		err = os.Setenv("PERIOD", "1")
		if err != nil {
			log.Fatalf("Failed to set env var: %v", err)
		}
		err = os.Setenv("PYTEMP", "1")
		if err != nil {
			log.Fatalf("Failed to set env var: %v", err)
		}
		app.Stop()
	})
	buttonList.SetStyle(tcell.StyleDefault.Background(tcell.ColorYellow.TrueColor()))                        //.Foreground(tcell.ColorWhite))
	buttonList.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorLightGoldenrodYellow.TrueColor())) //.Foreground(tcell.ColorBlack))
	buttonList.SetLabelColor(tcell.ColorBlack.TrueColor())
	buttonList.SetLabelColorActivated(tcell.ColorBlack.TrueColor())

	flexButtonTop := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
	flexButtonMiddle := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonSort, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonList, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonExit, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)
	flexButtonBottom := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
	flextButton := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flexButtonTop, 0, 2, false).
		AddItem(flexButtonMiddle, 0, 2, false).
		AddItem(flexButtonBottom, 0, 2, false)
	flextButton.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flexChecks, 0, 1, false).
		AddItem(flexForms, 0, 1, false).
		AddItem(flextButton, 0, 1, false).
		AddItem(flextText, 0, 3, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	if !exit {
		list, err := os.ReadDir(fmt.Sprintf(dir + "\\tmp"))
		if err != nil {
			log.Fatal(err, list)
		}

		// waiting := tview.NewApplication()
		// go Downloader(waiting, dir)
		// waiter := tview.NewModal().
		// 	SetText("Downloading Temp Files: Skipping now may cause issues").
		// 	AddButtons([]string{"Skip"})
		// waiter.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		// 	if buttonLabel == "Skip" {
		// 		waiting.Stop()
		// 	}
		// })
		// if err := waiting.SetRoot(waiter, false).EnableMouse(true).SetFocus(modal).Run(); err != nil {
		// 	panic(err)
		// }

		formatEntrys(dir, list)
	}

}

func recieve(message, erros *tview.List) {
	for {
		select {
		case msg, ok := <-channel:
			if !ok {
				channel = nil // Disable this case when the channel is closed
			} else {
				message.AddItem(msg, "", 0, nil)
			}
		case err, ok := <-errorChannel:
			if !ok {
				errorChannel = nil // Disable this case when the channel is closed
			} else {
				erros.AddItem(err, "", 0, nil)
			}
		}

		// Exit the loop when both channels are closed
		if channel == nil && errorChannel == nil {
			break
		}
	}
}

// split the results into message and error message

func Downloader(wait *tview.Application, dir string) {
	pyRun := filepath.Join(dir, "python.exe")

	// Run the Python script
	cmd := exec.Command(pyRun)

	// Run the command and ignore the output
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run python.exe: %v", err)
	}

	wait.Stop()

}
