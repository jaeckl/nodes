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
func (ad *AdderObject) Init(ctx api.RuntimeContext,args string) {
    ctx.NewInputNode("a1","int")
    ctx.NewInputNode("a2","int")
    ctx.NewOutputNode("a1","int")
}

func (ad *AdderObject) ReceiveMessage(ctx api.RuntimeContext,msg string) {
    fmt.Printf("Receiving Message: %v\n",msg)
    val := ctx.Inputs("a1").GetValue().(int) + ctx.Inputs("a2").GetValue().(int)
    ctx.Outputs("a1").SetValue(val)
    ctx.Broadcast("Pulse")
    fmt.Println(val)
}
