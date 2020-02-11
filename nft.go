package henesis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Pagination struct {
	TotalCount int64  `json:"totalCount"`
	PrevURL    string `json:"prevUrl"`
	NextURL    string `json:"nextUrl"`
}

type queries struct {
	Page           int
	Size           int
	OrderBy        string
	OrderDirection string
}

func (q queries) Encode() string {
	v := make(url.Values)
	v.Set("page", strconv.Itoa(q.Page))
	v.Set("size", strconv.Itoa(q.Size))
	v.Set("order_by", q.OrderBy)
	v.Set("order_direction", q.OrderDirection)
	return v.Encode()
}

type Date string

func ParseDate(year int, month time.Month, day int) Date {
	return ParseDateFromInt(year, int(month), day)
}

func ParseDateFromInt(year int, month int, day int) Date {
	return Date(fmt.Sprintf("%d-%02d-%02d", year, month, day))
}

func (d Date) String() string {
	return string(d)
}

type Usage struct {
	Count int64 `json:"count"`
	Date  Date  `json:"date"`
}

func (h Henesis) NFTUsage(start Date, end Date) (usages []*Usage, err error) {
	v := make(url.Values)
	v.Set("start", start.String())
	v.Set("end", end.String())
	v.Set("page", "0")
	v.Set("size", "31")
	v.Set("orderBy", "date")
	v.Set("orderDirection", "desc")

	ctx := context.TODO()
	b, err := h.getPath(ctx, "/nft/v1/stats/jsonRpcDailyStats?"+v.Encode())
	if err != nil {
		return
	}
	o := &struct {
		Usages []*Usage `json:"data"`
	}{}
	return o.Usages, json.Unmarshal(b, o)
}

func (h Henesis) GetContract(contractAddress string) (contract *Contract, err error) {
	ctx := context.TODO()
	return h.GetContractWithContext(ctx, contractAddress)
}

func (h Henesis) GetContractWithContext(ctx context.Context, contractAddress string) (contract *Contract, err error) {
	path := fmt.Sprintf("/nft/v1/contracts/%s", contractAddress)
	b, err := h.getPath(ctx, path)
	if err != nil {
		return
	}
	contract = new(Contract)
	return contract, json.Unmarshal(b, contract)
}

func (h Henesis) GetAllContracts() (contracts []*Contract, err error) {
	ctx := context.TODO()
	return h.GetAllContractsWithContext(ctx)
}

func (h Henesis) GetAllContractsWithContext(ctx context.Context) (contracts []*Contract, err error) {
	path := fmt.Sprintf("/nft/v1/contracts/")
	b, err := h.getPath(ctx, path)
	if err != nil {
		return
	}
	contracts = []*Contract{}
	return contracts, json.Unmarshal(b, &contracts)
}

func (h Henesis) GetContractsByAccountAddresss(accountAddress string) (contracts []*Contract, err error) {
	ctx := context.TODO()
	return h.GetContractsByAccountAddresssWithContext(ctx, accountAddress)
}

func (h Henesis) GetContractsByAccountAddresssWithContext(ctx context.Context, accountAddress string) (contracts []*Contract, err error) {
	path := fmt.Sprintf("/nft/v1/accounts/%s/contracts", accountAddress)
	b, err := h.getPath(ctx, path)
	if err != nil {
		return
	}
	o := &struct {
		Contracts []*Contract `json:"data"`
	}{}
	return o.Contracts, json.Unmarshal(b, o)
}

func (h Henesis) GetTokensByAccountAddress(accountAddress string, contractAddresses []string) (tokens []*Token, err error) {
	ctx := context.TODO()
	return h.GetTokensByAccountAddressWithContext(ctx, accountAddress, contractAddresses)
}

func (h Henesis) GetTokensByAccountAddressWithContext(ctx context.Context, accountAddress string, contractAddresses []string) (tokens []*Token, err error) {
	q := queries{
		Page:           0,
		Size:           200,
		OrderBy:        "transfer_block_number",
		OrderDirection: "desc",
	}
	contracts := "&contractAddresses=" + strings.Join(contractAddresses, ",")
	path := fmt.Sprintf("/nft/v1/accounts/%s/tokens?", accountAddress) + q.Encode() + contracts
	next := h.API + path
	i := 0
	init := true
	if next != "" {
		b, err := h.getURL(ctx, next)
		if err != nil {
			return nil, err
		}
		out := &struct {
			Tokens     []*Token    `json:"data"`
			Pagination *Pagination `json:"pagination"`
		}{}
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}

		if init {
			tokens = make([]*Token, out.Pagination.TotalCount)
			init = false
		}

		for _, d := range out.Tokens {
			tokens[i] = d
			tokens[i].ContractAddress = d.Contract.Address
			i++
		}

		next = out.Pagination.NextURL
	}

	return tokens, nil
}

func (h Henesis) GetOwnersByContractAddress(contractAddress string) (owners []*Owner, err error) {
	ctx := context.TODO()
	return h.GetOwnersByContractAddressWithContext(ctx, contractAddress)
}

func (h Henesis) GetOwnersByContractAddressWithContext(ctx context.Context, contractAddress string) (owners []*Owner, error error) {
	q := queries{
		Page:           0,
		Size:           200,
		OrderBy:        "token_count",
		OrderDirection: "desc",
	}
	path := fmt.Sprintf("/nft/v1/contracts/%s/owners?", contractAddress) + q.Encode()
	next := h.API + path
	i := 0
	init := true
	if next != "" {
		b, err := h.getURL(ctx, next)
		if err != nil {
			return nil, err
		}
		out := &struct {
			Owners     []*Owner    `json:"data"`
			Pagination *Pagination `json:"pagination"`
		}{}
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}

		if init {
			owners = make([]*Owner, out.Pagination.TotalCount)
			init = false
		}

		for _, d := range out.Owners {
			owners[i] = d
			i++
		}

		next = out.Pagination.NextURL
	}

	return
}

func (h Henesis) GetTokensByContractAddress(contractAddress string) (tokens []*Token, err error) {
	ctx := context.TODO()
	return h.GetTokensByContractAddressWithContext(ctx, contractAddress)
}

func (h Henesis) GetTokensByContractAddressWithContext(ctx context.Context, contractAddress string) (tokens []*Token, error error) {
	q := queries{
		Page:           0,
		Size:           200,
		OrderBy:        "transfer_block_number",
		OrderDirection: "desc",
	}
	path := fmt.Sprintf("/nft/v1/contracts/%s/tokens?", contractAddress) + q.Encode()
	next := h.API + path
	i := 0
	init := true
	if next != "" {
		b, err := h.getURL(ctx, next)
		if err != nil {
			return nil, err
		}
		out := &struct {
			Data struct {
				Address     string   `json:"address"`
				Name        string   `json:"name"`
				Symbol      string   `json:"symbol"`
				TotalSupply string   `json:"totalSupply"`
				Tokens      []*Token `json:"tokens"`
			}
			Pagination *Pagination `json:"pagination"`
		}{}
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, out)
		if err != nil {
			return nil, err
		}

		if init {
			tokens = make([]*Token, out.Pagination.TotalCount)
			init = false
		}

		for _, d := range out.Data.Tokens {
			tokens[i] = d
			tokens[i].Contract = &Contract{
				Address:     out.Data.Address,
				Name:        out.Data.Name,
				Symbol:      out.Data.Symbol,
				TotalSupply: out.Data.TotalSupply,
			}
			i++
		}

		next = out.Pagination.NextURL
	}

	return
}
