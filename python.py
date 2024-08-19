import yfinance as yf
import datetime
import os

# Define parameters for the query
ticker = 'AAPL'  # Stock ticker symbol

# Get today's date
end_date = datetime.datetime.now().strftime('%Y-%m-%d')

# Define a very early start date to get maximum historical data
# You can use a date several years in the past; here we use 1900-01-01
start_date = '1900-01-01'

tickerAcc = ['AAPL', 'NKE', 'MCD', 'JNJ','ADSK','MSFT','USDV.L','PGR','ASML','AT1.L','MPWR','O','LOW','ADP','DHR','TMO','PEP','MAIN']
tickerDis = ['POOL','KNOS.L','PEN','VRSK','LITE','MA','AVGO','BDX','PAYC','NEE','V','NXST','BAH']
tickerMid = ['CRM','SHOP','LPSN','DAY','HCA','GO','QLYS','TEAM','WEN','DTE','GRMN', 'BRK-B','MTH','COST','BSX','NBIX','CUS','AMD','NVDA','SRPT','ZS','ENTG','PINS']

print("Downloading Accs")
for ticker in tickerAcc:
    data = yf.download(ticker, start=start_date, end=end_date)
    output_dir = 'acc/data/'
    os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
    csv_path = os.path.join(output_dir, f'{ticker}.csv')
    # Export data to a CSV file
    data.to_csv(csv_path)

print("Downloading Dis")
for ticker in tickerDis:
    data = yf.download(ticker, start=start_date, end=end_date)
    output_dir = 'dis/data/'
    os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
    csv_path = os.path.join(output_dir, f'{ticker}.csv')
    # Export data to a CSV file
    data.to_csv(csv_path)

print("Downloading Mids")
for ticker in tickerMid:
    data = yf.download(ticker, start=start_date, end=end_date)
    output_dir = 'mid/data/'
    os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
    csv_path = os.path.join(output_dir, f'{ticker}.csv')
    # Export data to a CSV file
    data.to_csv(csv_path)
    

