package api
type InputNode interface {
    GetValue()interface{}
    Name() string
}

type OutputNode interface {
    SetValue(value interface{})
    Name() string
}

type RuntimeContext interface {
    Inputs(nodename string) InputNode
    Outputs(nodename string) OutputNode
    NewInputNode(nodename string,typename string)
    NewOutputNode(nodename string,typename string)
    Broadcast(msg string)
}

type NodeObject interface {
    Init(ctx RuntimeContext,args string)
    ReceiveMessage(ctx RuntimeContext,msg string)
}

type RuntimeAPI interface {
    RegisterObject(objectname string,constructor func() NodeObject)
    New(id string, objectname string,args string)
    Delete(id string)
    Connect(from string,outname string, to string, inname string)
    SendMessage(id string,msg string)
    ReceiveMessage(id string,msg string)
}

