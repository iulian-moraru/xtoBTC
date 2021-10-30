package main

import (
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
	"os"
	"sort"
	"strings"
)

func main() {
	apiKey := os.Getenv("APIKEY")
	apiSecret := os.Getenv("APISECRET")
	//
	client := bitfinex.NewClient().Auth(apiKey, apiSecret)
	allSymbols, err := client.Symbols.GetSymbols()

	if err != nil {
		os.Exit(1)
	}

	sortSymbols(allSymbols)

}

func sortSymbols(allSymbols []string) {
	sort.Strings(allSymbols)
	var symbsList []string
	for _, symb := range allSymbols {
		switch {
		case strings.Contains(symb, "f0") || strings.Contains(symb, "test") ||
			strings.Contains(symb, "aaa") || strings.Contains(symb, "eut") ||
			strings.Contains(symb, "eut"):
			continue
		case symb == "btcusd" || symb == "btcust":
			continue
		case strings.Contains(symb, "btc") || strings.Contains(symb, "usd"):
			symbsList = append(symbsList, symb)
		}
	}
	var finalSymbs []string
	c := make(chan []string)
	go checkPair(symbsList, c)
	finalSymbs = <-c
	sort.Strings(finalSymbs)
	fmt.Println(finalSymbs)
}

func checkPair(symbs []string, c chan<- []string) {
	var finalSymb []string
	for _, symb := range symbs {
		if strings.Contains(symb, "btc") {
			finalSymb = append(finalSymb, symb)
			continue
		} else {
			if strings.HasSuffix(symb, "usd") {
				scndSymb := getScndSymb(symb)
				pair := ""
				if strings.Contains(symb, ":") {
					pair = fmt.Sprintf("%s:btc", scndSymb)
				} else {
					pair = fmt.Sprintf("%sbtc", scndSymb)
				}

				for _, symb2 := range symbs {
					if strings.Contains(symb2, scndSymb) {
						if pair == symb2 {
							break
						} else {
							chkSymb := getScndSymb(symb2)
							if chkSymb == scndSymb {
								finalSymb = append(finalSymb, symb2)
								break
							}
						}
					}
				}
			}
		}
	}
	c <- finalSymb
}

func getScndSymb(symb string) string {
	scndSymb := ""
	if strings.Contains(symb, ":") {
		scndSymb = strings.Split(symb, ":")[0]
	} else {
		scndSymb = strings.Split(symb, "usd")[0]
	}
	return scndSymb
}
