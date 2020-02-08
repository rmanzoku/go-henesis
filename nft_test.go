package henesis

import (
	"testing"
)

var (
	owner    = "0xd868711BD9a2C6F1548F5f4737f71DA67d821090"
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
	print(*c[0])
}

func TestGetTokensByAccountAddress(t *testing.T) {
	is := initializeTest(t)
	var err error
	tokens, err := h.GetTokensByAccountAddress(
		owner,
		[]string{contract},
	)
	is.Nil(err)
	tk := *tokens[0]
	print(tk)
	// m, err := erc721metadata.FetchERC721Metadata(tk.URI)
	// is.Nil(err)
	// print(*m)
}
