package henesis

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	MainnetAPI = "https://eth-mainnet.api.henesis.io"
	RinkebyAPI = "https://eth-rinkeby.api.henesis.io"
)

type Henesis struct {
	API      string
	ClientID string
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

func NewHenesis(clientID string) (*Henesis, error) {
	h := &Henesis{
		API:      MainnetAPI,
		ClientID: clientID,
	}
	return h, nil
}

func NewHenesisRinkeby(clientID string) (*Henesis, error) {
	h := &Henesis{
		API:      RinkebyAPI,
		ClientID: clientID,
	}
	return h, nil
}

func (h Henesis) getPath(path string) ([]byte, error) {
	return h.getURL(h.API + path)
}

func (h Henesis) getURL(url string) ([]byte, error) {
	client := new(http.Client)
	if strings.Contains(url, "?") {
		url = url + "&clientId=" + h.ClientID
	} else {
		url = url + "?clientId=" + h.ClientID
	}

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
