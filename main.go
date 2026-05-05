package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/integrii/flaggy"
)

var app_name string = "ccnv"
var app_version string = "1.0.0"

type cliArgs struct {
	FromTo        CurrencyPair
	Amount        float64
	PrintCurrList bool
	PrintUsage    bool
	PrintVersion  bool
}

func main() {
	cliArgs := parseCliArgs()
	currRateSvc := newCurrRateSvc(5)
	exitCode := 0

	switch {
	case cliArgs.PrintVersion:
		{
			showVersion()
			os.Exit(exitCode)
		}
	case cliArgs.PrintCurrList:
		{
			list, err := currRateSvc.getAllCurrList()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				exitCode = 1
			} else {
				printCurrList(list)
			}
			os.Exit(exitCode)
		}
	}

	if cliArgs.Amount > 0 && !cliArgs.FromTo.isEmpty() {
		rate, err := currRateSvc.rate(cliArgs.FromTo.From, cliArgs.FromTo.To)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			exitCode = 1
		} else {
			valueFrom := NewCurrency(cliArgs.Amount, rate.FromTo.From, 2)
			valueTo := ConvertCurrency(valueFrom, rate.FromTo.To, rate.rate, 2)
			printRatingResult(rate, &valueFrom, &valueTo)
		}
	} else {
		fmt.Println("Invalid arguments. Use -h for help.")
		exitCode = 1
	}
	os.Exit(exitCode)
}

func parseCliArgs() *cliArgs {
	flaggy.DefaultParser.ShowVersionWithVersionFlag = false
	flaggy.DisableCompletion()

	flaggy.SetName(app_name)
	flaggy.SetVersion(app_version)
	flaggy.SetDescription("Currency converter")

	currList := false
	flaggy.Bool(&currList, "l", "list", "list all available currencies")

	version := false
	flaggy.Bool(&version, "v", "version", "print version")

	var amount float64 = 1
	flaggy.Float64(&amount, "a", "amount", "amount (10, 50, 100, ...)")

	currFrom := ""
	flaggy.String(&currFrom, "f", "from", "currency convert from (usd, eur, ...)")

	currTo := ""
	flaggy.String(&currTo, "t", "to", "currency convert to (usd, eur, ... )")

	flaggy.Parse()

	return &cliArgs{
		FromTo: CurrencyPair{
			From: strings.ToLower(currFrom),
			To:   strings.ToLower(currTo)},
		Amount:        amount,
		PrintCurrList: currList,
		PrintVersion:  version,
	}
}

func showVersion() {
	fmt.Printf("%s version %s\n", app_name, app_version)
}

func printCurrList(m map[string]string) {

	for code, name := range m {
		fmt.Println(code, " : ", name)
	}
}

func printRatingResult(rateInfo *CurrencyRate, valueFrom *Currency, valueTo *Currency) {
	fmt.Printf("[%s] [rate: %f]  %s -> %s \n", rateInfo.date, rateInfo.rate, valueFrom.String(), valueTo.String())
}
