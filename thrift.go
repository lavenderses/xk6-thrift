package thrift

import (
	"fmt"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/thrift", &TModule{})
}

type TModule struct{}

func (m *TModule) Echo(arg string) {
	fmt.Println("Hello " + arg)
}
