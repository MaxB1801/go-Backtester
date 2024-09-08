package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	RSI_VALUE      int = 35
	RSI_VALUE_EXIT int = 65
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

func isItExit(data []dayData, group string) bool {
	n := len(data) - 1
	exit := false

	if data[n].close < data[n].ema && group != "acc" {
		exit = true
		return exit
	} else {
		for data[n].close > data[n].ema {
			if data[n].rsi > float64(RSI_VALUE_EXIT) {
				exit = true
				return exit
			}
			n -= 1
		}
	}

	return exit

}

// return number for list
func GetTickerValues(isEntry, isExit bool, dir, ticker string, rawData [][]string) string {

	for _, row := range rawData {

		if row[0] == ticker {
			switch row[1] {
			case "3":
				if isExit {
					return "3"
				} else {
					return "1"
				}
			case "2":
				if isEntry {
					return "2"
				} else {
					return "0"
				}
			case "1":
				if isExit {
					return "3"
				} else {
					return "1"
				}
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

func addEntryExit(dir, group, ticker, number string, data [][]string) {

	folder := fmt.Sprintf(dir + "\\" + group + "\\tickers.csv")
	for n, row := range data {
		if row[0] == ticker {
			data[n][1] = number
		}
	}

	err := writeToCSV(data, folder)
	if err != nil {
		fmt.Println("")
	}

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
	var isExit bool

	var tmpData [][]string

	for _, folder := range list {
		if folder.Name() == "chartFile" {
			continue
		}
		tickers, err := os.ReadDir(fmt.Sprintf(dir + "\\tmp" + "\\" + folder.Name()))
		if err != nil {
			log.Fatal(err, list)
		}

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

		for _, ticker := range tickers {

			// newDir = fmt.Sprintf(dir + "\\" + folder.Name() + "\\tickers.csv")

			// fileOpen, err := os.Open(newDir)
			// if err != nil {

			// 	errorChannel <- fmt.Sprintf("ERROR OPENING FILE: %s", err)
			// }

			// defer fileOpen.Close()

			// // Create a new CSV reader
			// reader := csv.NewReader(fileOpen)

			// // Read all rows from the CSV
			// rawData, err := reader.ReadAll()
			// if err != nil {
			// 	errorChannel <- fmt.Sprintf("ERROR READING CSV: %s", err)
			// }

			data = getData(fmt.Sprintf(dir+"\\tmp"+"\\"+folder.Name()), ticker.Name(), false)

			addMaths(data)

			isEntry = isItEntry(data)

			isExit = isItExit(data, folder.Name())

			number = GetTickerValues(isEntry, isExit, dir, strings.TrimSuffix(ticker.Name(), ".csv"), rawData)

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

// func test(groups2, ticker2 string, list2 [][]string, app *tview.Application) {
// 	app.Stop()
// 	fmt.Println(groups2, ticker2, list2)
// }

func listItem(list [][]string, app *tview.List) ([]string, float64) {
	var new []string
	var count float64
	count = 0
	for _, data := range list {
		new = append(new, data[0])
		switch data[1] {
		case "0":
			app.AddItem(fmt.Sprintf("[white]%s", data[0]), "", 0, nil)
		case "1":
			app.AddItem(fmt.Sprintf("[green]%s", data[0]), "", 0, nil)
		case "2":
			count += 1
			app.AddItem(fmt.Sprintf("[#FFA500]%s", data[0]), "", 0, nil) // .SetMainTextColor(tcell.ColorOrange.TrueColor())
		case "3":
			app.AddItem(fmt.Sprintf("[green]*%s*", data[0]), "", 0, nil)
		}
	}
	return new, count
}

func createGUI(dir string, acc_List, dis_List, mid_List [][]string) {
	var acc []string
	var dis []string
	var mid []string
	var list [][]string
	var list2 [][]string
	var freeFunds float64
	var accR float64
	var disR float64
	var midR float64
	var accCount float64
	var disCount float64
	var midCount float64
	var floatError error
	var groups string
	var groups2 string
	var ticker string
	var ticker2 string
	var stop = false

	accR = 2
	disR = 1
	midR = 1

	for {
		app := tview.NewApplication()

		accInv := tview.NewTextView()
		accInv.SetLabel("Invest £0.00")
		accList := tview.NewList().SetSelectedBackgroundColor(tcell.ColorBlack)
		accRatio := tview.NewForm()
		accRatio.AddInputField("Ratio", "2", 10, nil, func(text string) {
			if text == "" {
				text = "0"
			}
			text = strings.TrimSpace(text)
			text = removeNonNumericChars(text)
			accR, floatError = strconv.ParseFloat(text, 64)
			if floatError != nil {
				freeFunds = 0
			}

		})
		accFlexForms := tview.NewFlex()
		accFlexForms.SetDirection(tview.FlexRow).AddItem(accRatio, 0, 1, false).AddItem(accInv, 0, 1, false)
		flexAcc := tview.NewFlex()
		flexAcc.SetDirection(tview.FlexColumn).AddItem(accList, 0, 1, false).AddItem(accFlexForms, 0, 1, false).SetBorder(true).SetTitle("Core Acc").SetBorderColor(tcell.ColorDarkCyan.TrueColor())

		disInv := tview.NewTextView()
		disInv.SetLabel("Invest £0.00")
		disList := tview.NewList().SetSelectedBackgroundColor(tcell.ColorBlack)
		disRatio := tview.NewForm()
		disRatio.AddInputField("Ratio", "1.5", 10, nil, func(text string) {
			if text == "" {
				text = "0"
			}
			text = strings.TrimSpace(text)
			text = removeNonNumericChars(text)
			disR, floatError = strconv.ParseFloat(text, 64)
			if floatError != nil {
				freeFunds = 0
			}

		})
		disFlexForms := tview.NewFlex()
		disFlexForms.SetDirection(tview.FlexRow).AddItem(disRatio, 0, 1, false).AddItem(disInv, 0, 1, false)
		flexDis := tview.NewFlex()
		flexDis.SetDirection(tview.FlexColumn).AddItem(disList, 0, 1, false).AddItem(disFlexForms, 0, 1, true).SetBorder(true).SetTitle("Core Dis").SetBorderColor(tcell.ColorYellow.TrueColor())

		midInv := tview.NewTextView()
		midInv.SetLabel("Invest £0.00")
		midList := tview.NewList().SetSelectedBackgroundColor(tcell.ColorBlack)
		midRatio := tview.NewForm()
		midRatio.AddInputField("Ratio", "1", 10, nil, func(text string) {
			if text == "" {
				text = "0"
			}
			text = strings.TrimSpace(text)
			text = removeNonNumericChars(text)
			midR, floatError = strconv.ParseFloat(text, 64)
			if floatError != nil {
				freeFunds = 0
			}

		})
		midFlexForms := tview.NewFlex()
		midFlexForms.SetDirection(tview.FlexRow).AddItem(midRatio, 0, 1, false).AddItem(midInv, 0, 1, false)
		flexMid := tview.NewFlex()
		flexMid.SetDirection(tview.FlexColumn).AddItem(midList, 0, 1, false).AddItem(midFlexForms, 0, 1, false).SetBorder(true).SetTitle("Mid").SetBorderColor(tcell.ColorLightCoral.TrueColor())

		acc, accCount = listItem(acc_List, accList)
		dis, disCount = listItem(dis_List, disList)
		mid, midCount = listItem(mid_List, midList)

		// entry group
		percentWin := tview.NewTextView()
		percentWin.SetLabel("0% Avg Return")
		winRates := tview.NewTextView()
		winRates.SetLabel("0% Win Rate")
		entryDropDown := tview.NewForm()
		group := []string{"Core Acc", "Core Dis", "Mid"}
		combDropDown := tview.NewForm()
		combDropDown.AddDropDown("Group: ", group, 0, func(optionGroup string, optionIndex int) {
			entryDropDown.Clear(true)
			switch optionGroup {
			case "Core Acc":
				groups = "acc"
				entryDropDown.AddDropDown("Add Entry: ", acc, 0, func(option string, optionIndex int) {
					list = acc_List
					ticker = option
					winRate(dir, groups, ticker, percentWin, winRates)
				})
			case "Core Dis":
				groups = "dis"
				entryDropDown.AddDropDown("Add Entry: ", dis, 0, func(option string, optionIndex int) {
					list = dis_List
					ticker = option
					winRate(dir, groups, ticker, percentWin, winRates)
				})
			case "Mid":
				groups = "mid"
				entryDropDown.AddDropDown("Add Entry: ", mid, 0, func(option string, optionIndex int) {
					list = mid_List
					ticker = option
					winRate(dir, groups, ticker, percentWin, winRates)
				})
			}
		})
		combDropDown.AddButton("Save", func() {
			addEntryExit(dir, groups, ticker, "1", list)
			for n, row := range list {
				if row[0] == ticker {
					list[n][1] = "1"
				}
			}
			switch groups {
			case "acc":
				accList.Clear()
				acc, accCount = listItem(acc_List, accList)
			case "dis":
				disList.Clear()
				dis, disCount = listItem(dis_List, disList)
			case "mid":
				midList.Clear()
				mid, midCount = listItem(mid_List, midList)
			}

		})

		flexCombs := tview.NewFlex()
		flexCombs.SetBorder(true).SetTitle("Add an Entry")
		flexCombs.SetDirection(tview.FlexColumn).AddItem(combDropDown, 0, 1, true).AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(entryDropDown, 0, 2, true).AddItem(percentWin, 0, 1, false).AddItem(winRates, 0, 1, false), 0, 1, false)

		// entry group
		entryDropDownExit := tview.NewForm()
		combDropDownExit := tview.NewForm()
		combDropDownExit.AddDropDown("Group:", group, 0, func(optionGroup1 string, optionIndex int) {
			entryDropDownExit.Clear(true)
			switch optionGroup1 {
			case "Core Acc":
				groups2 = "acc"
				entryDropDownExit.AddDropDown("Remove Entry: ", acc, 0, func(option2 string, optionIndex int) {
					list2 = acc_List
					ticker2 = option2
				})
			case "Core Dis":
				groups2 = "dis"
				entryDropDownExit.AddDropDown("Remove Entry: ", dis, 0, func(option2 string, optionIndex int) {
					list2 = dis_List
					ticker2 = option2
				})
			case "Mid":
				groups2 = "mid"
				entryDropDownExit.AddDropDown("Remove Entry: ", mid, 0, func(option2 string, optionIndex int) {
					list2 = mid_List
					ticker2 = option2
				})
			}
		})
		combDropDownExit.AddButton("Save", func() {
			addEntryExit(dir, groups2, ticker2, "0", list2)

			for n, row := range list2 {
				if row[0] == ticker2 {
					list[n][1] = "0"
				}
			}
			switch groups2 {
			case "acc":
				accList.Clear()
				acc, accCount = listItem(acc_List, accList)
			case "dis":
				disList.Clear()
				dis, disCount = listItem(dis_List, disList)
			case "mid":
				midList.Clear()
				mid, midCount = listItem(mid_List, midList)
			}
		})

		firstForm := tview.NewForm()
		firstForm.AddInputField("Free Funds £", "0", 10, nil, func(text string) {
			if text == "" {
				text = "0"
			}
			text = strings.TrimSpace(text)
			text = removeNonNumericChars(text)
			freeFunds, floatError = strconv.ParseFloat(text, 64)
			if floatError != nil {
				freeFunds = 0
			}

			findInvestAmountsPerStock(freeFunds, accCount, disCount, midCount, accR, disR, midR, accInv, disInv, midInv)
		})

		firstForm.AddButton("Exit", func() {
			stop = false
			app.Stop()
		})
		firstForm.AddButton("Reconfigure", func() {
			stop = true

			app.Stop()
		})

		formFlex := tview.NewFlex()
		formFlex.SetBorder(true)
		formFlex.AddItem(tview.NewBox(), 0, 1, false).AddItem(firstForm, 0, 3, false).AddItem(tview.NewBox(), 0, 1, false)

		flexCombsExit := tview.NewFlex()
		flexCombsExit.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Remove an Entry")
		flexCombsExit.AddItem(combDropDownExit, 0, 1, true).AddItem(entryDropDownExit, 0, 1, true)

		flexTopRow := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(formFlex, 0, 1, false).AddItem(flexCombs, 0, 2, false).AddItem(flexCombsExit, 0, 2, false)

		flexLists := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(flexAcc, 0, 2, false).AddItem(flexDis, 0, 2, false).AddItem(flexMid, 0, 2, false)

		FinalFlex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(flexTopRow, 0, 1, false).AddItem(flexLists, 0, 4, false)

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			switch event.Key() {
			case tcell.KeyEnter:
				findInvestAmountsPerStock(freeFunds, accCount, disCount, midCount, accR, disR, midR, accInv, disInv, midInv)
			default:
				findInvestAmountsPerStock(freeFunds, accCount, disCount, midCount, accR, disR, midR, accInv, disInv, midInv)
				enterEvent := tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
				// Post this event to the application's event queue
				app.QueueEvent(enterEvent)

			}
			return event
		})

		if err := app.SetRoot(FinalFlex, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		if stop {

			waiting := tview.NewApplication()
			go Downloader(waiting, dir)
			waiter := tview.NewModal().
				SetText("Downloading Temp Files: Skipping now may cause issues").
				AddButtons([]string{"Skip"})
			waiter.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Skip" {
					waiting.Stop()
				}
			})
			if err := waiting.SetRoot(waiter, false).EnableMouse(true).Run(); err != nil {
				panic(err)
			}

			list, err := os.ReadDir(fmt.Sprintf(dir + "\\tmp"))
			if err != nil {
				log.Fatal(err, list)
			}

		} else {
			break
		}
	}

}

func winRate(dir, group, ticker string, app *tview.TextView, app2 *tview.TextView) {

	fileOpen, err := os.Open(fmt.Sprintf(dir + "\\" + group + "\\trades\\" + ticker + ".csv"))
	if err != nil {
		app.SetLabel("1% Win Rate")
	}

	defer fileOpen.Close()

	// Create a new CSV reader
	reader := csv.NewReader(fileOpen)

	// Read all rows from the CSV
	rawData, err := reader.ReadAll()
	if err != nil {
		app.SetLabel("2% Win Rate")
	}

	var total float64 = 0
	var tmp float64 = 0
	var count float64

	for _, row := range rawData {
		tmp, err = strconv.ParseFloat(row[5], 64)
		if err != nil {
			app.SetLabel("3% Win Rate")
		}
		if tmp >= 0 {
			count += 1
		}
		total += tmp
	}
	// app2.SetLabel(fmt.Sprintf("%.2f -- %.2f", count, (float64)(len(rawData))))

	app2.SetLabel(fmt.Sprintf("%.2f%% Win Rate", 100*(count/(float64)(len(rawData)))))

	app.SetLabel(fmt.Sprintf("%.2f%% Avg Return", total/(float64)(len(rawData)-1)))

}

func findInvestAmountsPerStock(funds, accCount, disCount, midCount, accR, disR, midR float64, accText, disText, midText *tview.TextView) { // (float64, float64, float64, float64) { //, freeFundChan, accRChan, disRChan, midRChan chan float64, boolChan chan bool, app *tview.Application) {

	x := funds / ((accR * accCount) + (disR * disCount) + (midR * midCount))

	invAcc := accR * x
	invDis := disR * x
	invMid := midR * x

	if accCount == 0 {
		invAcc = 0
	}
	if disCount == 0 {
		invDis = 0
	}
	if midCount == 0 {
		invMid = 0
	}

	accText.Clear()
	disText.Clear()
	midText.Clear()
	accText.SetLabel(fmt.Sprintf("Invest £%.2f", invAcc))
	disText.SetLabel(fmt.Sprintf("Invest £%.2f", invDis))
	midText.SetLabel(fmt.Sprintf("Invest £%.2f", invMid))

}

func removeNonNumericChars(input string) string {

	// if multiple "." return 0
	if strings.Count(input, ".") > 1 {
		return "0"
	}

	// Compile a regular expression to match non-numeric characters
	re := regexp.MustCompile(`[^\d.]`)

	// Replace all non-numeric characters with an empty string
	return re.ReplaceAllString(input, "")
}
