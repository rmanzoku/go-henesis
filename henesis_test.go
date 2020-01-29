package henesis

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/cheekybits/is"
)

var (
	h = &Henesis{}
)

func initializeTest(t *testing.T) is.I {
	is := is.New(t)
	var err error

	h, err = NewHenesis(os.Getenv("CLIENT_ID"))
	is.Nil(err)
	return is
}

func print(in interface{}) {
	if reflect.TypeOf(in).Kind() == reflect.Struct {
		in, _ = json.Marshal(in)
		in = string(in.([]byte))
	}
	fmt.Println(in)
}
