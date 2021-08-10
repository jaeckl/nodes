package main

import (
    "godspnodes/core"
    "godspnodes/api"
    "fmt"
    "strconv"
)


func main() {
    rt := core.NewRuntime()
    rt.RegisterObject("math.adder",NewAdder)
    rt.RegisterObject("lang.value",NewNumber)

    rt.New("math.adder:1","math.adder","")
    rt.New("lang.value:1","lang.value","5")
    rt.New("lang.value:2","lang.value","0")

    rt.Connect("lang.value:1","a1","math.adder:1","a1")
    rt.Connect("lang.value:1","a1","math.adder:1","a2")
    rt.Connect("math.adder:1","a1","lang.value:2","a1")
    rt.ConnectMsg("math.adder:1","lang.value:2")

    rt.ReceiveMessage("math.adder:1","Pulse")
}





type NumberObject struct {
    data int
}

func NewNumber() api.NodeObject {
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





type AdderObject struct {

}

func NewAdder() api.NodeObject {
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
