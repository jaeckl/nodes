package core
import "github.com/jaeckl/nodes/api"
import "fmt"
import "log"
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
    log.Printf("Registered %v successfull\n",objectname)
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
        succ := datamodel.Object.Init(&RuntimeContextImpl{Runtime:rt,scope:datamodel},args)
        if !succ {
            log.Printf("Error: Bad arguments\n")
            return
        }
        rt.memoryModel[id] = datamodel
        log.Printf("Created new object of type %v successfull\n",objectname)
    } else {
        log.Printf("Error: Oject of type %v not registered\n",objectname)
    }
}

func (rt *Runtime) Delete(id string) {
    delete(rt.memoryModel,id)
    log.Printf("Removed oject with id %v\n",id)
}

/**************************** Object Messaging ****************************/

func (rt *Runtime) SendMessage(id string,msg string) {
    if val, ok := rt.memoryModel[id]; ok {
        for _,receiver := range val.MsgReceivers {
        receiver.Object.ReceiveMessage(&RuntimeContextImpl{Runtime:rt,scope:receiver },msg)
    }
    }
}
func (rt *Runtime) ReceiveMessage(id string,msg string) {
    if val, ok := rt.memoryModel[id]; ok {
        val.Object.ReceiveMessage(&RuntimeContextImpl{Runtime:rt,scope:val},msg)
    }
}

/**************************** Node Management ****************************/

func (rt *Runtime) Connect(objectFrom string,nodeFrom string, objectTo string, nodeTo string){
    ObjectFrom,ok := rt.memoryModel[objectFrom]
    if !ok {
        log.Printf("Error: Requested sender %v does not exist%v\n",objectFrom)
        return
    }
    ObjectTo,ok := rt.memoryModel[objectTo]
    if !ok {
        log.Printf("Error: Requested receiver %v does not exist%v\n",objectFrom)
        return
    }
    var NodeFrom *Node
    for _,node := range ObjectFrom.Outputs {
        if nodeFrom== node.Name {
            NodeFrom = node
        }
    }
    if NodeFrom == nil {
        log.Printf("Error: Requested sender %v has no node named %v\n",objectFrom,nodeFrom)
        return
    }

    var NodeTo *Node
    for _,node := range ObjectTo.Inputs {
        if nodeTo== node.Name {
            NodeTo = node
        }
    }
    if NodeTo == nil {
        log.Printf("Error: Requested receiver %v has no node named %v\n",objectTo,nodeTo)
        return
    }
    if NodeFrom.Type != NodeTo.Type {
        log.Printf("Error: Type of %v is %v but type of %v is %v\n",nodeFrom,NodeFrom.Type,nodeTo,NodeTo.Type)
        return
    }
    NodeFrom.Connection = append(NodeFrom.Connection,NodeTo)
    NodeTo.Connection[0] = NodeFrom
    log.Printf("Connecting %v:%v to %v:%v successfull\n",objectFrom,nodeFrom,objectTo,nodeTo)
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
    *Runtime
    scope *ObjectData
}
func NewRuntimeContext(runtime *Runtime,scope *ObjectData) *RuntimeContextImpl {
    return &RuntimeContextImpl{Runtime:runtime,scope:scope}
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
        receiver.Object.ReceiveMessage(&RuntimeContextImpl{Runtime:ctx.Runtime,scope:receiver },msg)
    }
}
