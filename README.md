# stockgame


## Dev Installation

### Data

Download the Kaggle Dataset from the Stock Market Dataset: https://www.kaggle.com/datasets/jacksoncrow/stock-market-dataset, unzip and place the csv files in the data folder.

```sh
unzip stock-market-dataset.zip
mv stock-market-dataset/* data/raw
go run cmd/data-loader/main.go
```

The script will insert the data into the SQL Lite database

