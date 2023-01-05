package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	krakenAPI              = "https://api.kraken.com/"
	krakenSystemStatusURL  = krakenAPI + "0/public/SystemStatus"
	krakenTradablePairsURL = krakenAPI + "0/public/AssetPairs"
	krakenAssetsTickerURL  = krakenAPI + "0/public/Ticker?pair="
)

type KrakenSystemStatusResponse struct {
	Error  []string `json:"error"`
	Result struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	} `json:"result"`
}

type SystemStatus struct {
	Status    string
	Timestamp string
}

type KrakenTradablePairsResponse struct {
	Error  []string `json:"error"`
	Result map[string]struct {
		Altname string `json:"altname"`
		Wsname  string `json:"wsname"`
		Base    string `json:"base"`
		Quote   string `json:"quote"`
	} `json:"result"`
}

type AssetPair struct {
	Altname string
	Wsname  string
	Base    string
	Quote   string
}

type KrakenTickerResponse struct {
	Error  []string `json:"error"`
	Result map[string]struct {
		C []string `json:"c"`
		V []string `json:"v"`
		T []int    `json:"t"`
		L []string `json:"l"`
		H []string `json:"h"`
	} `json:"result"`
}

type Ticker struct {
	Price      float64
	Volume     float64
	NbOfTrades int
	LowPrice   float64
	HighPrice  float64
}

func getSystemStatus() (SystemStatus, error) {
	var statusServer KrakenSystemStatusResponse

	response, err := http.Get(krakenSystemStatusURL)
	if err != nil {
		fmt.Println(err)
		return SystemStatus{}, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return SystemStatus{}, err
	}

	err = json.Unmarshal(body, &statusServer)
	if err != nil {
		fmt.Println(err)
		return SystemStatus{}, err
	}

	timestampFormatted, _ := time.Parse(time.RFC3339, statusServer.Result.Timestamp)

	return SystemStatus{
		Status:    statusServer.Result.Status,
		Timestamp: timestampFormatted.Format("2006-01-02 15:04:05"),
	}, nil
}

func getAssets() []AssetPair {
	var tradablePairs KrakenTradablePairsResponse

	response, err := http.Get(krakenTradablePairsURL)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = json.Unmarshal(body, &tradablePairs)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	myAssets := []AssetPair{}

	for _, asset := range tradablePairs.Result {

		myAssets = append(myAssets, AssetPair{
			Altname: asset.Altname,
			Wsname:  asset.Wsname,
			Base:    asset.Base,
			Quote:   asset.Quote,
		})

	}

	return myAssets
}

func getTickerInfos(asset string) Ticker {
	var tickerInfos KrakenTickerResponse

	response, err := http.Get(krakenAssetsTickerURL + asset)
	if err != nil {
		fmt.Println(err)
		return Ticker{}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return Ticker{}
	}

	err = json.Unmarshal(body, &tickerInfos)
	if err != nil {
		fmt.Println(err)
		return Ticker{}
	}

	ticker := tickerInfos.Result[asset]

	if len(ticker.C) == 0 {
		return Ticker{}
	}

	price, _ := strconv.ParseFloat(ticker.C[0], 64)
	volume, _ := strconv.ParseFloat(ticker.V[1], 64)
	nbOfTrades := ticker.T[1]
	lowPrice, _ := strconv.ParseFloat(ticker.L[1], 64)
	highPrice, _ := strconv.ParseFloat(ticker.H[1], 64)

	return Ticker{
		Price:      price,
		Volume:     volume,
		NbOfTrades: nbOfTrades,
		LowPrice:   lowPrice,
		HighPrice:  highPrice,
	}
}

func getAssetsTickers(myAssets []AssetPair) {
	for _, asset := range myAssets {
		getTickerInfos(asset.Altname)
	}
}
