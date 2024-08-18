package main

import (
	"math"
)

func round2dp(value float64) float64 {
	return math.Round(value*100) / 100
}

func get20(n int, data []dayData) [20]dayData {
	var days [20]dayData
	for i := 0; i < 20; i++ {
		//fmt.Println(i)
		days[19-i] = data[n-i]
	}
	return days
}

func getEMA(data [20]dayData) float64 {
	var total float64 = 0
	for n := 0; n < 20; n++ {
		total += data[n].close
	}

	return round2dp(total / 20)
}

// rma calculates the Relative Moving Average for a slice of floats
func rma(data []float64, length int) []float64 {
	alpha := 1.0 / float64(length)
	rma := make([]float64, len(data))
	rma[0] = data[0] // SMA for the first value

	for i := 1; i < len(data); i++ {
		rma[i] = alpha*data[i] + (1-alpha)*rma[i-1]
	}

	return rma
}

// calculateRSI calculates the RSI for a slice of dayData
func calculateRSI(data []dayData, length int) {
	// Calculate the differences (delta)
	delta := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		delta[i-1] = data[i].close - data[i-1].close
	}

	// Separate the gains and losses
	gains := make([]float64, len(delta))
	losses := make([]float64, len(delta))

	for i := 0; i < len(delta); i++ {
		if delta[i] > 0 {
			gains[i] = delta[i]
		} else {
			losses[i] = -delta[i]
		}
	}

	// Calculate RMA for gains and losses
	avgGain := rma(gains, length)
	avgLoss := rma(losses, length)

	// Calculate RSI
	for i := length; i < len(data); i++ {
		rs := avgGain[i-1] / avgLoss[i-1]
		data[i].rsi = round2dp(100 - (100 / (1 + rs)))
	}
}

func getLow(data []dayData, start, end int) float64 {
	low := data[start].low
	for start < end {
		if data[start].low < low {
			low = data[start].low
		}
		start += 1
	}
	return low
}

func getROI(trader []trades) float64 {
	var totalROI float64
	for _, trade := range trader {
		totalROI += trade.ROI
	}

	return totalROI
}
