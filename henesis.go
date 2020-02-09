package henesis

import (
	"encoding/json"
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
	Network  string
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
	Address             string `json:"address"`
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	Owners              string `json:"owners"`
	TotalSupply         string `json:"totalSupply"`
	TokenCountByAccount uint64 `json:"tokenCountByAccount,omitempty"`
}

func NewHenesis(clientID string) (*Henesis, error) {
	h := &Henesis{
		Network:  "mainnet",
		API:      MainnetAPI,
		ClientID: clientID,
	}
	return h, nil
}

func NewHenesisRinkeby(clientID string) (*Henesis, error) {
	h := &Henesis{
		Network:  "rinkeby",
		API:      RinkebyAPI,
		ClientID: clientID,
	}
	return h, nil
}

func (h Henesis) TrustedNodeRPC() string {
	return "https://tn.henesis.io/ethereum/" + h.Network + "?clientId=" + h.ClientID
}

type errorResponse struct {
	Body *errorBody `json:"error"`
}

func (e errorResponse) Error() error {
	return fmt.Errorf("henesis: status %d %s", e.Body.Status, e.Body.Message)
}

type errorBody struct {
	Message string `json:"message"`
	Status  int    `json:"code"`
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
	// fmt.Println(url)
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

	// If it returns array, not an error
	if string(body[0]) != "[" {
		e := new(errorResponse)
		err = json.Unmarshal(body, e)
		if err != nil {
			return nil, err
		}
		if e.Body != nil {
			return nil, e.Error()
		}
	}

	return body, nil
}
