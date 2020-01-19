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

func NewHenesis(apikey, network string) (*Henesis, error) {
	h := &Henesis{
		API:     API,
		Network: network,
	}
	return h, nil
}

func (h Henesis) get(path string) ([]byte, error) {
	client := new(http.Client)
	req, err := http.NewRequest("GET", h.API+path, nil)
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
