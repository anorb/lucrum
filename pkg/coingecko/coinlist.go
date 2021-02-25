package coingecko

import (
	"encoding/json"
	"errors"
)

type CoinList []struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func FetchCoinList() (CoinList, error) {
	cl := CoinList{}

	body, err := makeCall("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		return cl, err
	}

	if err = json.Unmarshal(body, &cl); err != nil {
		return cl, errors.New("Failed to unmarshal: " + err.Error())
	}

	return cl, nil
}
