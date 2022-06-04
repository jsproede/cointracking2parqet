package main

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type Trades struct {
	TradeType string `csv:"Typ"`
	Amount    string `csv:"Kauf"`
	CoinType  string `csv:"Cur."`
	Price     string `csv:"Preis (Kauf)"`
	TradeDate string `csv:"Trade Datum"`
}

type ParqetTrades struct {
	DateTime   string `csv:"datetime"`
	Price      string `csv:"price"`
	Shares     string `csv:"shares"`
	Amount     string `csv:"amount"`
	Tax        string `csv:"tax"`
	Fee        string `csv:"fee"`
	Type       string `csv:"type"`
	AssetType  string `csv:"assettype"`
	Identifier string `csv:"identifier"`
	Currency   string `csv:"currency"`
}

const (
	layoutCointracking = "02.01.2006 15:04"
	layoutParqet       = "2006-01-02T15:04:05.000Z"
)

func main() {
	tradesFile, err := os.OpenFile("trades.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer tradesFile.Close()

	trades := []*Trades{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.LazyQuotes = true
		r.Comma = ';'
		return r
	})

	if err := gocsv.UnmarshalFile(tradesFile, &trades); err != nil {
		panic(err)
	}

	parqetTrades := []*ParqetTrades{}
	for _, trade := range trades {
		if trade.TradeType == "Trade" {
			time, _ := time.Parse(layoutCointracking, trade.TradeDate)

			priceString := trade.Price
			eurIndex := strings.Index(priceString, " EUR")
			price := priceString[0:eurIndex]

			parqetTrades = append(parqetTrades, &ParqetTrades{
				DateTime:   time.Format(layoutParqet),
				Price:      price,
				Shares:     trade.Amount,
				Amount:     "1",
				Tax:        "0",
				Fee:        "0",
				Type:       "Buy",
				AssetType:  "Crypto",
				Identifier: trade.CoinType,
				Currency:   "EUR",
			})
		}
	}

	file, err := os.Create("converted.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	csvContent, _ := gocsv.MarshalString(parqetTrades)
	file.WriteString(csvContent)
}
