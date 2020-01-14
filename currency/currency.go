package swap

import (
	"io/ioutil"
	"os"
)

func GetCurrency() ([]byte, error) {
	f, err := os.Open("./currency/currency.json")
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
