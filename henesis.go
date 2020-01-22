package henesis

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	API = "http://api.henesis.io/"
)

type Henesis struct {
	API     string
	Network string
}

type Token struct {
	ID       string   `json:"id"`
	Owner    string   `json:"owner"`
	URI      string   `json:"uri"`
	Contract Contract `json:"contract"`
}

type Contract struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owners      string `json:"owners"`
	TotalSupply string `json:"totalSupply"`
}

func NewHenesis(apikey, network string) (*Henesis, error) {
	h := &Henesis{
		API:     API,
		Network: network,
	}
	return h, nil
}

func (h Henesis) getPath(path string) ([]byte, error) {
	return h.getURL(h.API + path)
}

func (h Henesis) getURL(url string) ([]byte, error) {
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
