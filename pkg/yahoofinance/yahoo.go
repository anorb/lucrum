package yahoofinance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Stock struct {
	Ask                               float64 `json:"ask"`
	AskSize                           int     `json:"askSize"`
	AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
	AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
	Bid                               float64 `json:"bid"`
	BidSize                           int     `json:"bidSize"`
	BookValue                         float64 `json:"bookValue"`
	Currency                          string  `json:"currency"`
	EpsForward                        float64 `json:"epsForward"`
	EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
	EsgPopulated                      bool    `json:"esgPopulated"`
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
	Exchange                          string  `json:"exchange"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FinancialCurrency                 string  `json:"financialCurrency"`
	FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
	ForwardPE                         float64 `json:"forwardPE"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
	Language                          string  `json:"language"`
	LongName                          string  `json:"longName"`
	MarketCap                         int64   `json:"marketCap"`
	MarketState                       string  `json:"marketState"`
	Market                            string  `json:"market"`
	MessageBoardID                    string  `json:"messageBoardId"`
	PriceHint                         int     `json:"priceHint"`
	PriceToBook                       float64 `json:"priceToBook"`
	QuoteSourceName                   string  `json:"quoteSourceName"`
	QuoteType                         string  `json:"quoteType"`
	Region                            string  `json:"region"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	RegularMarketTime                 int     `json:"regularMarketTime"`
	RegularMarketVolume               int     `json:"regularMarketVolume"`
	SharesOutstanding                 int     `json:"sharesOutstanding"`
	ShortName                         string  `json:"shortName"`
	SourceInterval                    int     `json:"sourceInterval"`
	Symbol                            string  `json:"symbol"`
	Tradeable                         bool    `json:"tradeable"`
	TrailingPE                        float64 `json:"trailingPE"`
	Triggerable                       bool    `json:"triggerable"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	YtdReturn                         float64 `json:"ytdReturn"`
	PostMarketChange                  float64 `json:"postMarketChange"`
	PostMarketChangePercent           float64 `json:"postMarketChangePercent"`
	PostMarketTime                    int     `json:"postMarketTime"`
	PostMarketPrice                   float64 `json:"postMarketPrice"`
	PreMarketChange                   float64 `json:"preMarketChange"`
	PreMarketChangePercent            float64 `json:"preMarketChangePercent"`
	PreMarketTime                     int     `json:"preMarketTime"`
	PreMarketPrice                    float64 `json:"preMarketPrice"`
	FormattedRegularMarketPrice       string
	FormattedRegularMarketChange      string
	FormattedRegularMarketChangePct   string
	FormattedRegularMarketDayHigh     string
	FormattedRegularMarketDayLow      string
	FormattedRegularMarketDayOpen     string
}

type Query struct {
	Quote struct {
		Result []Stock     `json:"result"`
		Error  interface{} `json:"error"`
	} `json:"quoteResponse"`
}

func FetchQuote(symbols []string) ([]Stock, error) {
	q := Query{}

	resp, err := http.Get(fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s", strings.Join(symbols[:], ",")))
	if err != nil {
		return q.Quote.Result, errors.New("Failed to get json: " + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return q.Quote.Result, errors.New("Failed to read body: " + err.Error())
	}
	json.Unmarshal(body, &q)

	if q.Quote.Error != nil {
		return q.Quote.Result, errors.New(q.Quote.Error.(string))
	}

	for i, s := range q.Quote.Result {
		s.FormattedRegularMarketPrice = formatCash(s.RegularMarketPrice)
		s.FormattedRegularMarketChange = formatCash(s.RegularMarketChange)
		s.FormattedRegularMarketChangePct = formatPercentage(s.RegularMarketChangePercent)
		s.FormattedRegularMarketDayHigh = formatCash(s.RegularMarketDayHigh)
		s.FormattedRegularMarketDayLow = formatCash(s.RegularMarketDayLow)
		s.FormattedRegularMarketDayOpen = formatCash(s.RegularMarketOpen)
		q.Quote.Result[i] = s
	}

	return q.Quote.Result, nil
}

func formatPercentage(p float64) string {
	return fmt.Sprintf("%.2f%%", p)
}

func formatCash(c float64) string {
	if c < 0 {
		return fmt.Sprintf("-$%.2f", c*-1)
	}
	return fmt.Sprintf("$%.2f", c)
}
