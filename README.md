# go-Backtester
go_tradingBacktester


Workflow

main -
sets up configurables to be passed to backtester. Will be where i create the tview configuration. Tview layout will be lifts where i can turn on true of false configurations, set what type i will be doing it oncore,dev,mid etc. If backtests are iterations have a box with text displaying the backtests completed, have a box displaying the current and final best backtest. Have a button to start the backtest function which will start by calling the read function. Thie folowing functions will be nested in a for loop depending on how many iterations they are
Remember a very specific backtest does neceraly make it best because it was so niche it was the best accidnetally and noth through smart logic, so use ball park numbers like 30,32,34 rsi and not 32.23,32.34. Nah d0 25, 30,35,40

read function - reads the csv file and creats a slice of structs. 
    - once read the read function will be passed back to main with slice of structs

main  - slice of structs retrieved from read function. Add rsi and ema's functions will be made which uses the slice of structs and appends the values. More functions can be made for adding more options like maybe bollinger bands. Returns slice of structs with added values

main - starts backtest function which uses the slice of structs. Will iterate throught the slice until it finds an entry. Ane ntry being when all the functions are true, functions are the configurables so if rsi is less than number

Split folders into accs, dis, mid

tbh python good enough for now. Will do this so i can add better tview stuff and graphs and stuff, maybe graphs as it iterates, just want vetter logic

GRAPHSSS

# notes use yahoo finance api to load files into their folder with the updated data

# make ticker list for all my entrys, with different colours for different things. Use drop down table for selection

# implement upcoming trades logic aswell for my ticker list by downloading stocks in the last year, and doing the maths on it. MAke use a tmp folder (use a year for short download)

# implement the trading view position sizing calculator logic

# builds

pyinstaller --onefile python.py

go build -o tradeHelper.exe

rsrc -ico design.ico -o appicon.syso
go build -o finance.exe



# Known bugs
The channles get messed up if enter backtest then exit







