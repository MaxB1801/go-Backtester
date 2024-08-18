package main

import (
	"fmt"
	"log"
	"os"
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

const (
// rsi_low  int = 24
// rsi_high int = 40

// rsi_exit_low  = 60
// rsi_exit_high = 76

// rsi_increment int = 2

// folder string = "\\acc"

// stop_loss bool = true
)

// find entry function

// find exit function

func main() {

	var err error
	var stop_loss bool
	var rsi_low int
	var rsi_high int
	var rsi_exit_low int
	var rsi_exit_high int
	var folder string
	var rsi_increment int
	var buttonText string

	// Create a channel for messages
	channel := make(chan string)

	errorChannel := make(chan string)

	app := tview.NewApplication()

	// create flex to center
	message := tview.NewList()
	message.SetBorder(true).SetTitle("Results").SetBorderColor(tcell.ColorDarkCyan.TrueColor())
	// error channels texts
	errors := tview.NewList()
	errors.SetBorder(true).SetTitle("Error Messages").SetBorderColor(tcell.ColorRed.TrueColor())
	// making flext
	flextText := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(message, 0, 2, false).
		AddItem(errors, 0, 2, false)

	// define button here
	buttonSort := tview.NewButton("Start Backtests")
	// consume messages
	go recieve(message, errors, channel, errorChannel)

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
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		list, err := os.ReadDir(fmt.Sprintf(dir + folder + "\\data"))
		if err != nil {
			log.Fatal(err, list)
		}

		if rsi_low < 0 || rsi_high < 0 || rsi_exit_high < 0 || rsi_exit_low < 0 {
			errorChannel <- "Err: Not a Number"
		}

		startBacktests(fmt.Sprintf(dir+folder), list, stop_loss, rsi_low, rsi_high, rsi_exit_low, rsi_exit_high, rsi_increment, channel, errorChannel)

	})
	buttonSort.SetStyle(tcell.StyleDefault.Background(tcell.ColorTurquoise))      //.Foreground(tcell.ColorWhite))
	buttonSort.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGreen)) //.Foreground(tcell.ColorBlack))
	buttonSort.SetLabelColor(tcell.ColorBlack)
	buttonSort.SetLabelColorActivated(tcell.ColorBlack)
	flexButtonTop := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
	flexButtonMiddle := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(buttonSort, 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
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
		AddItem(flexChecks, 0, 2, false).
		AddItem(flexForms, 0, 2, false).
		AddItem(flextButton, 0, 2, false).
		AddItem(flextText, 0, 2, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}

func recieve(message, erros *tview.List, channel, errorChannel chan string) {
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
