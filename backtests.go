package main

import (
	"encoding/csv"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type dayData struct {
	date  string
	low   float64
	close float64
	rsi   float64
	ema   float64
}

type trades struct {
	Entry_Date   string
	Entry_Price  float64
	Exit_Date    string
	Exit_Price   float64
	Trade_Length int
	ROI          float64
}

type results struct {
	stock    string
	rsi      int
	rsi_exit int
	roi      float64
}

// retur data as struct
func getData(dir string, stock fs.DirEntry) []dayData {

	dir = fmt.Sprintf(dir + "\\data")

	file := filepath.Join(dir, stock.Name())

	fileOpen, err := os.Open(file)
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

	var data []dayData
	var day dayData
	var lowV float64
	var closeV float64
	var err1 error
	for _, row := range rawData {
		lowV, err1 = strconv.ParseFloat(row[3], 64)
		if err != nil {
			fmt.Printf("Error converting string to float64: %v\n", err1)
		}

		closeV, err1 = strconv.ParseFloat(row[4], 64)
		if err != nil {
			fmt.Printf("Error converting string to float64: %v\n", err1)
		}

		day = dayData{
			date:  row[0],
			low:   lowV,
			close: closeV,
		}
		data = append(data, day)
	}
	return data
}

func outputFiles(dir string, stock fs.DirEntry, results []trades) {
	fileName := filepath.Join(dir, stock.Name())
	file, err := os.Create(fileName)
	if err != nil {
		errorChannel <- fmt.Sprintf("ERROR Writing to CSV: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Entry Date", "Entry Price", "Exit Date", "Exit Price", "Trade Length in Days", "ROI %"}
	if err := writer.Write(header); err != nil {
		errorChannel <- fmt.Sprintf("error writing header to CSV: %s", err)
	}

	// Write the data rows
	for _, trade := range results {
		record := []string{
			trade.Entry_Date,
			strconv.FormatFloat(trade.Entry_Price, 'f', 2, 64),
			trade.Exit_Date,
			strconv.FormatFloat(trade.Exit_Price, 'f', 2, 64),
			strconv.Itoa(trade.Trade_Length),
			strconv.FormatFloat(trade.ROI, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

}

func resultsFile(dir string, trades []results) {
	dir = fmt.Sprintf(dir + "\\Trade-Parameters.csv")

	file, err := os.Create(dir)
	if err != nil {
		errorChannel <- fmt.Sprintf("Error creating results file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	header := []string{"Stock", "RSI Entry", "RSI Exit", "ROI %"}
	if err := writer.Write(header); err != nil {
		errorChannel <- fmt.Sprintf("error writing header to CSV: %s", err)
	}

	// Write the data rows
	for _, result := range trades {
		record := []string{
			result.stock,
			strconv.Itoa(result.rsi),
			strconv.Itoa(result.rsi_exit),
			strconv.FormatFloat(result.roi, 'f', 2, 64),
		}
		if err := writer.Write(record); err != nil {
			errorChannel <- fmt.Sprintf("error writing record to CSV: %s", err)
		}
	}

}

func addMaths(data []dayData) {

	//calculating ema
	for n, _ := range data {
		if n < 19 {
			continue
		}

		last20 := get20(n, data)

		data[n].ema = getEMA(last20)
	}

	calculateRSI(data, 14)

}

func findEntry(n int, data []dayData, rsi int, stop_loss bool) (int, trades, float64) {

	var trade trades
	var ranges int

	// iterates n until entries can be looked for
	for {
		if n > len(data)-1 {
			return n, trades{}, 0
		}
		if data[n].rsi < float64(rsi) && data[n].close < data[n].ema {
			ranges = n
			break
		}
		n += 1
	}

	for {

		if n > len(data)-1 {
			return n, trades{}, 0
		}

		if data[n].close > data[n].ema {
			trade = trades{
				Entry_Date:   data[n].date,
				Entry_Price:  data[n].close,
				Exit_Date:    "",
				Exit_Price:   0,
				Trade_Length: 0,
				ROI:          0,
			}

			if stop_loss {
				return n, trade, getLow(data, ranges, n)
			} else {
				return n, trade, 0
			}
		}
		n += 1
	}

}

func findExit(n int, data []dayData, rsi int, trade trades, low float64) (int, trades) {
	length := n
	// iterates n until exits can be looked for
	for {
		if n > len(data)-1 {
			return n, trades{}
		}

		if data[n].close < low {
			trade.Exit_Date = data[n].date
			trade.Exit_Price = data[n].close
			trade.Trade_Length = n - length                                              // Calculate the length of the trade
			trade.ROI = (trade.Exit_Price - trade.Entry_Price) / trade.Entry_Price * 100 // Calculate ROI
			return n, trade
		}

		if data[n].rsi > float64(rsi) {
			break
		}
		n += 1

	}

	for {
		if n > len(data)-1 {
			return n, trades{}
		}
		if data[n].close < data[n].ema {
			trade.Exit_Date = data[n].date
			trade.Exit_Price = data[n].close
			trade.Trade_Length = n - length                                              // Calculate the length of the trade
			trade.ROI = (trade.Exit_Price - trade.Entry_Price) / trade.Entry_Price * 100 // Calculate ROI
			return n, trade

		}
		n += 1
	}

}

func startBacktests(dir string, list []fs.DirEntry, stop_loss bool, rsi_low, rsi_high, rsi_exit_low, rsi_exit_high, rsi_increment int) {

	var result results
	var resultSLice []results

	var test []trades
	var data []dayData
	var rsi int
	var rsi_exit int
	var totalROI float64
	for _, stock := range list {

		if rsi_high < rsi_low {
			errorChannel <- "RSI High Lower Than the RSI Low"
			return
		}
		if rsi_exit_high < rsi_exit_low {
			errorChannel <- "RSI Exit High Lower Than the RSI Exit Low"
			return
		}
		if rsi_high == 0 && rsi_exit_high == 0 && rsi_low == 0 && rsi_exit_low == 0 {
			errorChannel <- "No RSI Input Detected"
			return
		}

		data = getData(dir, stock)
		stockName := strings.TrimSuffix(stock.Name(), ".csv") // Adjust extension as needed

		// adds rsi and ema to structs
		addMaths(data)

		test, rsi, rsi_exit, totalROI = backtest(data, stop_loss, rsi_low, rsi_high, rsi_exit_low, rsi_exit_high, rsi_increment)

		// doing output logic
		outputFiles(fmt.Sprintf(dir+"\\trades"), stock, test)

		result = results{
			stock:    stockName,
			rsi:      rsi,
			rsi_exit: rsi_exit,
			roi:      totalROI,
		}

		channel <- fmt.Sprintf("%s: [RSI: %d] [RSI EXIT: %d] [ROI %.2f%%]", stockName, rsi, rsi_exit, totalROI)

		resultSLice = append(resultSLice, result)

	}

	// add reults to results file
	resultsFile(dir, resultSLice)

}

func backtest(data []dayData, stop_loss bool, rsi_low, rsi_high, rsi_exit_low, rsi_exit_high, rsi_increment int) ([]trades, int, int, float64) {
	var test []trades
	var best []trades
	var trade trades
	var low float64
	var returnOnInvestment float64
	var bestRSI int
	var bestExitRSI int
	var totalROI float64 = 0

	var pointer int
	for rsi := rsi_low; rsi <= rsi_high; rsi += rsi_increment {
		for rsi_exit := rsi_exit_low; rsi_exit <= rsi_exit_high; rsi_exit += rsi_increment {
			n := 19
			for {

				// sends n and slice
				n, trade, low = findEntry(n, data, rsi, stop_loss)
				if n > len(data)-1 {
					break
				}

				// used for marking the spot for entering trades whilst in trades
				pointer = n

				n, trade = findExit(n, data, rsi_exit, trade, low)

				if n > len(data)-1 {
					break
				}

				test = append(test, trade)
				n = pointer

			}

			returnOnInvestment = getROI(test)

			// use channels here maybe to display info
			if returnOnInvestment > totalROI {

				best = test
				bestRSI = rsi
				bestExitRSI = rsi_exit
				totalROI = returnOnInvestment
			}
			test = []trades{}
		}

	}

	return best, bestRSI, bestExitRSI, totalROI
}
