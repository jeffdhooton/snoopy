# Snoopy - SEO First Page Rank Checking

CLI tool that takes in either a query & url or a csv and determines whether or not you have a first page ranking for it.

## Running

Clone this project down, make sure you have go installed, and run with `go run main.go`.

## Building

Clone this project down, make sure you have go installed, and run `go build`. The executable will be called `snoopy` and put right in the project folder. You can move it to a location thats in your PATH to run it as `snoopy -file ./input.csv`

## Flags

1. **`-query`** The search query for a single run
1. **`-url`** The domain that you want to check for (example.com)
1. **`-file`** Runs a csv through the checks, example csv structure provided at input.csv
