import csv

file_path = "./data/raw/stocks/BEP.csv"

with open(file_path, "r") as file:
    reader = csv.reader(file)
    for line_number, row in enumerate(reader, start=1):
        if len(row) == 7 and (row[1].strip() == "" or row[1].lower() == "null"):
            print(f"Line {line_number} has an empty or invalid 'open' value: {row}")