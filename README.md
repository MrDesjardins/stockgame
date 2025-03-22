# stockgame


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

## Backlog

- [ ] Load the stock information into the database
- [ ] Create an api that take the user inputs and resutn a score
- [ ] Create user tables (user, user scores)

## Done

- [x] Create a database
- [x] Read CSV and load prices into the database
- [x] Create a simple API that returns the stock price for a given date
- [x] Create a Makefile