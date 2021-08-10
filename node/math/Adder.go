//go:generate go build -buildmode=plugin -o node/math/add.nd node/math/Adder.go
package main

import "github.com/jaeckl/nodes/api"
import "fmt"
type AdderObject struct {

}

func New() api.NodeObject {
    object := &AdderObject{}
    return object
}
func (ad *AdderObject) Init(ctx api.RuntimeContext,args string) bool {
    ctx.NewInputNode("left","int")
    ctx.NewInputNode("right","int")
    ctx.NewOutputNode("out","int")
    return true
}

func (ad *AdderObject) ReceiveMessage(ctx api.RuntimeContext,msg string) {
    fmt.Printf("Receiving Message: %v\n",msg)

    val := ctx.Inputs("left").GetValue().(int64) + ctx.Inputs("right").GetValue().(int64)
    ctx.Outputs("out").SetValue(val)
    ctx.Broadcast("Pulse")
}
