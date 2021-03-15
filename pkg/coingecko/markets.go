package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type MarketsResponse struct {
	ID                                  string    `json:"id"`
	Symbol                              string    `json:"symbol"`
	Name                                string    `json:"name"`
	Image                               string    `json:"image"`
	CurrentPrice                        float64   `json:"current_price"`
	MarketCap                           int64     `json:"market_cap"`
	MarketCapRank                       int64     `json:"market_cap_rank"`
	FullyDilutedValuation               int64     `json:"fully_diluted_valuation"`
	TotalVolume                         float64   `json:"total_volume"`
	High24H                             float64   `json:"high_24h"`
	Low24H                              float64   `json:"low_24h"`
	PriceChange24H                      float64   `json:"price_change_24h"`
	PriceChangePercentage24H            float64   `json:"price_change_percentage_24h"`
	MarketCapChange24H                  float64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H        float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply                   float64   `json:"circulating_supply"`
	TotalSupply                         float64   `json:"total_supply"`
	MaxSupply                           float64   `json:"max_supply"`
	Ath                                 float64   `json:"ath"`
	AthChangePercentage                 float64   `json:"ath_change_percentage"`
	AthDate                             time.Time `json:"ath_date"`
	Atl                                 float64   `json:"atl"`
	AtlChangePercentage                 float64   `json:"atl_change_percentage"`
	AtlDate                             time.Time `json:"atl_date"`
	ROI                                 ROI       `json:"roi"`
	LastUpdated                         time.Time `json:"last_updated"`
	PriceChangePercentage14DInCurrency  float64   `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage1HInCurrency   float64   `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage1YInCurrency   float64   `json:"price_change_percentage_1y_in_currency"`
	PriceChangePercentage200DInCurrency float64   `json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage24HInCurrency  float64   `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage30DInCurrency  float64   `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage7DInCurrency   float64   `json:"price_change_percentage_7d_in_currency"`
}

func FetchMarkets(resultsPerPage, displayPage int) ([]MarketsResponse, error) {
	m := []MarketsResponse{}

	if resultsPerPage > 250 {
		return m, errors.New("Results per page must be 250 or less")
	}

	body, err := makeCall(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=%d&page=%d&sparkline=false&price_change_percentage=%s", resultsPerPage, displayPage, "1h,24h,7d,14d,30d,200d,1y"))
	if err != nil {
		return m, err
	}
	if err = json.Unmarshal(body, &m); err != nil {
		return m, errors.New("Failed to unmarshal market response: " + err.Error())
	}

	return m, nil
}
