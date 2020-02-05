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
	TokenCount  uint64 `json:"tokenCount"`
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

type errorResponse struct {
	Body errorBody `json:"error"`
}

func (e errorResponse) Error() error {
	return fmt.Errorf("henesis: status %d %s", e.Body.Status, e.Body.Message)
}

type errorBody struct {
	Message string
	Status  int
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
	fmt.Println(url)
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
	if string(body)[0:9] == "{\"error\":" {
		err := new(errorResponse)
		if err2 := json.Unmarshal(body, err); err2 != nil {
			return nil, err2
		}
		return nil, err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Backend returns status %d msg: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
