package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const CURR_LIST_REF = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.json"
const CURR_LIST_RATE_BY_REF = "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/"

type CurrencyPair struct {
	From string
	To   string
}

func (c CurrencyPair) isEmpty() bool {
	return c.From == "" || c.To == ""
}

type CurrencyRate struct {
	FromTo CurrencyPair
	rate   float64
	date   string
}

func NewDefaultCurrencyRate() *CurrencyRate {
	return &CurrencyRate{FromTo: CurrencyPair{From: "", To: ""}, rate: 1, date: ""}
}

func NewCurrencyRate(from, to string, rate float64, date string) *CurrencyRate {
	return &CurrencyRate{
		FromTo: CurrencyPair{From: from, To: to},
		rate:   rate,
		date:   date,
	}
}

type currRateSvc struct {
	client *http.Client
}

func newCurrRateSvc(timeout int64) currRateSvc {
	return currRateSvc{
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

var (
	ErrInvalidResponse  = errors.New("invalid response")
	ErrParseFailed      = errors.New("parse response body")
	ErrCurrencyNotFound = errors.New("currency not found")
)

type RatesResponse map[string]interface{}

func (self currRateSvc) getAllCurrList() (map[string]string, error) {

	body, err := self.fetchURLBody(CURR_LIST_REF)
	if err != nil {
		return nil, fmt.Errorf("Cannot obtain currencies list: %w ", err)
	}
	return self.parseCurrList(body)
}

func (self currRateSvc) rate(from, to string) (*CurrencyRate, error) {

	var remoteRef string = CURR_LIST_RATE_BY_REF + from + ".json"
	body, err := self.fetchURLBody(remoteRef)
	if err != nil {
		return nil, fmt.Errorf("%w: (%s %v)", ErrCurrencyNotFound, from, err)
	}

	var currMap RatesResponse
	err = json.Unmarshal(body, &currMap)
	if err != nil {
		return nil, fmt.Errorf("parse exchange rate response: %w: %v", ErrParseFailed, err)
	}

	jsFrom, ok := currMap[from].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("obtain exchange rate: %w: %s", ErrCurrencyNotFound, from)
	}

	rateVal, ok := jsFrom[to].(float64)
	if !ok {
		return nil, fmt.Errorf("unexpected rate value for currency: %s", to)
	}

	dateVal, ok := currMap["date"].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected date value for currency: %s/%s", from, to)
	}

	return NewCurrencyRate(from, to, rateVal, dateVal), nil
}

func (currRateSvc) parseCurrList(data []byte) (map[string]string, error) {

	var currMap map[string]string

	err := json.Unmarshal(data, &currMap)
	if err != nil {
		return nil, fmt.Errorf("parse currency list%w: %v", ErrParseFailed, err)
	}
	return currMap, nil
}

func (self currRateSvc) fetchURLBody(ref string) ([]byte, error) {

	resp, err := self.client.Get(ref)
	if err != nil {
		return nil, fmt.Errorf("fetching URL data: %w: %v", ErrInvalidResponse, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching URL data: %w: Unexpected status code %d",
			ErrInvalidResponse, resp.StatusCode)
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}
	return body, nil
}
