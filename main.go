package Backtester

import "fmt"

// NO GPT, EXCEPT FOR HELP WITH RSI MATHS LOGIC
// aims of this code. Use tview for configurables of how to do back test. For now set up with the constants for back testing
// improve code logic for entrys so entrys are neat little functions that must all become true for entry same with exits
// sort lev, accs,dis and mids all into seperate folders for backtests, will be a tview configurable
// add starting pot function for testing that way, might need to be a seperacode

const (
	rsiEntry   bool = true
	ema20Entry bool = true
	rsiExit    bool = true
	eme20Exit  bool = true
)

func main() {
	fmt.Println("lets start")
}
