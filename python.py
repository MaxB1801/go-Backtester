import yfinance as yf
from datetime import datetime, timedelta
import os
import csv


# Get today's date
end_date = datetime.now()

# Define a very early start date to get maximum historical data
# You can use a date several years in the past; here we use 1900-01-01
# maybe make this configurable at some point, through go file shares or env variables. Yes env vars
period_str= os.getenv('PERIOD', '0')
period = int(period_str)
tmp_files_str = os.getenv('PYTEMP', '0')
tmp_files = int(tmp_files_str)


if period <= 0:
    start_date = end_date - timedelta(days=365 * 150)
else:
    start_date = end_date - timedelta(days=365 * period)

start_data = start_date.strftime('%Y-%m-%d')
end_date = end_date.strftime('%Y-%m-%d')


tickerAcc = []
file_path = 'acc/tickers.csv'
# Open the CSV file
with open(file_path, mode='r') as file:
    # Create a CSV reader object
    csv_reader = csv.reader(file)

    # Iterate over the rows in the CSV file
    for row in csv_reader:
        # Append each row to the data list
        tickerAcc.append(row[0].strip())

tickerDis = []
file_path = 'dis/tickers.csv'
# Open the CSV file
with open(file_path, mode='r') as file:
    # Create a CSV reader object
    csv_reader = csv.reader(file)

    # Iterate over the rows in the CSV file
    for row in csv_reader:
        # Append each row to the data list
        tickerDis.append(row[0].strip())

tickerMid = []
file_path = 'mid/tickers.csv'
# Open the CSV file
with open(file_path, mode='r') as file:
    # Create a CSV reader object
    csv_reader = csv.reader(file)

    # Iterate over the rows in the CSV file
    for row in csv_reader:
        # Append each row to the data list
        tickerMid.append(row[0].strip())

if tmp_files == 0:
    for ticker in tickerAcc:
        data = yf.download(ticker, start=start_date, end=end_date)
        output_dir = 'acc/data/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)


    for ticker in tickerDis:
        data = yf.download(ticker, start=start_date, end=end_date)
        output_dir = 'dis/data/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)


    for ticker in tickerMid:
        data = yf.download(ticker, start=start_date, end=end_date)
        output_dir = 'mid/data/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)
else:
    for ticker in tickerAcc:
        output_dir = 'tmp/acc/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        data = yf.download(ticker, start=start_date, end=end_date)
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)

    for ticker in tickerDis:
        output_dir = 'tmp/dis/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        data = yf.download(ticker, start=start_date, end=end_date)
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)


    for ticker in tickerMid:
        output_dir = 'tmp/mid/'
        os.makedirs(output_dir, exist_ok=True)  # Create the directory if it doesn't exist
        data = yf.download(ticker, start=start_date, end=end_date)
        csv_path = os.path.join(output_dir, f'{ticker}.csv')
        # Export data to a CSV file
        data.to_csv(csv_path)