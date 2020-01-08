package swap

import (
	"io/ioutil"
	"os"
)

func GetSwap() ([]byte, error) {
	f, err := os.Open("./swap/swap.json")
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
