//go:generate go build -buildmode=plugin -o node/math/add.nd node/math/Adder.go
package main
import "github.com/jaeckl/nodes/api"
import "strconv"
import "fmt"
type NumberObject struct {
    data int
}

func New() api.NodeObject {
    object := &NumberObject{}
    return object
}

func (ad *NumberObject) Init(ctx api.RuntimeContext,args string) {
    ctx.NewInputNode("a1","int")
    ctx.NewOutputNode("a1","int")
    i, _ := strconv.Atoi(args)
    ctx.Outputs("a1").SetValue(i)
}

func (ad *NumberObject) ReceiveMessage(ctx api.RuntimeContext,msg string) {
    fmt.Println(ctx.Inputs("a1").GetValue().(int))
}
