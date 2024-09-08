import plotly.graph_objects as go
from plotly.subplots import make_subplots
import pandas as pd
import os

folder = os.getenv('FOLDER', '0')
ticker = os.getenv('TICKER', '0')

# Get the current working directory
dir= os.getcwd()

tmpcsv = "{}\\tmp\\chartFile\\tmpData.csv".format(dir)
tradeFile = "{}\\{}\\trades\\{}.csv".format(dir,folder,ticker)

# Load the data
csvFile = pd.read_csv(tmpcsv)
csvFile['Date'] = pd.to_datetime(csvFile['Date'])  # Ensure dates are in datetime format

name = "{} Candlestick Chart".format(ticker)

# Load trades and create arrows
trades = pd.read_csv(tradeFile)
trades['Entry Date'] = pd.to_datetime(trades['Entry Date'])

# Create subplots with shared x-axis
fig = make_subplots(rows=2, cols=1, shared_xaxes=True, 
                    row_heights=[0.7, 0.3], vertical_spacing=0.05, 
                    subplot_titles=(name, '14-Day RSI'))

# Add candlestick chart to the first subplot
fig.add_trace(go.Candlestick(
    x=csvFile['Date'],
    open=csvFile['Open'],
    high=csvFile['High'],
    low=csvFile['Low'],
    close=csvFile['Close'],
    name='Candlestick'
), row=1, col=1)

# Add 50-day moving average to the first subplot
fig.add_trace(go.Scatter(
    x=csvFile['Date'],
    y=csvFile['Ema'],
    mode='lines',
    name='20-Day MA',
    line=dict(color='blue', width=2, dash='solid')
), row=1, col=1)

# Add RSI to the second subplot
fig.add_trace(go.Scatter(
    x=csvFile['Date'],
    y=csvFile['Rsi'],  # Replace with actual RSI calculation if needed
    mode='lines',
    name='RSI',
    line=dict(color='green', width=2)
), row=2, col=1)

# Update RSI subplot with overbought/oversold lines
fig.add_hline(y=70, line_dash='dash', row=2, col=1, line_color='red')
fig.add_hline(y=30, line_dash='dash', row=2, col=1, line_color='blue')

# Update layout for better interactivity and appearance
fig.update_layout(
    title=name,
    xaxis2_title='Date',
    yaxis_title='Close Price',
    xaxis_rangeslider_visible=False,  # Optionally hide range slider
    height=900,
    dragmode='zoom',  # Enable zooming by dragging

)

for index, trade in trades.iterrows():

    # Define arrow positions
    low_price = trade['Entry Price']
    high_price = trade['Exit Price']
    
   # arrowH_price = high_price
    arrowL_price = low_price - (low_price * 0.025)
    arrowH_price = high_price + (high_price * 0.025)
    
    # Add a small green vertical arrow below the candle at the specified date
    fig.add_shape(
        dict(
            type="line",
            x0=trade['Entry Date'], x1=trade['Entry Date'],  # Position the arrow along the date
            y0=arrowL_price, y1=arrowL_price + (low_price * 0.1),  # Small arrow length
            line=dict(color="green", width=2),
            opacity=0,
            xref="x", yref="y"  # Referencing the main x and y axis (first subplot)
        )
    )

    # Add a small arrowhead to the bottom
    fig.add_annotation(
        dict(
            x=trade['Entry Date'],  # Position on the x-axis (the date)
            y=arrowL_price,  # Position just below the start of the arrow
            ax=trade['Entry Date'],  # X-axis of arrow tail (relative to x)
            ay=arrowL_price - (low_price * 0.05),  # Y-axis of arrow tail (adjust for arrowhead size)
            xref="x", yref="y", axref="x", ayref="y",  # References for the arrow coordinates
            showarrow=True,
            arrowhead=2,  # Arrowhead style
            arrowsize=3,  # Make arrowhead small
            arrowcolor="green",
            text="T{}".format(index+1),  # Text for the trade
            font=dict(color="green", size=14)
        )
    )
    # exits

        # Add a small green vertical arrow below the candle at the specified date
    fig.add_shape(
        dict(
            type="line",
            x0=trade['Exit Date'], x1=trade['Exit Date'],  # Position the arrow along the date
            y0=arrowH_price, y1=arrowH_price - (high_price * 0.1),  # Small arrow length
            line=dict(color="red", width=2),
            opacity=0,
            xref="x", yref="y"  # Referencing the main x and y axis (first subplot)
        )
    )

    # Add a small arrowhead to the bottom
    fig.add_annotation(
        dict(
            x=trade['Exit Date'],  # Position on the x-axis (the date)
            y=arrowH_price,  # Position just below the start of the arrow
            ax=trade['Exit Date'],  # X-axis of arrow tail (relative to x)
            ay=arrowH_price + (high_price * 0.05),  # Y-axis of arrow tail (adjust for arrowhead size)
            xref="x", yref="y", axref="x", ayref="y",  # References for the arrow coordinates
            showarrow=True,
            arrowhead=2,  # Arrowhead style
            arrowsize=3,  # Make arrowhead small
            arrowcolor="red",
            text="",  # Text for the trade
            font=dict(color="red", size=14)
        )
    )

# Show the plot
fig.show()

print("Stop the loading on go")