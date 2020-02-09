package henesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"context"
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

func (h Henesis) getPath(ctx context.Context, path string) ([]byte, error) {
	return h.getURL(ctx, h.API+path)
}

func (h Henesis) getURL(ctx context.Context, url string) ([]byte, error) {
	client := httpClient()
	if strings.Contains(url, "?") {
		url = url + "&clientId=" + h.ClientID
	} else {
		url = url + "?clientId=" + h.ClientID
	}
	// fmt.Println(url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

func httpClient() *http.Client {
	client := new(http.Client)
	var transport http.RoundTripper = &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		DisableKeepAlives:  false,
		DisableCompression: false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 300 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client.Transport = transport
	return client
}
