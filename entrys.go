package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	RSI_VALUE int = 35
)

func writeToCSV(data [][]string, folder string) error {

	// Open the file with write and truncate mode to clear its contents
	file, err := os.OpenFile(folder, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write the data to the CSV file
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Ensure all writes are flushed
	writer.Flush()

	// Check if there were any errors during the flush
	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func isItEntry(data []dayData) bool {

	n := len(data) - 1
	enter := false
	for data[n].close < data[n].ema {
		if data[n].rsi < float64(RSI_VALUE) {
			enter = true
		}
		n -= 1
	}
	return enter
}

// return number for list
func GetTickerValues(isEntry bool, dir, ticker string, rawData [][]string) string {

	for _, row := range rawData {

		if row[0] == ticker {
			switch row[1] {
			case "2":
				if isEntry {
					return "2"
				} else {
					return "0"
				}
			case "1":
				return "1"
			case "0":
				if isEntry {
					return "2"
				} else {
					return "0"
				}
			}
		}

	}
	return "0"
}

// finds entrys and appends to csv. Formats the data essentially
func formatEntrys(dir string, list []fs.DirEntry) {
	var number string
	var isEntry bool
	var data []dayData
	var newDir string
	var acc_list [][]string
	var dis_list [][]string
	var mid_list [][]string

	var tmpData [][]string

	for _, folder := range list {
		tickers, err := os.ReadDir(fmt.Sprintf(dir + "\\tmp" + "\\" + folder.Name()))
		if err != nil {
			log.Fatal(err, list)
		}
		for _, ticker := range tickers {

			newDir = fmt.Sprintf(dir + "\\" + folder.Name() + "\\tickers.csv")

			fileOpen, err := os.Open(newDir)
			if err != nil {

				errorChannel <- fmt.Sprintf("ERROR OPENING FILE: %s", err)
			}

			defer fileOpen.Close()

			// Create a new CSV reader
			reader := csv.NewReader(fileOpen)

			// Read all rows from the CSV
			rawData, err := reader.ReadAll()
			if err != nil {
				errorChannel <- fmt.Sprintf("ERROR READING CSV: %s", err)
			}

			data = getData(fmt.Sprintf(dir+"\\tmp"+"\\"+folder.Name()), ticker)

			addMaths(data)

			isEntry = isItEntry(data)

			number = GetTickerValues(isEntry, dir, strings.TrimSuffix(ticker.Name(), ".csv"), rawData)

			tmpData = append(tmpData, []string{
				strings.TrimSuffix(ticker.Name(), ".csv"),
				(number),
			})

		}
		// fmt.Println("\n", tmpData, folder.Name())

		if folder.Name() == "acc" {
			acc_list = tmpData
		}
		if folder.Name() == "dis" {
			dis_list = tmpData
		}
		if folder.Name() == "mid" {
			mid_list = tmpData
		}

		err1 := writeToCSV(tmpData, newDir)
		if err1 != nil {
			fmt.Println("Error Creating CSV")
		}

		tmpData = [][]string{}
	}

	// gui needs to read the data at tickers.csv for formatting the colours and therefore that is the only thing it will be passed. It will also change the data
	// incorrect it will be passed three lists. Acc lists, dis list and mid list. Then its easier cause it just needs to edit the list and run the writer code again
	createGUI(dir, acc_list, dis_list, mid_list)
}

func createGUI(dir string, acc_List, dis_List, mid_List [][]string) {
	var acc []string
	var dis []string
	var mid []string
	// var enterAcc string
	// var enterDis string
	// var enterMid string
	// var exitAcc string
	// var exitMid string
	// var exitDis string

	// fmt.Println(acc_List, "\n")
	// fmt.Println(dis_List, "\n")
	// fmt.Println(mid_List, "\n")

	app := tview.NewApplication()

	accList := tview.NewList()

	disList := tview.NewList()
	midList := tview.NewList()
	accList.SetBorder(true).SetTitle("Core Acc").SetBorderColor(tcell.ColorDarkCyan.TrueColor())
	disList.SetBorder(true).SetTitle("Core Dis").SetBorderColor(tcell.ColorYellow.TrueColor())
	midList.SetBorder(true).SetTitle("Mid").SetBorderColor(tcell.ColorLightCoral.TrueColor())

	for _, data := range acc_List {
		acc = append(acc, data[0])
		switch data[1] {
		case "0":
			accList.AddItem(fmt.Sprintf("[white]%s", data[0]), "", 0, nil)
		case "1":
			accList.AddItem(fmt.Sprintf("[green]%s", data[0]), "", 0, nil)
		case "2":
			accList.AddItem(data[0], "", 0, nil).SetMainTextColor(tcell.ColorOrange.TrueColor())
		}
	}

	for _, data := range dis_List {
		dis = append(dis, data[0])
		if data[1] == "0" {
			disList.AddItem(fmt.Sprintf("[white]%s", data[0]), "", 0, nil)
		}
		if data[1] == "1" {
			disList.AddItem(fmt.Sprintf("[green]%s", data[0]), "", 0, nil)
		}
		if data[1] == "2" {
			disList.AddItem(data[0], "", 0, nil).SetMainTextColor(tcell.ColorOrange.TrueColor())
		}
	}

	for _, data := range mid_List {
		mid = append(mid, data[0])
		if data[1] == "0" {
			midList.AddItem(fmt.Sprintf("[white]%s", data[0]), "", 0, nil)
		}
		if data[1] == "1" {
			midList.AddItem(fmt.Sprintf("[green]%s", data[0]), "", 0, nil)
		}
		if data[1] == "2" {
			midList.AddItem(data[0], "", 0, nil).SetMainTextColor(tcell.ColorOrange.TrueColor())
		}
	}

	// entry group
	entryDropDown := tview.NewForm()
	group := []string{"Core Acc", "Core Dis", "Mid"}
	combDropDown := tview.NewForm()
	combDropDown.AddDropDown("Group: ", group, 0, func(optionGroup string, optionIndex int) {
		entryDropDown.Clear(true)
		switch optionGroup {
		case "Core Acc":
			entryDropDown.AddDropDown("Add Entry: ", acc, 0, func(optionAcc string, optionIndex int) {})
		case "Core Dis":
			entryDropDown.AddDropDown("Add Entry: ", dis, 0, func(optionAcc string, optionIndex int) {})
		case "Mid":
			entryDropDown.AddDropDown("Add Entry: ", mid, 0, func(optionAcc string, optionIndex int) {})
		}
	})
	combDropDown.AddButton("Save", nil)

	flexCombs := tview.NewFlex()
	flexCombs.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Add an Entry")
	flexCombs.AddItem(combDropDown, 0, 1, true).AddItem(entryDropDown, 0, 1, true)
	///////////////////////////////////////////

	// entry group
	entryDropDownExit := tview.NewForm()
	combDropDownExit := tview.NewForm()
	combDropDownExit.AddDropDown("Group:", group, 0, func(optionGroup1 string, optionIndex int) {
		entryDropDownExit.Clear(true)
		switch optionGroup1 {
		case "Core Acc":
			entryDropDownExit.AddDropDown("Add Entry: ", acc, 0, func(optionAcc string, optionIndex int) {})
		case "Core Dis":
			entryDropDownExit.AddDropDown("Add Entry: ", dis, 0, func(optionAcc string, optionIndex int) {})
		case "Mid":
			entryDropDownExit.AddDropDown("Add Entry: ", mid, 0, func(optionAcc string, optionIndex int) {})
		}
	})
	combDropDownExit.AddButton("Save", nil)

	flexCombsExit := tview.NewFlex()
	flexCombsExit.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Remove an Entry")
	flexCombsExit.AddItem(combDropDownExit, 0, 1, true).AddItem(entryDropDownExit, 0, 1, true)
	///////////////////
	buttonExit := tview.NewButton("Exit").SetSelectedFunc(func() {
		app.Stop()
	})
	buttonExit.SetBorder(true)
	flexButtonTop := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
	flexButtonMiddle := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonExit, 0, 2, false).
		AddItem(tview.NewBox(), 0, 1, false)
	flexButtonBottom := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(tview.NewBox(), 0, 2, false)
	flextButton := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flexButtonTop, 0, 1, false).
		AddItem(flexButtonMiddle, 0, 2, false).
		AddItem(flexButtonBottom, 0, 1, false)

	////////////
	buttonConfigure := tview.NewButton("Reconfigure").SetSelectedFunc(func() {
		app.Stop()
	})
	buttonConfigure.SetBorder(true)

	flexButtonMiddleC := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(buttonConfigure, 0, 2, false).
		AddItem(tview.NewBox(), 0, 1, false)

	flextButtonC := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(flexButtonTop, 0, 1, false).
		AddItem(flexButtonMiddleC, 0, 2, false).
		AddItem(flexButtonBottom, 0, 1, false)

	///////

	flexTopRow := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(flextButton, 0, 1, false).AddItem(flextButtonC, 0, 1, false).AddItem(flexCombs, 0, 3, false).AddItem(flexCombsExit, 0, 3, false)

	flexLists := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(accList, 0, 2, false).AddItem(disList, 0, 2, false).AddItem(midList, 0, 2, false)

	FinalFlex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(flexTopRow, 0, 1, false).AddItem(flexLists, 0, 4, false)

	if err := app.SetRoot(FinalFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
