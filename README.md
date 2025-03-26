# Stock Guessing Game
The Stock Game is a simple game where the user has to guess the price of a stock for a given date. The user will see the stock price for a given date and will have to draw the price on a canvas. The user will get points based on how close the user's guess is to the actual price.

## Dev Installation

### Database

You can install and create the database with the following commands:

```sh
curl https://install.duckdb.org | sh
duckdb ./data/db/stockgame.duckdb
```

### Data

Download the Kaggle Dataset from the Stock Market Dataset: https://www.kaggle.com/datasets/jacksoncrow/stock-market-dataset, unzip and place the csv files in the data folder.

```sh
unzip stock-market-dataset.zip
mv stock-market-dataset/* data/raw
go run cmd/data-loader/main.go
```

The script will insert the data into the SQL Lite database (about 2 minutes)

```sh
go run cmd/data-loader/main.go
```

Confirming the data is loaded:

```sh
duckdb data/db/stockgame.duckdb
select count(*) from stocks;
┌─────────────────┐
│  count_star()   │
│      int64      │
├─────────────────┤
│    24186113     │
│ (24.19 million) │
└─────────────────┘
```

# Todo

## Backlog Top Priorities

- [ ] Remove hardcodeing URL from StockCanvas.tsx to use a environment variable
- [ ] Obfuscate the Stock to avoid people cheating (remove stock name and date on the initial load)
- [ ] Make the canvas more responsive (mobile friendly?)
- [ ] Make the canvas draw with touch (mobile support?)

## Backlog Lower Priorities
- [ ] Animate the score into the canvas making it more "game" like
- [ ] Make sure the user sent 10 days of prices on the submission
- [ ] Avoid rendering the Canvas HTML over and over again. Should have it once and then just update the canvas
- [ ] Load the stock information into the database (name of the company, see symbols_valid_meta.csv)
- [ ] Create user tables (user, user scores)
- [ ] Add a delay between submission to avoid people cheating
- [ ] Add a leaderboard
- [ ] Add a user registration
- [ ] Add a user login
- [ ] Add a user logout

## Done

- [x] Create a database
- [x] Read CSV and load prices into the database
- [x] Create a simple API that returns the stock price for a given date
- [x] Create a Makefile
- [x] Draw on a Canvas the price of a stock
- [x] Allow the user to draw on the canvas (only after the stock price)
- [x] Clarify each day for the user area (vertical lines)
- [x] Create an API endoint that take the user inputs and result a score
- [x] Show the solution that diff the user input and the stock price
- [x] Determine a logic to assign point (inside low/high gives X points, outside gives Y points)
- [x] Send a score that is not a single number but the details of the score
- [x] Display the volume on the canvas