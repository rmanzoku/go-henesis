package henesis

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Token struct {
	ID    *big.Int       `json:"id"`
	Owner common.Address `json:"owner"`
	URI   string         `json:"uri"`
}

type getTokensByAccountAddressInput struct {
	AccountAddress    string
	Page              int
	Size              int
	OrderBy           string
	OrderDirection    string
	ContractAddresses []string
}

func (in getTokensByAccountAddressInput) Path() string {
	q := make(url.Values)
	q.Set("page", strconv.Itoa(in.Page))
	q.Set("size", strconv.Itoa(in.Size))
	q.Set("order_by", in.OrderBy)
	q.Set("order_direction", in.OrderDirection)
	q.Set("contractAddresses", strings.Join(in.ContractAddresses, ","))
	return fmt.Sprintf("/nft/v1/accounts/%s/tokens?", in.AccountAddress) + q.Encode()
}

type getTokensByAccountAddressOutput struct {
	Data       []Datum    `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Datum struct {
	ID       string   `json:"id"`
	Owner    string   `json:"owner"`
	URI      string   `json:"uri"`
	Contract Contract `json:"contract"`
}

type Contract struct {
	Address     string      `json:"address"`
	Name        string      `json:"name"`
	Symbol      string      `json:"symbol"`
	Owners      interface{} `json:"owners"`
	TotalSupply string      `json:"totalSupply"`
}

type Pagination struct {
	TotalCount int64  `json:"totalCount"`
	PrevURL    string `json:"prevUrl"`
	NextURL    string `json:"nextUrl"`
}

func parseDatum(d Datum) *Token {
	id, _ := new(big.Int).SetString(d.ID, 10)
	return &Token{
		ID:    id,
		Owner: common.HexToAddress(d.Owner),
		URI:   d.URI,
	}
}

func (h Henesis) GetTokensByAccountAddress(accountAddress string, contractAddresses []string) (tokens []*Token, err error) {
	in := &getTokensByAccountAddressInput{
		AccountAddress:    accountAddress,
		Page:              0,
		Size:              200,
		OrderBy:           "transfer_block_number",
		OrderDirection:    "desc",
		ContractAddresses: contractAddresses,
	}

	next := h.API + in.Path()
	i := 0
	if next != "" {
		b, err := h.getURL(next)
		if err != nil {
			return nil, err
		}
		out := new(getTokensByAccountAddressOutput)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}

		if tokens == nil {
			tokens = make([]*Token, out.Pagination.TotalCount)
		}

		for _, d := range out.Data {
			tokens[i] = parseDatum(d)
			i++
		}

		next = out.Pagination.NextURL
	}

	return tokens, nil
}
