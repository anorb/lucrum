package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type SimpleResponse map[string]map[string]float64

func FetchSimplePrice(coins []string) (SimpleResponse, error) {
	s := SimpleResponse{}

	body, err := makeCall(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", strings.Join(coins[:], ",")))
	if err != nil {
		return s, err
	}

	if err = json.Unmarshal(body, &s); err != nil {
		return s, errors.New("Failed to unmarshal: " + err.Error())
	}

	return s, nil
}
