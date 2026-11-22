package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type XMLCurrency struct {
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  int
	CharCode string
	Value    float64
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to YAML config")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	if config.InputFile == "" || config.OutputFile == "" {
		panic("input-file and output-file are required")
	}

	currencies := readCurrencies(config.InputFile)
	for _, currency := range currencies {
		fmt.Println(currency.NumCode, currency.CharCode, currency.Value)
	}
}

func readCurrencies(path string) []Currency {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(encoding string, reader io.Reader) (io.Reader, error) {
		if strings.EqualFold(encoding, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(reader), nil
		}

		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}

	var rates struct {
		XMLName xml.Name      `xml:"ValCurs"`
		Items   []XMLCurrency `xml:"Valute"`
	}

	if err := decoder.Decode(&rates); err != nil {
		panic(err)
	}

	currencies := make([]Currency, 0, len(rates.Items))
	for _, item := range rates.Items {
		value, err := strconv.ParseFloat(strings.ReplaceAll(item.Value, ",", "."), 64)
		if err != nil {
			panic(err)
		}

		currencies = append(currencies, Currency{
			NumCode:  item.NumCode,
			CharCode: item.CharCode,
			Value:    value,
		})
	}

	sort.SliceStable(currencies, func(first, second int) bool {
		return currencies[first].Value > currencies[second].Value
	})

	return currencies
}
