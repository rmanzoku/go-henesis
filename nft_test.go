package henesis

import (
	"testing"
)

func TestGetTokensByAccountAddress(t *testing.T) {
	is := initializeTest(t)
	var err error
	tokens, err := h.GetTokensByAccountAddress(
		"0xd868711BD9a2C6F1548F5f4737f71DA67d821090",
		[]string{"0xdceaf1652a131f32a821468dc03a92df0edd86ea"},
	)
	is.Nil(err)
	tk := *tokens[0]
	print(tk)
	// m, err := erc721metadata.FetchERC721Metadata(tk.URI)
	// is.Nil(err)
	// print(*m)
}
