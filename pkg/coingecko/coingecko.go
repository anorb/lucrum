package coingecko

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func makeCall(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return []byte{}, errors.New("Failed to get json: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errors.New("Failed to read body: " + err.Error())
	}
	return body, nil
}
