package henesis

import (
	"testing"
)

var (
	owner = "0xd868711BD9a2C6F1548F5f4737f71DA67d821090"
	// owner    = "0xd868711BD9a2C6F1548F5f4737f71DA67d821091"
	contract = "0xdceaf1652a131f32a821468dc03a92df0edd86ea"
)

func TestGetContract(t *testing.T) {
	is := initializeTest(t)
	var err error
	c, err := h.GetContract(contract)
	is.Nil(err)
	print(*c)
}

func TestGetAllContracts(t *testing.T) {
	is := initializeTest(t)
	var err error
	c, err := h.GetAllContracts()
	is.Nil(err)
	print(*c[0])
}

func TestGetContractsByAccountAddresss(t *testing.T) {
	is := initializeTest(t)
	var err error
	c, err := h.GetContractsByAccountAddresss(owner)
	is.Nil(err)
	if len(c) != 0 {
		print(*c[0])
	}
}

func TestGetTokensByAccountAddress(t *testing.T) {
	is := initializeTest(t)
	var err error
	tokens, err := h.GetTokensByAccountAddress(
		owner,
		[]string{contract},
	)
	is.Nil(err)
	if len(tokens) != 0 {
		tk := *tokens[0]
		print(tk)
		// m, err := erc721metadata.FetchERC721Metadata(tk.URI)
		// is.Nil(err)
		// print(*m)
	}
}

func TestGetOwnersByContractAddress(t *testing.T) {
	is := initializeTest(t)
	var err error
	ret, err := h.GetOwnersByContractAddress(contract)
	is.Nil(err)
	print(len(ret))
	if len(ret) != 0 {
		print(*ret[0])
	}
}

func TestGetTokensByContractAddress(t *testing.T) {
	is := initializeTest(t)
	var err error
	ret, err := h.GetTokensByContractAddress(contract)
	is.Nil(err)
	print(len(ret))
	if len(ret) != 0 {
		print(*ret[0])
	}
}
