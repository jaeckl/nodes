package core
import "github.com/jaeckl/nodes/api"
import "fmt"

type ObjectName = string
type NodeName = string

type Runtime struct {
    classRegistry map[string] func()api.NodeObject
    memoryModel map[ObjectName] *ObjectData

}

func NewRuntime() *Runtime {
    return &Runtime{
                classRegistry:  make(map[string] func() api.NodeObject),
                memoryModel:    make(map[string] *ObjectData),

    }
}
/**************************** Object Registry ****************************/

func (rt *Runtime) RegisterObject(objectname string,constructor func() api.NodeObject){
    rt.classRegistry[objectname] = constructor
}


/**************************** Object Lifetime ****************************/

func (rt *Runtime) New(id string, objectname string,args string) {
    if val, ok := rt.classRegistry[objectname]; ok {
        datamodel := &ObjectData{
            Object:         val(),
            Inputs:         make([]*Node,0),
            Outputs:        make([]*Node,0),
            MsgReceivers:   make([]*ObjectData,0),
        }

        datamodel.Object.Init(&RuntimeContextImpl{runtime:rt,scope:datamodel},args)

        rt.memoryModel[id] = datamodel
    } else {
        // LOG Object does not exists.
    }
}

func (rt *Runtime) Delete(id string) {
    delete(rt.memoryModel,id)
}

/**************************** Object Messaging ****************************/

func (rt *Runtime) SendMessage(id string,msg string) {
    if val, ok := rt.memoryModel[id]; ok {
        for _,receiver := range val.MsgReceivers {
        receiver.Object.ReceiveMessage(&RuntimeContextImpl{runtime:rt,scope:receiver },msg)
    }
    }
}
func (rt *Runtime) ReceiveMessage(id string,msg string) {
    if val, ok := rt.memoryModel[id]; ok {
        val.Object.ReceiveMessage(&RuntimeContextImpl{runtime:rt,scope:val},msg)
    }
}

/**************************** Node Management ****************************/

func (rt *Runtime) Connect(objectFrom string,nodeFrom string, objectTo string, nodeTo string){
    ObjectFrom,ok := rt.memoryModel[objectFrom]
    if !ok {
        fmt.Println("Sending Object does not exists")
        return
    }
    ObjectTo,ok := rt.memoryModel[objectTo]
    if !ok {
        fmt.Println("Receiving Object does not exists")
        return
    }
    var NodeFrom *Node
    for _,node := range ObjectFrom.Outputs {
        if nodeFrom== node.Name {
            NodeFrom = node
        }
    }

    var NodeTo *Node
    for _,node := range ObjectTo.Inputs {
        if nodeTo== node.Name {
            NodeTo = node
        }
    }

    NodeFrom.Connection = append(NodeFrom.Connection,NodeTo)
    NodeTo.Connection[0] = NodeFrom
}

func (rt *Runtime) ConnectMsg(objectFrom string,objectTo string) {
    ObjectFrom,ok := rt.memoryModel[objectFrom]
    if !ok {
        fmt.Println("Sending Object does not exists")
        return
    }
    ObjectTo,ok := rt.memoryModel[objectTo]
    if !ok {
        fmt.Println("Receiving Object does not exists")
        return
    }
    ObjectFrom.MsgReceivers = append(ObjectFrom.MsgReceivers,ObjectTo)
}

//
/**************************** Node Structure ****************************/
//
type TypeString = string

type Node struct {
    Name string
    Type TypeString
    Data interface{}
    Connection []*Node
    Object api.NodeObject
}

type InputNodeImpl struct{
    node *Node
}
type OutputNodeImpl struct{
    node *Node
}

func (nd *Node) SetValue(value interface{}) {
    nd.Data = value
}

func (nd *Node) GetValue() interface{} {
    return nd.Data
}

func (nd *OutputNodeImpl) SetValue(value interface{}) {
    nd.node.SetValue(value)
}

func (nd *OutputNodeImpl) Name()string {
    return nd.node.Name
}

func (nd *InputNodeImpl) GetValue() interface{} {
    return nd.node.Connection[0].GetValue()
}

func (nd *InputNodeImpl) Name() string{
    return nd.node.Name
}

type ObjectData struct {
    Object api.NodeObject
    Inputs []*Node
    Outputs []*Node
    MsgReceivers []*ObjectData
}

func (od *ObjectData) Input(nodeName string) api.InputNode {
    for _,node := range od.Inputs {
        if nodeName == node.Name {
            return &InputNodeImpl{node:node}
        }
    }
    return nil
}

func (od *ObjectData) Output(nodeName string) api.OutputNode {
    for _,node := range od.Outputs {
        if nodeName == node.Name {
            return &OutputNodeImpl{node:node}
        }
    }
    return nil
}

//
/***********************************************************/
//



type RuntimeContextImpl struct {
    runtime *Runtime
    scope *ObjectData
}
func NewRuntimeContext(runtime *Runtime,scope *ObjectData) *RuntimeContextImpl {
    return &RuntimeContextImpl{runtime:runtime,scope:scope}
}
func (ctx *RuntimeContextImpl) NewInputNode(nodename string,typename TypeString) {
    ctx.scope.Inputs = append(ctx.scope.Inputs,&Node{Name:nodename,Data:nil,Type:typename,Connection:make([]*Node,1),Object:ctx.scope.Object})
}

func (ctx *RuntimeContextImpl) NewOutputNode(nodename string,typename TypeString) {
    ctx.scope.Outputs = append(ctx.scope.Outputs,&Node{Name:nodename,Data:nil,Type:typename,Connection:make([]*Node,1),Object:ctx.scope.Object})
}
func (ctx *RuntimeContextImpl) Inputs(nodeName string) api.InputNode {
    return ctx.scope.Input(nodeName)
}

func (ctx *RuntimeContextImpl) Outputs(nodeName string) api.OutputNode {
    return ctx.scope.Output(nodeName)
}

func (ctx *RuntimeContextImpl) Broadcast(msg string) {
    for _,receiver := range ctx.scope.MsgReceivers {
        receiver.Object.ReceiveMessage(&RuntimeContextImpl{runtime:ctx.runtime,scope:receiver },msg)
    }
}
