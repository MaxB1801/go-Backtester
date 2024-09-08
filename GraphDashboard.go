package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func graphDashboard(dir string) {
	var tmp []fs.DirEntry
	var accList []string
	var disList []string
	var midList []string
	var selected string
	var folder string
	var loader bool = false

	groups := []string{"Core Acc", "Core Dis", "Mid"}

	for {
		// creating lists for of tickers that have backtests
		for _, group := range groups {
			switch group {
			case "Core Acc":
				tmp = tradeList(dir, "acc")
				for _, ticker := range tmp {
					accList = append(accList, strings.TrimSuffix(ticker.Name(), ".csv"))
				}

			case "Core Dis":
				tmp = tradeList(dir, "dis")
				for _, ticker := range tmp {
					disList = append(disList, strings.TrimSuffix(ticker.Name(), ".csv"))
				}

			case "Mid":
				tmp = tradeList(dir, "mid")
				for _, ticker := range tmp {
					midList = append(midList, strings.TrimSuffix(ticker.Name(), ".csv"))
				}
			}

		}

		graph := tview.NewApplication()

		groupForm := tview.NewForm()
		selectForm := tview.NewForm()
		groupForm.AddDropDown("Group: ", groups, 1, func(optionGroup string, optionIndex int) {
			selectForm.Clear(true)
			switch optionGroup {
			case "Core Acc":
				selectForm.AddDropDown("Add Entry: ", accList, 0, func(option string, optionIndex int) {
					folder = "acc"
					selected = option

				})

			case "Core Dis":
				selectForm.AddDropDown("Add Entry: ", disList, 0, func(option string, optionIndex int) {
					folder = "dis"
					selected = option

				})

			case "Mid":
				selectForm.AddDropDown("Add Entry: ", midList, 0, func(option string, optionIndex int) {
					folder = "mid"
					selected = option

				})
			}
		})

		seeChart := tview.NewButton("Load Chart of Trades")
		seeChart.SetSelectedFunc(func() {

			loader = true
			graph.Stop()
			// loadChart(dir, folder, selected, seeChart)

		})
		seeChart.SetBorder(false)
		seeChart.SetStyle(tcell.StyleDefault.Background(tcell.ColorGreen.TrueColor()))
		seeChart.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorGreen.TrueColor()))

		exitChart := tview.NewButton("Exit")
		exitChart.SetSelectedFunc(func() {
			loader = false
			graph.Stop()
		})
		exitChart.SetBorder(false)
		exitChart.SetStyle(tcell.StyleDefault.Background(tcell.ColorRed.TrueColor()))
		exitChart.SetActivatedStyle(tcell.StyleDefault.Background(tcell.ColorRed.TrueColor()))

		flexForm := tview.NewFlex()

		flexForm.SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox(), 0, 2, false).AddItem(groupForm, 0, 2, false).AddItem(selectForm, 0, 2, false).AddItem(tview.NewBox(), 0, 2, false)

		flexButton := tview.NewFlex()

		flexButton.SetDirection(tview.FlexColumn).
			AddItem(tview.NewBox(), 0, 2, false).AddItem(exitChart, 0, 2, false).AddItem(tview.NewBox(), 0, 2, false).AddItem(seeChart, 0, 2, false).AddItem(tview.NewBox(), 0, 2, false)

		flexTotal := tview.NewFlex()
		flexTotal.SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 2, false).AddItem(flexForm, 0, 2, false).AddItem(flexButton, 0, 2, false).AddItem(tview.NewBox(), 0, 2, false)

		if err := graph.SetRoot(flexTotal, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		if loader {
			wait := tview.NewApplication()
			go loadChart(dir, folder, selected, wait)
			/////////
			waiter := tview.NewModal().
				SetText("%% Loading Chart %%")

			if err := wait.SetRoot(waiter, false).EnableMouse(true).Run(); err != nil {
				panic(err)
			}

			continue

		} else if !loader {
			break
		}
	}

}

func tradeList(homeDir, group string) []fs.DirEntry {
	reader, err := os.ReadDir(fmt.Sprintf(homeDir + "\\" + group + "\\trades"))
	if err != nil {
		fmt.Println(err)
		main()
	}

	return reader

}

func loadChart(dir, folder, ticker string, wait *tview.Application) {

	data := getData(fmt.Sprintf(dir+"\\"+folder+"\\data"), fmt.Sprintf("%s.csv", ticker), true)

	data = data[1:]

	addMaths(data)

	err := os.RemoveAll(fmt.Sprintf(dir + "\\tmp\\chartFile"))
	if err != nil {
		noFunc()
	}

	err = os.Mkdir(fmt.Sprintf(dir+"\\tmp\\chartFile"), os.ModePerm)
	if err != nil {
		noFunc()
	}

	mathsFile(data, fmt.Sprintf(dir+"\\tmp\\chartFile\\tmpData.csv"))

	err = os.Setenv("FOLDER", folder)
	if err != nil {
		log.Fatalf("Failed to set env var: %v", err)
	}

	err = os.Setenv("TICKER", ticker)
	if err != nil {
		log.Fatalf("Failed to set env var: %v", err)
	}

	pyRun := filepath.Join(dir, "graph.exe")

	// Run the Python script
	cmd := exec.Command(pyRun)

	// Run the command and ignore the output
	err1 := cmd.Run()
	if err1 != nil {
		log.Fatalf("Failed to run python.exe: %v", err)
	}

	// clear files stored in chartFile

	// save data to folder chartfile

	// set folder and ticker as os.Env vars

	// run adapted python

	wait.Stop()

}

func noFunc() {}

func mathsFile(data []dayData, dir string) {

	file, err := os.Create(dir)
	if err != nil {
		errorChannel <- fmt.Sprintf("ERROR Writing to CSV: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Date", "Open", "High", "Low", "Close", "Ema", "Rsi"}
	if err := writer.Write(header); err != nil {
		errorChannel <- fmt.Sprintf("error writing header to CSV: %s", err)
	}

	// Write the data rows
	for _, day := range data {
		record := []string{
			day.date,
			strconv.FormatFloat(day.open, 'f', 2, 64),
			strconv.FormatFloat(day.high, 'f', 2, 64),
			strconv.FormatFloat(day.low, 'f', 2, 64),
			strconv.FormatFloat(day.close, 'f', 2, 64),
			strconv.FormatFloat(day.ema, 'f', 2, 64),
			strconv.FormatFloat(day.rsi, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

}
