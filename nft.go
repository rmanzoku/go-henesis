package henesis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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

type getTokensByAccountAddressInput struct {
	queries
	AccountAddress    string
	ContractAddresses []string
}

func (in getTokensByAccountAddressInput) Path() string {
	contracts := "&contractAddresses=" + strings.Join(in.ContractAddresses, ",")
	return fmt.Sprintf("/nft/v1/accounts/%s/tokens?", in.AccountAddress) + in.queries.Encode() + contracts
}

func (h Henesis) GetTokensByAccountAddress(accountAddress string, contractAddresses []string) (tokens []*Token, err error) {
	ctx := context.TODO()
	return h.GetTokensByAccountAddressWithContext(ctx, accountAddress, contractAddresses)
}

func (h Henesis) GetTokensByAccountAddressWithContext(ctx context.Context, accountAddress string, contractAddresses []string) (tokens []*Token, err error) {
	in := &getTokensByAccountAddressInput{
		AccountAddress: accountAddress,
		queries: queries{
			Page:           0,
			Size:           200,
			OrderBy:        "transfer_block_number",
			OrderDirection: "desc",
		},
		ContractAddresses: contractAddresses,
	}

	next := h.API + in.Path()
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
			i++
		}

		next = out.Pagination.NextURL
	}

	return tokens, nil
}
